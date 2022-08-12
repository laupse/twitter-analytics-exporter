package services

import (
	"reflect"
	"time"

	"github.com/laupse/twitter-analytics-exporter/ports"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type MetricsService struct {
	metrics *prometheus.GaugeVec
	repo    ports.AnalyticsRepository
}

func NewMetricsService(ar ports.AnalyticsRepository) *MetricsService {
	opts := prometheus.GaugeOpts{
		Name: "tweet_analytics",
	}
	metricsGauge := prometheus.NewGaugeVec(opts, []string{"id", "conversation_id", "type"})
	prometheus.MustRegister(metricsGauge)
	return &MetricsService{
		metrics: metricsGauge,
		repo:    ar,
	}

}

func (ms *MetricsService) Collect(userId string, refreshInterval time.Duration) {
	for {
		tweets, err := ms.repo.GetAnalytics(userId)
		if err != nil {
			log.Errorf("%s", err)
			time.Sleep(refreshInterval)
			continue
		}
		for _, tweet := range tweets {
			n := reflect.TypeOf(tweet.OrganicMetrics)
			v := reflect.ValueOf(tweet.OrganicMetrics)

			for i := 0; i < n.NumField(); i++ {
				metricValue := v.Field(i).Interface().(float64)
				gauge, err := ms.metrics.GetMetricWithLabelValues(tweet.Id, tweet.ConversationId, n.Field(i).Tag.Get("json"))
				if err != nil {
					log.Errorf("%s", err)
					continue
				}

				gauge.Set(metricValue)
			}
		}
		time.Sleep(refreshInterval)
	}
}
