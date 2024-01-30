package srv

import (
	"encoding/json"
	"net/http"
	"time"

	"my-little-olap/internal/core"
	"my-little-olap/internal/utils"
)

type OLAPServer struct {
	repo   core.MetricsRepository
	logger *utils.Logger
}

func NewOLAPServer(r core.MetricsRepository, logger *utils.Logger) OLAPServer {
	return OLAPServer{r, logger}
}

func (s *OLAPServer) Run() error {
	srv := http.Server{
		Addr:     ":8080",
		ErrorLog: s.logger.Error,
		Handler:  s.newMux(),
	}
	return srv.ListenAndServe()
}

func (s *OLAPServer) newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", s.uploadMetrics)
	mux.HandleFunc("/last_screen_time", s.lastScreenOpenTimeMetrics)
	mux.HandleFunc("/last_request_time", s.lastRequestTimeMetrics)
	return mux
}

type RawMetric struct {
	Type       string        `json:"type"`
	SessionID  string        `json:"session_id"`
	Platform   string        `json:"platform"`
	Timestamp  time.Time     `json:"timestamp"`
	Duration   time.Duration `json:"duration"`
	ScreenName string        `json:"screen_name"`
	RequestURL string        `json:"request_url"`
}

func (s *OLAPServer) uploadMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Error.Printf(
			"Metrics upload: wrong method type: %s",
			r.Method,
		)
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var rBody struct {
		Metrics []RawMetric `json:"metrics"`
	}
	err := json.NewDecoder(r.Body).Decode(&rBody)
	if err != nil {
		s.logger.Error.Printf("Malformed upload metric request: %s\n", err)
		http.Error(w, "Request body malformed.", http.StatusBadRequest)
		return
	}

	var screenMetrics []core.ScreenOpeningTimeMetric
	var requestMetrics []core.RequestTimeMetric
	for _, m := range rBody.Metrics {
		switch m.Type {
		case "screen_open_time_metric":
			screenMetrics = append(screenMetrics, core.ScreenOpeningTimeMetric{
				DurationMetricBase: core.DurationMetricBase{
					MetricBase: core.MetricBase{
						SessionID: m.SessionID,
						Platform:  m.Platform,
						Timestamp: m.Timestamp,
					},
					Duration: m.Duration,
				},
				ScreenName: m.ScreenName,
			})
		case "request_time_metric":
			requestMetrics = append(requestMetrics, core.RequestTimeMetric{
				DurationMetricBase: core.DurationMetricBase{
					MetricBase: core.MetricBase{
						SessionID: m.SessionID,
						Platform:  m.Platform,
						Timestamp: m.Timestamp,
					},
					Duration: m.Duration,
				},
				RequestURL: m.RequestURL,
			})
		default:
			s.logger.Error.Printf("Unknown metric type: %s\n", m.Type)
		}
	}
	if len(screenMetrics) > 0 {
		s.logger.Info.Println("Adding screen opening time metrics.")
		s.repo.AddScreenOpeningTimeMetrics(screenMetrics)
	}
	if len(requestMetrics) > 0 {
		s.logger.Info.Println("Adding request time metrics.")
		s.repo.AddRequestTimeMetrics(requestMetrics)
	}
}

func (s *OLAPServer) lastScreenOpenTimeMetrics(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	metrics := s.repo.GetLastScreenOpeningTimeMetrics(10)
	raw := make([]RawMetric, len(metrics))
	for i, m := range metrics {
		raw[i] = RawMetric{
			Type:       "screen_open_time_metric",
			SessionID:  m.SessionID,
			Platform:   m.Platform,
			Timestamp:  m.Timestamp,
			Duration:   m.Duration,
			ScreenName: m.ScreenName,
			RequestURL: "",
		}
	}
	err := json.NewEncoder(w).Encode(struct {
		Metrics []RawMetric `json:"metrics"`
	}{raw})
	if err != nil {
		s.logger.Error.Printf("Screen time metrics marshaling error: %s\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *OLAPServer) lastRequestTimeMetrics(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	metrics := s.repo.GetLastAddRequestTimeMetrics(10)
	raw := make([]RawMetric, len(metrics))
	for i, m := range metrics {
		raw[i] = RawMetric{
			Type:       "request_time_metric",
			SessionID:  m.SessionID,
			Platform:   m.Platform,
			Timestamp:  m.Timestamp,
			Duration:   m.Duration,
			ScreenName: "",
			RequestURL: m.RequestURL,
		}
	}
	err := json.NewEncoder(w).Encode(struct {
		Metrics []RawMetric `json:"metrics"`
	}{raw})
	if err != nil {
		s.logger.Error.Printf(
			"Screen time metrics marshaling error: %s\n",
			err,
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
