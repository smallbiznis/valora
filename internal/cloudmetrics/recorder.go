package cloudmetrics

import (
	"strings"
	"sync"
)

type Recorder interface {
	RecordUsageEvent(orgID, meter string)
	RecordInvoiceGenerated(orgID string)
	UpdateActiveSubscriptions(orgID string, count int)
	RecordEngineError(orgID, operation string)
}

type recorder struct {
	metrics       *metrics
	defaultOrgID  string
	defaultOrgName string
}

type noopRecorder struct{}

func (noopRecorder) RecordUsageEvent(string, string)              {}
func (noopRecorder) RecordInvoiceGenerated(string)                {}
func (noopRecorder) UpdateActiveSubscriptions(string, int)         {}
func (noopRecorder) RecordEngineError(string, string)              {}

var (
	activeRecorder Recorder = noopRecorder{}
	recorderMu     sync.RWMutex
)

func setRecorder(rec Recorder) {
	if rec == nil {
		return
	}
	recorderMu.Lock()
	activeRecorder = rec
	recorderMu.Unlock()
}

func RecordUsageEvent(orgID, code string) {
	recorderMu.RLock()
	rec := activeRecorder
	recorderMu.RUnlock()
	rec.RecordUsageEvent(orgID, code)
}

func RecordInvoiceGenerated(orgID string) {
	recorderMu.RLock()
	rec := activeRecorder
	recorderMu.RUnlock()
	rec.RecordInvoiceGenerated(orgID)
}

func UpdateActiveSubscriptions(orgID string, count int) {
	recorderMu.RLock()
	rec := activeRecorder
	recorderMu.RUnlock()
	rec.UpdateActiveSubscriptions(orgID, count)
}

func RecordEngineError(orgID, operation string) {
	recorderMu.RLock()
	rec := activeRecorder
	recorderMu.RUnlock()
	rec.RecordEngineError(orgID, operation)
}

func (r *recorder) RecordUsageEvent(orgID, meter string) {
	if r == nil || r.metrics == nil {
		return
	}
	org := r.normalizeOrg(orgID)
	meterLabel := normalizeLabel(meter)
	r.metrics.usageEvents.WithLabelValues(org, meterLabel).Inc()
}

func (r *recorder) RecordInvoiceGenerated(orgID string) {
	if r == nil || r.metrics == nil {
		return
	}
	org := r.normalizeOrg(orgID)
	r.metrics.invoicesGenerated.WithLabelValues(org).Inc()
}

func (r *recorder) UpdateActiveSubscriptions(orgID string, count int) {
	if r == nil || r.metrics == nil {
		return
	}
	if count < 0 {
		count = 0
	}
	org := r.normalizeOrg(orgID)
	r.metrics.activeSubscriptions.WithLabelValues(org).Set(float64(count))
}

func (r *recorder) RecordEngineError(orgID, operation string) {
	if r == nil || r.metrics == nil {
		return
	}
	org := r.normalizeOrg(orgID)
	opLabel := normalizeLabel(operation)
	r.metrics.engineErrors.WithLabelValues(org, opLabel).Inc()
}

func (r *recorder) normalizeOrg(orgID string) string {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		orgID = strings.TrimSpace(r.defaultOrgID)
	}
	if orgID == "" {
		return "unknown"
	}
	return orgID
}

func normalizeLabel(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	return value
}
