package cloudmetrics

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/snappy"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
	"github.com/smallbiznis/valora/internal/config"
	collectormetricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

const (
	exportInterval = 15 * time.Second
	exportTimeout  = 5 * time.Second
)

var registerOnce sync.Once

// Register configures cloud metrics emission. Failures are logged and never block billing.
func Register(lc fx.Lifecycle, cfg config.Config, registry *prometheus.Registry, logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	if !shouldEnable(cfg) {
		return
	}

	exporterCfg, err := parseExporterConfig(cfg)
	if err != nil {
		logger.Warn("cloud metrics disabled", zap.Error(err))
		return
	}

	registerOnce.Do(func() {
		setRecorder(&recorder{
			metrics:      newMetrics(registry),
			defaultOrgID: cfg.Cloud.OrganizationID,
			defaultOrgName: cfg.Cloud.OrganizationName,
		})

		exp := newExporter(registry, exporterCfg, logger)
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				exp.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return exp.Stop(ctx)
			},
		})
	})
}

func shouldEnable(cfg config.Config) bool {
	fmt.Printf("cfg.IsCloud() && cfg.Cloud.Metrics.Enabled: %v\n", cfg.IsCloud())
	fmt.Printf("cfg.IsCloud() && cfg.Cloud.Metrics.Enabled: %v\n", cfg.Cloud.Metrics.Enabled)
	return cfg.IsCloud() && cfg.Cloud.Metrics.Enabled
}

type exporterConfig struct {
	kind           string
	endpoint       string
	authToken      string
	otlpAddress    string
	otlpSecure     bool
	serviceName    string
	serviceVersion string
	environment    string
}

func parseExporterConfig(cfg config.Config) (exporterConfig, error) {
	kind := strings.ToLower(strings.TrimSpace(cfg.Cloud.Metrics.Exporter))
	if kind == "" {
		return exporterConfig{}, errors.New("cloud.metrics.exporter is required")
	}
	endpoint := strings.TrimSpace(cfg.Cloud.Metrics.Endpoint)
	if endpoint == "" {
		return exporterConfig{}, errors.New("cloud.metrics.endpoint is required")
	}

	out := exporterConfig{
		kind:           kind,
		endpoint:       endpoint,
		authToken:      strings.TrimSpace(cfg.Cloud.Metrics.AuthToken),
		serviceName:    cfg.AppName,
		serviceVersion: cfg.AppVersion,
		environment:    cfg.Environment,
	}

	switch kind {
	case exporterPrometheusRemoteWrite:
		if _, err := url.ParseRequestURI(endpoint); err != nil {
			return exporterConfig{}, fmt.Errorf("invalid cloud.metrics.endpoint: %w", err)
		}
	case exporterOTLP:
		addr, secure, err := parseOTLPEndpoint(endpoint)
		if err != nil {
			return exporterConfig{}, err
		}
		out.otlpAddress = addr
		out.otlpSecure = secure
	default:
		return exporterConfig{}, fmt.Errorf("unsupported cloud.metrics.exporter: %s", kind)
	}

	return out, nil
}

func parseOTLPEndpoint(endpoint string) (string, bool, error) {
	if strings.Contains(endpoint, "://") {
		parsed, err := url.Parse(endpoint)
		if err != nil {
			return "", false, fmt.Errorf("invalid cloud.metrics.endpoint: %w", err)
		}
		if parsed.Host == "" {
			return "", false, errors.New("cloud.metrics.endpoint host is required")
		}
		secure := parsed.Scheme == "https" || parsed.Scheme == "grpcs"
		return parsed.Host, secure, nil
	}
	if strings.TrimSpace(endpoint) == "" {
		return "", false, errors.New("cloud.metrics.endpoint is required")
	}
	return endpoint, false, nil
}

type exporter struct {
	kind       string
	endpoint   string
	authToken  string
	registry   *prometheus.Registry
	logger     *zap.Logger
	httpClient *http.Client
	resource   *resourcepb.Resource

	otlpAddress string
	otlpSecure  bool
	grpcConn    *grpc.ClientConn

	stopCh    chan struct{}
	doneCh    chan struct{}
	errorOnce atomic.Bool
}

func newExporter(registry *prometheus.Registry, cfg exporterConfig, logger *zap.Logger) *exporter {
	resource := buildResource(cfg.serviceName, cfg.serviceVersion, cfg.environment)
	return &exporter{
		kind:        cfg.kind,
		endpoint:    cfg.endpoint,
		authToken:   cfg.authToken,
		registry:    registry,
		logger:      logger,
		httpClient:  &http.Client{Timeout: exportTimeout},
		resource:    resource,
		otlpAddress: cfg.otlpAddress,
		otlpSecure:  cfg.otlpSecure,
	}
}

