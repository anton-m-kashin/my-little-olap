package db

import (
	"context"
	"my-little-olap/internal/core"
	"my-little-olap/internal/utils"
	"time"
)

func (ch *ClickhouseDB) AddScreenOpeningTimeMetrics(
	metrics []core.ScreenOpeningTimeMetric,
) {
	err := insert[core.ScreenOpeningTimeMetric](
		ch,
		"screen_opening_time",
		metrics,
		func(m core.ScreenOpeningTimeMetric) []any {
			return []any{
				m.SessionID,
				m.Platform,
				m.Timestamp,
				m.Duration,
				m.ScreenName,
			}
		},
	)
	if err != nil {
		ch.logger.Error.Fatalf(
			"Screen opening time batch sending error: %s\n",
			err,
		)
	}
}

func (ch *ClickhouseDB) AddRequestTimeMetrics(
	metrics []core.RequestTimeMetric,
) {
	err := insert[core.RequestTimeMetric](
		ch,
		"screen_opening_time",
		metrics,
		func(m core.RequestTimeMetric) []any {
			return []any{
				m.SessionID,
				m.Platform,
				m.Timestamp,
				m.Duration,
				m.RequestURL,
			}
		},
	)
	if err != nil {
		ch.logger.Error.Fatalf(
			"Screen opening time batch sending error: %s\n",
			err,
		)
	}
}

func (ch *ClickhouseDB) GetLastScreenOpeningTimeMetrics(
	count uint,
) []core.ScreenOpeningTimeMetric {
	var results []struct {
		SessionID  string        `ch:"session_id"`
		Platform   string        `ch:"platform"`
		Timestamp  time.Time     `ch:"timestamp"`
		Duration   time.Duration `ch:"duration"`
		ScreenName string        `ch:"screen_name"`
	}
	err := ch.conn.Select(
		context.Background(),
		&results,
		`SELECT * FROM screen_opening_time LIMIT ?`,
		count,
	)
	if err != nil {
		ch.logger.Error.Printf("Screen opening time metrics: %s\n", err)
		return []core.ScreenOpeningTimeMetric{}
	}
	metrics := make([]core.ScreenOpeningTimeMetric, len(results))
	for i, r := range results {
		metrics[i] = core.ScreenOpeningTimeMetric{
			DurationMetricBase: core.DurationMetricBase{
				MetricBase: core.MetricBase{
					SessionID: r.SessionID,
					Platform:  r.Platform,
					Timestamp: r.Timestamp,
				},
				Duration: r.Duration,
			},
			ScreenName: r.ScreenName,
		}
	}
	return metrics
}

func (ch *ClickhouseDB) GetLastAddRequestTimeMetrics(
	count uint,
) []core.RequestTimeMetric {
	var results []struct {
		SessionID  string        `ch:"session_id"`
		Platform   string        `ch:"platform"`
		Timestamp  time.Time     `ch:"timestamp"`
		Duration   time.Duration `ch:"duration"`
		RequestURL string        `ch:"request_url"`
	}
	err := ch.conn.Select(
		context.Background(),
		&results,
		`SELECT * FROM request_time LIMIT ?`,
		count,
	)
	if err != nil {
		ch.logger.Error.Println("No screen opening time metrics.")
		return []core.RequestTimeMetric{}
	}
	metrics := make([]core.RequestTimeMetric, len(results))
	for i, r := range results {
		metrics[i] = core.RequestTimeMetric{
			DurationMetricBase: core.DurationMetricBase{
				MetricBase: core.MetricBase{
					SessionID: r.SessionID,
					Platform:  r.Platform,
					Timestamp: r.Timestamp,
				},
				Duration: r.Duration,
			},
			RequestURL: r.RequestURL,
		}
	}
	return metrics
}

var _ = (*core.MetricsRepository)(nil)

func insert[Item any](
	ch *ClickhouseDB,
	table string,
	items []Item,
	transform func(Item) []any,
) error {
	nextItem := utils.Iterate(items)
	return ch.insertBatch(
		table,
		func() *[]any {
			item := nextItem()
			if item == nil {
				return nil
			}
			fields := transform(*item)
			return &fields
		},
	)
}
