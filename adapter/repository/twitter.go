package repository

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/gomodule/oauth1/oauth"
	"github.com/laupse/twitter-analytics-exporter/application/entity"
	log "github.com/sirupsen/logrus"
)

type TwitterRepository struct {
	restyClient          *resty.Client
	OAuthClient          *oauth.Client
	OAuthUserCredentials *oauth.Credentials
}

type TwitterTimelineResponse struct {
	Data []entity.Tweet `json:"data"`
}

func NewTwitterRepository(conusumerKey, consumerSecret, accessToken, tokenSecret string) *TwitterRepository {
	restyClient := resty.New()
	restyClient.SetAuthScheme("OAuth")
	restyClient.SetBaseURL("https://api.twitter.com/2")

	oauthClient := &oauth.Client{
		SignatureMethod: oauth.HMACSHA1,
		Credentials: oauth.Credentials{
			Token:  conusumerKey,
			Secret: consumerSecret,
		},
	}
	return &TwitterRepository{
		restyClient: restyClient,
		OAuthClient: oauthClient,
		OAuthUserCredentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: tokenSecret,
		},
	}
}

func (t *TwitterRepository) GetAnalytics(userId string) ([]entity.Tweet, error) {
	path := fmt.Sprintf("/users/%s/tweets", userId)

	queryParams := map[string]string{
		"max_results":  "100",
		"tweet.fields": "organic_metrics,conversation_id",
		"exclude":      "replies",
	}

	values := &url.Values{}
	for k, v := range queryParams {
		values.Add(k, v)
	}

	req := t.restyClient.R().SetQueryParamsFromValues(*values)
	fullUrl := t.restyClient.BaseURL + path
	t.OAuthClient.SignForm(t.OAuthUserCredentials, "GET", fullUrl, *values)

	for k, v := range queryParams {
		values.Add(k, v)
	}

	oauthHeader := ""
	for k, v := range *values {
		oauthHeader = oauthHeader + fmt.Sprintf("%s=\"%s\",", k, url.QueryEscape(v[0]))
	}
	oauthHeader = strings.TrimRight(oauthHeader, ",")

	resp, err := req.SetAuthToken(oauthHeader).Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		log.WithFields(log.Fields{
			"content": string(resp.Body()),
		}).Error("failed to retrieve timeline")
	}
	twitterResponse := TwitterTimelineResponse{}
	json.Unmarshal(resp.Body(), &twitterResponse)

	return twitterResponse.Data, nil
}
