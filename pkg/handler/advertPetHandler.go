package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/underbeers/AdvertService/pkg/models"
	"net/http"
	"strconv"
)

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

	input.Status = "На модерации"

	err = h.services.AdvertPet.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) getAllAdverts(c *gin.Context) {

	filter := models.AdvertPetFilter{}

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

func (h *Handler) updateAdvert(c *gin.Context) {

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
