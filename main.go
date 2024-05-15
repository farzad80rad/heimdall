package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"heimdall/internal/config"
	"heimdall/internal/heimdall"
	"log"
)

func main() {

	apiConfigs, err := config.ReadConfig("/app/config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}
	r := gin.Default()
	for _, apiConfig := range apiConfigs.ApisConfig {
		log.Println(apiConfig)
		func() {
			defer func() {
				err := recover()
				if err != nil {
					log.Println(fmt.Sprintf("unable to start porxy for service: %s for reason :%v", apiConfig.Match.Name, err))
				}
			}()
			h, err := heimdall.NewApiGateway(apiConfig, r)
			if err != nil {
				log.Println(fmt.Sprintf("creating gateway for url %s failed for reason: %v", apiConfig.Match.Path, err))
				return
			}
			if err := h.Run(); err != nil {
				log.Println(fmt.Sprintf("running gateway for url %s failed for reason: %v", apiConfig.Match.Path, err))
			}
		}()
	}

	r.Run(":80")
}
