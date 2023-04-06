package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/underbeers/AdvertService/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{

		advertPet := api.Group("adverts")
		{
			advertPet.POST("/new", h.createNewAdvert).OPTIONS("/new", h.createNewAdvert)
			advertPet.GET("", h.getAllAdverts).OPTIONS("", h.getAllAdverts)
			advertPet.POST("/changeStatus/:id", h.changeStatus).OPTIONS("/changeStatus/:id", h.changeStatus)
			advertPet.PUT("/update/:id", h.updateAdvert).OPTIONS("/update/:id", h.updateAdvert)
			advertPet.DELETE("/delete/:id", h.deleteAdvert).OPTIONS("/delete/:id", h.deleteAdvert)
		}

	}

	h.services.Router = router

	return router
}