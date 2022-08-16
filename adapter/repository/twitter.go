package repository

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/gomodule/oauth1/oauth"
	"github.com/laupse/twitter-analytics-exporter/application/entity"
)

type TwitterRepository struct {
	RestyClient          *resty.Client
	OAuthClient          *oauth.Client
	OAuthUserCredentials *oauth.Credentials
}

type TwitterTimelineResponse struct {
	Data []entity.Tweet `json:"data"`
}

func NewTwitterRepository(conusumerKey, consumerSecret, accessToken, tokenSecret string) *TwitterRepository {
	RestyClient := resty.New()
	RestyClient.SetAuthScheme("OAuth")
	RestyClient.SetBaseURL("https://api.twitter.com/2")
	RestyClient.SetQueryParams(map[string]string{
		"max_results":  "100",
		"tweet.fields": "organic_metrics,conversation_id",
		"exclude":      "replies",
	})

	oauthClient := &oauth.Client{
		SignatureMethod: oauth.HMACSHA1,
		Credentials: oauth.Credentials{
			Token:  conusumerKey,
			Secret: consumerSecret,
		},
	}

	return &TwitterRepository{
		RestyClient: RestyClient,
		OAuthClient: oauthClient,
		OAuthUserCredentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: tokenSecret,
		},
	}
}

func (t *TwitterRepository) GetAnalytics(userId string) ([]entity.Tweet, error) {
	path := fmt.Sprintf("/users/%s/tweets", userId)
	fullUrl, err := url.Parse(t.RestyClient.BaseURL + path)
	if err != nil {
		return nil, err
	}

	req := t.RestyClient.R()
	if err := t.OAuthClient.SetAuthorizationHeader(req.Header, t.OAuthUserCredentials, "GET", fullUrl, t.RestyClient.QueryParam); err != nil {
		return nil, err
	}

	resp, err := req.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("http code is not 200")
	}

	twitterResponse := TwitterTimelineResponse{}
	if err := json.Unmarshal(resp.Body(), &twitterResponse); err != nil {
		return nil, err
	}

	return twitterResponse.Data, nil
}
