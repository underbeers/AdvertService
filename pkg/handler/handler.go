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
			advertPet.GET("/full", h.getFullAdvert).OPTIONS("/full", h.getFullAdvert)
			advertPet.POST("/changeStatus", h.changeStatus).OPTIONS("/changeStatus", h.changeStatus)
			advertPet.PUT("/update", h.updateAdvert).OPTIONS("/update", h.updateAdvert)
			advertPet.DELETE("/delete", h.deleteAdvert).OPTIONS("/delete", h.deleteAdvert)
		}

		authAdvertPet := api.Group("auth/adverts")
		{
			authAdvertPet.GET("", h.getAuthAllAdverts).OPTIONS("", h.getAuthAllAdverts)
			authAdvertPet.GET("/full", h.getAuthFullAdvert).OPTIONS("/full", h.getAuthFullAdvert)
		}

		locality := api.Group("location")
		{
			locality.GET("/city", h.getCities).OPTIONS("/city", h.getCities)
			locality.GET("/district", h.getDistricts).OPTIONS("/district", h.getDistricts)
		}

		favorites := api.Group("favorites")
		{
			favorites.POST("/add", h.addFavorites).OPTIONS("/add", h.addFavorites)
			favorites.GET("", h.getFavorites).OPTIONS("", h.getFavorites)
			favorites.DELETE("/delete", h.deleteFavorites).OPTIONS("/delete", h.deleteFavorites)
		}

		gwConnect := api.Group("endpoint-info")
		{
			gwConnect.GET("/", h.handleInfo).OPTIONS("/", h.handleInfo)
		}

	}

	h.services.Router = router

	return router
}
