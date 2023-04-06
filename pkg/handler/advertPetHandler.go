package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/underbeers/AdvertService/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

const (
	moderationFailed = "Не прошло модерацию"
	archived         = "В архиве"
	published        = "Опубликовано"
)

func descriptionFilter(s string) string {
	s = strings.ToLower(s)
	words := strings.Fields(s)
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(banWords); j++ {
			if words[i] == banWords[j] {
				return "Не прошло модерацию"
			}
		}
	}
	return "Опубликовано"
}

func (h *Handler) createNewAdvert(c *gin.Context) {

	var input models.AdvertPet
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//userID := c.Request.Header.Get("userID")
	userID := "5cd754f9-d1aa-4b58-abc9-4d106be4d475"

	if len(userID) == 0 {
		c.JSON(http.StatusBadRequest, statusResponse{"invalid access token"})
		return
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input.UserId = id

	input.Status = descriptionFilter(input.Description)

	err = h.services.AdvertPet.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) getAllAdverts(c *gin.Context) {
	query := c.Request.URL.Query()
	filter := models.AdvertPetFilter{}

	if query.Has("id") {
		AdvertPetId, err := strconv.Atoi(query.Get("id"))
		if err != nil || AdvertPetId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid advert pet id param")
			return
		}
		filter.AdvertPetId = AdvertPetId
	}

	if query.Has("userID") {
		userID := query.Get("userID")
		if len(userID) == 0 {
			c.JSON(http.StatusBadRequest, statusResponse{"invalid access token"})
			return
		}
		id, err := uuid.Parse(userID)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		filter.UserId = id
	}

	if query.Has("minPrice") {
		minPrice, err := strconv.Atoi(query.Get("minPrice"))
		if err != nil || minPrice <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid advert pet id param")
			return
		}
		filter.MinPrice = minPrice
	}

	if query.Has("maxPrice") {
		maxPrice, err := strconv.Atoi(query.Get("maxPrice"))
		if err != nil || maxPrice <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid advert pet id param")
			return
		}
		filter.MaxPrice = maxPrice
	}

	if query.Has("region") {
		filter.Region = query.Get("region")
	}

	if query.Has("locality") {
		filter.Locality = query.Get("locality")
	}

	if query.Has("status") {

		switch query.Get("status") {
		case "moderationFailed":
			filter.Status = moderationFailed
		case "archived":
			filter.Status = archived
		case "published":
			filter.Status = published
		default:
			filter.Status = ""
		}

	}

	if query.Has("sort") {

		switch query.Get("sort") {
		case "minPrice":
			filter.MinPriceSort = true
		case "maxPrice":
			filter.MaxPriceSort = true
		case "publication":
			filter.PublicationSort = true
		}

	}

	advertPetList, err := h.services.AdvertPet.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(advertPetList) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "records not found")
		return
	}

	c.JSON(http.StatusOK, advertPetList)
}

func (h *Handler) changeStatus(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetAll(models.AdvertPetFilter{AdvertPetId: id})
	if len(advertPet) != 1 || err != nil {
		c.JSON(http.StatusBadRequest, statusResponse{"incorrect advert pet id"})
		return
	}

	//userID := c.Request.Header.Get("userID")
	userID := "5cd754f9-d1aa-4b58-abc9-4d106be4d475"
	if len(userID) == 0 {
		c.JSON(http.StatusBadRequest, statusResponse{"invalid access token"})
		return
	}

	parseUserID, err := uuid.Parse(userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	/*Проверка на то, что id из токена совпадает с id владельца объявления*/
	if advertPet[0].UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to change status")
		return
	}

	status := ""

	if advertPet[0].Status == archived {
		status = published
	} else if advertPet[0].Status == published {
		status = archived
	} else {
		newErrorResponse(c, http.StatusBadRequest, "Can't change status because ad moderation failed")
		return
	}

	if err := h.services.AdvertPet.ChangeStatus(id, status); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) updateAdvert(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetAll(models.AdvertPetFilter{AdvertPetId: id})
	if len(advertPet) != 1 || err != nil {
		c.JSON(http.StatusBadRequest, statusResponse{"incorrect advert pet id"})
		return
	}

	var input models.UpdateAdvertInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//userID := c.Request.Header.Get("userID")
	userID := "5cd754f9-d1aa-4b58-abc9-4d106be4d475"
	if len(userID) == 0 {
		c.JSON(http.StatusBadRequest, statusResponse{"invalid access token"})
		return
	}

	parseUserID, err := uuid.Parse(userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	/*Проверка на то, что id из токена совпадает с id владельца объявления*/
	if advertPet[0].UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to update")
		return
	}

	input.UserId = &parseUserID

	if input.Description != nil {
		getStatus := descriptionFilter(*input.Description)
		input.Status = &getStatus
	}

	if err := h.services.AdvertPet.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) deleteAdvert(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetAll(models.AdvertPetFilter{AdvertPetId: id})
	if len(advertPet) != 1 || err != nil {
		c.JSON(http.StatusBadRequest, statusResponse{"incorrect advert pet id"})
		return
	}

	err = h.services.AdvertPet.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}
