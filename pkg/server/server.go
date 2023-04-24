package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	attempts = 4
	timeout  = 5
	protocol = "http"
	baseURL  = "/api/v1/"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler *gin.Engine, config *repository.Config) error {
	logrus.Infof("Start to listen to port %s", config.Listen.Port)
	var errorCnt int
	for errorCnt < attempts {
		time.Sleep(time.Second * timeout)
		logrus.Infof("Attempt to connect to APIGateway %d of %d", errorCnt+1, attempts)
		if err := pingAPIGateway(config); err != nil {
			errorCnt++
		} else {
			break
		}
	}
	if errorCnt >= attempts {
		return errors.New("failed to send info to the ApiGateway")
	}
	err := HelloAPIGateway(config)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{
		Addr:           ":" + config.Listen.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func knock(url string, payload []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload)) //nolint: gosec, noctx
	if resp == nil {
		// FIXME:Super dirty. Need to handle error
		log.Println("can't say Hello to Gateway", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		log.Println("knock() Post Error", err)
	}
	if resp.StatusCode == http.StatusOK {
		log.Println("Successfully greet ApiGateway")
	}
}

func pingAPIGateway(config *repository.Config) error {
	gwURL, err := gatewayURL(config)
	if err != nil {
		return errors.New("can't generate pingApiGateway gwURL")

	}
	resp, err := http.Get(gwURL.String()) //nolint: noctx

	if resp == nil {
		return errors.New("can't ping ApiGateway")
	}
	if err != nil {
		return errors.New("error pingApiGateway http.Get(gwURL)")
	}
	if err := resp.Body.Close(); err != nil {
		log.Println("can't close response body")
	}

	return nil
}

func gatewayURL(config *repository.Config) (*url.URL, error) {
	var domain string
	domain = config.Gateway.IP
	gwURL, err := url.Parse(
		protocol + "://" + domain + ":" + config.Gateway.Port + baseURL + "hello?service=advert&port=" + config.Listen.Port)
	if err != nil {
		return nil, errors.New("can't connect to ApiGateway")
	}
	logrus.Infof("Connection gateway url %s", gwURL.String())

	return gwURL, nil
}

func HelloAPIGateway(config *repository.Config) error {
	var domain string

	domain = config.Gateway.IP
	gatewayURL, err := url.Parse(
		protocol + "://" + domain + ":" + config.Gateway.Port + baseURL + "hello")
	if err != nil {
		return errors.New("can't parse url for endpoint 'hello'")
	}

	info := &models.Hello{
		Name:      "advert",
		Label:     "pl_advert_service",
		IP:        config.Listen.IP,
		Port:      config.Listen.Port,
		Endpoints: nil,
	}
	jsonStr, err := json.Marshal(info)
	if err != nil {
		return errors.New("could not marshal data")
	}

	go knock(gatewayURL.String(), jsonStr)

	return nil
}
