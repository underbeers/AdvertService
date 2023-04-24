package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	contentType = "Content-Type"
	appJSON     = "application/json"
)

func (h *Handler) handleInfo(c *gin.Context) {
	c.Writer.Header().Add(contentType, appJSON)
	serviceInfo := GetServiceInfo(h.services.Config, h.services.Router)
	payload, err := json.Marshal(serviceInfo)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	_, err = c.Writer.Write(payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

func GetServiceInfo(config *repository.Config, e *gin.Engine) *repository.Service {
	handles, err := getHandles(e)
	if err != nil {
		logrus.Fatalf("failed to getHandles, %v", err)
	}

	instance := repository.Service{
		Name:      "advert",
		Label:     "pl_advert_service",
		IP:        config.Listen.IP,
		Port:      config.Listen.Port,
		Endpoints: nil,
	}
	unprotected, err := getUnprotected()
	if err != nil {
		logrus.Fatalf("failed to getHandles, %v", err)
	}

	for k, v := range handles {
		// skip endpoint-info
		if k == "endpoint-info/" {
			continue
		}
		endpoint := models.Endpoint{
			URL:       k,
			Protected: true,
			Methods:   v,
		}
		if unprotected[k] {
			endpoint.Protected = false
		}
		instance.Endpoints = append(instance.Endpoints, endpoint)
	}

	return &instance
}

func getHandles(e *gin.Engine) (map[string][]string, error) {
	data := make(map[string][]string)
	handlers := e.Routes()
	logrus.Info(handlers)

	for _, v := range handlers {
		path := v.Path
		method := []string{v.Method}
		path = strings.Split(path, "/api/v1/")[1]
		data[path] = append(data[path], method...)
	}

	return data, nil
}

func getUnprotected() (map[string]bool, error) {
	// Read's list of unprotected endpoints
	lst, err := os.OpenFile("service.json", os.O_RDONLY, 0o600) //nolint:gomnd
	if err != nil {
		return nil, errors.New("can't open file")
	}
	reader, err := io.ReadAll(lst)
	if err != nil {
		return nil, errors.New("can't read file")
	}
	data := struct {
		URLS []string `json:"urls"`
	}{}
	err = json.Unmarshal(reader, &data)
	if err != nil {
		return nil, errors.New("can't unmarshal data")
	}
	result := make(map[string]bool)
	for _, k := range data.URLS {
		result[k] = true
	}

	return result, nil
}