func buildResource(serviceName, serviceVersion, environment string) *resourcepb.Resource {
	attrs := make([]*commonpb.KeyValue, 0, 3)
	if serviceName != "" {
		attrs = append(attrs, &commonpb.KeyValue{
			Key:   "service.name",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: serviceName}},
		})
	}
	if serviceVersion != "" {
		attrs = append(attrs, &commonpb.KeyValue{
			Key:   "service.version",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: serviceVersion}},
		})
	}
	if environment != "" {
		attrs = append(attrs, &commonpb.KeyValue{
			Key:   "deployment.environment",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: environment}},
		})
	}
	if len(attrs) == 0 {
		return &resourcepb.Resource{}
	}
	return &resourcepb.Resource{Attributes: attrs}
}

func (e *exporter) Start() {
	if e == nil {
		return
	}
	if e.stopCh != nil {
		return
	}
	e.stopCh = make(chan struct{})
	e.doneCh = make(chan struct{})

	go func() {
		defer close(e.doneCh)
		ticker := time.NewTicker(exportInterval)
		defer ticker.Stop()
		e.exportOnce()
		for {
			select {
			case <-ticker.C:
				e.exportOnce()
			case <-e.stopCh:
				return
			}
		}
	}()
}

func (e *exporter) Stop(ctx context.Context) error {
	if e == nil || e.stopCh == nil {
		return nil
	}
	close(e.stopCh)
	if e.grpcConn != nil {
		_ = e.grpcConn.Close()
	}
	select {
	case <-e.doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *exporter) exportOnce() {
	if e == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), exportTimeout)
	defer cancel()
	families, err := e.registry.Gather()
	if err != nil {
		e.logExportError(err)
		return
	}
	if len(families) == 0 {
		return
	}

	switch e.kind {
	case exporterPrometheusRemoteWrite:
		err = e.exportRemoteWrite(ctx, families)
	case exporterOTLP:
		err = e.exportOTLP(ctx, families)
	default:
		err = fmt.Errorf("unsupported exporter: %s", e.kind)
	}

	if err != nil {
		e.logExportError(err)
		return
	}
	e.errorOnce.Store(false)
}

func (e *exporter) logExportError(err error) {
	if err == nil {
		return
	}
	if e.errorOnce.CompareAndSwap(false, true) {
		e.logger.Warn("cloud metrics export failed", zap.Error(err))
	}
}

func (e *exporter) exportRemoteWrite(ctx context.Context, families []*dto.MetricFamily) error {
	series := buildRemoteWriteSeries(families, time.Now().UnixMilli())
	if len(series) == 0 {
		return nil
	}

	req := &prompb.WriteRequest{Timeseries: series}
	payload, err := proto.Marshal(protoadapt.MessageV2Of(req))
	if err != nil {
		return err
	}

	compressed := snappy.Encode(nil, payload)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.endpoint, bytes.NewReader(compressed))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("Content-Encoding", "snappy")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	if e.authToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+e.authToken)
	}

	resp, err := e.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("remote write returned %s", resp.Status)
	}
	return nil
}

func buildRemoteWriteSeries(families []*dto.MetricFamily, timestampMs int64) []prompb.TimeSeries {
	series := make([]prompb.TimeSeries, 0, len(families))
	for _, family := range families {
		switch family.GetType() {
		case dto.MetricType_COUNTER, dto.MetricType_GAUGE:
		default:
			continue
		}
		for _, metric := range family.GetMetric() {
			value := extractMetricValue(family.GetType(), metric)
			if value == nil {
				continue
			}
			labels := make([]prompb.Label, 0, len(metric.GetLabel())+1)
			labels = append(labels, prompb.Label{Name: "__name__", Value: family.GetName()})
			for _, label := range metric.GetLabel() {
				labels = append(labels, prompb.Label{Name: label.GetName(), Value: label.GetValue()})
			}
			sort.Slice(labels, func(i, j int) bool {
				return labels[i].Name < labels[j].Name
			})

			series = append(series, prompb.TimeSeries{
				Labels: labels,
				Samples: []prompb.Sample{{
					Value:     *value,
					Timestamp: timestampMs,
				}},
			})
		}
	}
	return series
}

