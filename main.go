package main

import (
	"strings"

	"github.com/laupse/twitter-analytics-exporter/adapter/http"
	"github.com/laupse/twitter-analytics-exporter/adapter/repository"
	"github.com/laupse/twitter-analytics-exporter/application/services"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("consumer-key", "", "key used to represents your Twitter developer app")
	pflag.String("consumer-secret", "", "secret used to represents your Twitter developer app")
	pflag.String("user-access-token", "", "user-specific token credentials used to authenticate OAuth 1.0a")
	pflag.String("user-secret-token", "", "user-specific secret credentials used to authenticate OAuth 1.0a")
	pflag.String("user-id", "", "user-specific id from where to retrieve timeline")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("TAE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	twitterRepository := repository.NewTwitterRepository(viper.GetString("consumer-key"), viper.GetString("consumer-secret"), viper.GetString("user-access-token"), viper.GetString("user-secret-token"))
	metricsService := services.NewMetricsService(twitterRepository)

	go metricsService.Collect(viper.GetString("user-id"))

	ginHandler := http.NewGinHandler()
	ginHandler.SetupRoutes()
	ginHandler.Run(":3000")

}
