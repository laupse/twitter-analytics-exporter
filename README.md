# Twitter-analytics-exporter

This repository contains a prometheus "exporter" that will collect analytics about tweets for a user timeline

Analytics gathered belong to [organic metrics](https://developer.twitter.com/en/docs/twitter-api/metrics) from twitter

# Usage

## Configuration

You need to have a [twitter developper account](https://developer.twitter.com/en/docs/platform-overview) in order to use this service

And then you will be able to configured with below flags :
``` bash
Usage of ./twitter-analytics-exporter:
      --consumer-key string        key used to represents your Twitter developer app
      --consumer-secret string     secret used to represents your Twitter developer app
      --user-access-token string   user-specific token credentials used to authenticate OAuth 1.0a
      --user-id string             user-specific id from where to retrieve timeline
      --user-secret-token string   user-specific secret credentials used to authenticate OAuth 1.0a
pflag: help requested
```
You can also use environment variable by prefixing them by `TAE_` followed by uppercase value and `-` replace by `_`. Ex: `TAE_USER_ID="1477951133402513413"`

`consumer` information are credentials related to you twitter app

`user` credentials are related to twitter user

> This application uses Oauth1.0a without 3-legged Oauth meaning there is no "in-app" way of exchanging user credentials (for now ?)

## Deployment 

Best way to deploy this app is using docker image
``` bash
ghcr.io/laupse/twitter-analytics-exporter
```

This application is an exporter so it needs to be backed by a :
* `prometheus` server (or a grafana-agent with long-term storage) in order to scrape metrics, store them and make them available by querying 
* `grafana` in order to visualize

# Limitations

* Using OAuth1.0a comes with limitation, 
    * Only tweets younger than 30 days are gathered because of the organic metrics requested
    * No easy way to configure scraping user that are not the user behind the app 

Overcome this will needs developping a full Oauth 2.0 app
