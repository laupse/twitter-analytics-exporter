package ports

import "github.com/laupse/twitter-analytics-exporter/application/entity"

type AnalyticsRepository interface {
	GetAnalytics(userId string) ([]entity.Tweet, error)
}
