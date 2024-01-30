package db

import (
	"context"
	"my-little-olap/internal/core"
	"time"
)

func (ch *ClickhouseDB) AddScreenOpeningTimeMetrics(
	metrics []core.ScreenOpeningTimeMetric,
) {
	batch, err := ch.conn.PrepareBatch(
		context.Background(),
		`INSERT INTO screen_opening_time VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		ch.logger.Error.Fatalf(
			"Screen opening time metrics insertion error: %s\n",
			err,
		)
	}
	for _, m := range metrics {
		err := batch.Append(
			m.SessionID,
			m.Platform,
			m.Timestamp,
			m.Duration,
			m.ScreenName,
		)
		if err != nil {
			ch.logger.Error.Fatalf(
				"Screen opening time batch creation error: %s\n",
				err,
			)
		}
	}
	err = batch.Send()
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
	batch, err := ch.conn.PrepareBatch(
		context.Background(),
		`INSERT INTO request_time VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		ch.logger.Error.Fatalf(
			"Request time metrics insertion error: %s\n",
			err,
		)
	}
	for _, m := range metrics {
		err := batch.Append(
			m.SessionID,
			m.Platform,
			m.Timestamp,
			m.Duration,
			m.RequestURL,
		)
		if err != nil {
			ch.logger.Error.Fatalf(
				"Request time batch creation error: %s\n",
				err,
			)
		}
	}
	err = batch.Send()
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
