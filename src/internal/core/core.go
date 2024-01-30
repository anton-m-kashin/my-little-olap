package core

import (
	"time"
)

type (
	MetricBase struct {
		SessionID string
		Platform  string
		Timestamp time.Time
	}
	DurationMetricBase struct {
		MetricBase
		Duration time.Duration
	}
	ScreenOpeningTimeMetric struct {
		DurationMetricBase
		ScreenName string
	}
	RequestTimeMetric struct {
		DurationMetricBase
		RequestURL string
	}
)

type MetricsRepository interface {
	AddScreenOpeningTimeMetrics(metrics []ScreenOpeningTimeMetric)
	AddRequestTimeMetrics(metrics []RequestTimeMetric)
	GetLastScreenOpeningTimeMetrics(count uint) []ScreenOpeningTimeMetric
	GetLastAddRequestTimeMetrics(count uint) []RequestTimeMetric
}
