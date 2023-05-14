package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/underbeers/AdvertService/pkg/models"
	"net/http"
	"strconv"
)

func (h *Handler) getCities(c *gin.Context) {
	query := c.Request.URL.Query()
	filter := models.CityFilter{}

	if query.Has("id") {
		CityId, err := strconv.Atoi(query.Get("id"))
		if err != nil || CityId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid city id")
			return
		}
		filter.CityId = CityId
	}

	cityList, err := h.services.Location.GetCities(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(cityList) == 0 {
		newErrorResponse(c, http.StatusOK, "records not found")
		return
	}

	c.JSON(http.StatusOK, cityList)
}

func (h *Handler) getDistricts(c *gin.Context) {
	query := c.Request.URL.Query()
	filter := models.DistrictFilter{}

	if query.Has("id") {
		DistrictId, err := strconv.Atoi(query.Get("id"))
		if err != nil || DistrictId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid district id")
			return
		}
		filter.DistrictId = DistrictId
	}

	if query.Has("cityID") {
		cityID, err := strconv.Atoi(query.Get("cityID"))
		if err != nil || cityID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid city id")
			return
		}
		filter.CityId = cityID
	}

	cityList, err := h.services.Location.GetDistricts(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(cityList) == 0 {
		newErrorResponse(c, http.StatusOK, "records not found")
		return
	}

	c.JSON(http.StatusOK, cityList)
}
