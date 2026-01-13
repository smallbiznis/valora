package config

import (
	"errors"
	"log"
	"strings"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func DefaultBillingConfig() BillingConfig {
	return BillingConfig{
		AgingBuckets: []AgingBucket{
			{Label: "0-30", MinDays: 0, MaxDays: intPtr(30)},
			{Label: "31-60", MinDays: 31, MaxDays: intPtr(60)},
			{Label: "60+", MinDays: 61, MaxDays: nil},
		},
		RiskLevels: []RiskLevel{
			{Level: "high", MinOutstanding: 1_000_000, MinDays: 60},
			{Level: "medium", MinOutstanding: 250_000, MinDays: 31},
			{Level: "low", MinOutstanding: 0, MinDays: 0},
		},
	}
}

func intPtr(v int) *int { return &v }

type BillingConfigHolder struct {
	current atomic.Value // holds BillingConfig
}

func NewBillingConfigHolder() (*BillingConfigHolder, error) {
	v := viper.New()

	v.SetConfigName("billing")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/valora")

	// env hanya untuk path override (optional)
	v.SetEnvPrefix("VALORA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg BillingConfig
	if err := v.UnmarshalKey("billing", &cfg); err != nil {
		return nil, err
	}
	if err := validateBillingConfig(cfg); err != nil {
		return nil, err
	}

	holder := &BillingConfigHolder{}
	holder.current.Store(cfg)

	// 🔥 HOT RELOAD
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var updated BillingConfig
		if err := v.UnmarshalKey("billing", &updated); err != nil {
			log.Printf("[billing-config] reload failed: %v", err)
			return
		}
		if err := validateBillingConfig(updated); err != nil {
			log.Printf("[billing-config] invalid config ignored: %v", err)
			return
		}
		holder.current.Store(updated)
		log.Printf("[billing-config] reloaded from %s", e.Name)
	})

	return holder, nil
}

func (h *BillingConfigHolder) Get() BillingConfig {
	return h.current.Load().(BillingConfig)
}

func validateBillingConfig(cfg BillingConfig) error {
	if len(cfg.AgingBuckets) == 0 {
		return errors.New("billing.agingBuckets cannot be empty")
	}
	if len(cfg.RiskLevels) == 0 {
		return errors.New("billing.riskLevels cannot be empty")
	}
	return nil
}