func (e *exporter) exportOTLP(ctx context.Context, families []*dto.MetricFamily) error {
	if e.grpcConn == nil {
		if err := e.connectOTLP(ctx); err != nil {
			return err
		}
	}

	metrics := buildOTLPMetrics(families, uint64(time.Now().UnixNano()))
	if len(metrics) == 0 {
		return nil
	}

	scope := &commonpb.InstrumentationScope{Name: "valora.cloudmetrics"}
	rm := &metricspb.ResourceMetrics{
		Resource: e.resource,
		ScopeMetrics: []*metricspb.ScopeMetrics{
			{
				Scope:   scope,
				Metrics: metrics,
			},
		},
	}

	if e.authToken != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+e.authToken)
	}

	client := collectormetricspb.NewMetricsServiceClient(e.grpcConn)
	_, err := client.Export(ctx, &collectormetricspb.ExportMetricsServiceRequest{
		ResourceMetrics: []*metricspb.ResourceMetrics{rm},
	})
	return err
}

func (e *exporter) connectOTLP(ctx context.Context) error {
	if e.grpcConn != nil {
		return nil
	}
	var creds credentials.TransportCredentials
	if e.otlpSecure {
		creds = credentials.NewClientTLSFromCert(nil, "")
	} else {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.DialContext(ctx, e.otlpAddress, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		return err
	}
	e.grpcConn = conn
	return nil
}

func buildOTLPMetrics(families []*dto.MetricFamily, now uint64) []*metricspb.Metric {
	metrics := make([]*metricspb.Metric, 0, len(families))
	for _, family := range families {
		switch family.GetType() {
		case dto.MetricType_COUNTER:
			dataPoints := buildOTLPDataPoints(family.GetMetric(), now, true)
			if len(dataPoints) == 0 {
				continue
			}
			metrics = append(metrics, &metricspb.Metric{
				Name:        family.GetName(),
				Description: family.GetHelp(),
				Data: &metricspb.Metric_Sum{
					Sum: &metricspb.Sum{
						IsMonotonic:            true,
						AggregationTemporality: metricspb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
						DataPoints:             dataPoints,
					},
				},
			})
		case dto.MetricType_GAUGE:
			dataPoints := buildOTLPDataPoints(family.GetMetric(), now, false)
			if len(dataPoints) == 0 {
				continue
			}
			metrics = append(metrics, &metricspb.Metric{
				Name:        family.GetName(),
				Description: family.GetHelp(),
				Data: &metricspb.Metric_Gauge{
					Gauge: &metricspb.Gauge{
						DataPoints: dataPoints,
					},
				},
			})
		default:
			continue
		}
	}
	return metrics
}

func buildOTLPDataPoints(metrics []*dto.Metric, now uint64, isCounter bool) []*metricspb.NumberDataPoint {
	points := make([]*metricspb.NumberDataPoint, 0, len(metrics))
	for _, metric := range metrics {
		value := extractMetricValueWithCounterFlag(metric, isCounter)
		if value == nil {
			continue
		}
		points = append(points, &metricspb.NumberDataPoint{
			Attributes:   buildOTLPAttributes(metric.GetLabel()),
			TimeUnixNano: now,
			Value:        &metricspb.NumberDataPoint_AsDouble{AsDouble: *value},
		})
	}
	return points
}

func buildOTLPAttributes(labels []*dto.LabelPair) []*commonpb.KeyValue {
	if len(labels) == 0 {
		return nil
	}
	attrs := make([]*commonpb.KeyValue, 0, len(labels))
	for _, label := range labels {
		if label == nil {
			continue
		}
		attrs = append(attrs, &commonpb.KeyValue{
			Key:   label.GetName(),
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: label.GetValue()}},
		})
	}
	return attrs
}

func extractMetricValue(metricType dto.MetricType, metric *dto.Metric) *float64 {
	if metric == nil {
		return nil
	}
	switch metricType {
	case dto.MetricType_COUNTER:
		if metric.GetCounter() == nil {
			return nil
		}
		value := metric.GetCounter().GetValue()
		return &value
	case dto.MetricType_GAUGE:
		if metric.GetGauge() == nil {
			return nil
		}
		value := metric.GetGauge().GetValue()
		return &value
	default:
		return nil
	}
}

func extractMetricValueWithCounterFlag(metric *dto.Metric, isCounter bool) *float64 {
	if metric == nil {
		return nil
	}
	if isCounter {
		if metric.GetCounter() == nil {
			return nil
		}
		value := metric.GetCounter().GetValue()
		return &value
	}
	if metric.GetGauge() == nil {
		return nil
	}
	value := metric.GetGauge().GetValue()
	return &value
}
