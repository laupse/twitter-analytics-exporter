package entity

type Tweet struct {
	Id             string  `json:"id"`
	ConversationId string  `json:"conversation_id"`
	Text           string  `json:"text"`
	OrganicMetrics Metrics `json:"organic_metrics"`
}

type Metrics struct {
	ImpressionCount  float64 `json:"impression_count"`
	UrlLinkClicks    float64 `json:"url_link_clicks"`
	ReplyCount       float64 `json:"reply_count"`
	RetweetCount     float64 `json:"retweet_count"`
	LikeCount        float64 `json:"like_count"`
	UserProfileClick float64 `json:"user_profile_clicks"`
}
