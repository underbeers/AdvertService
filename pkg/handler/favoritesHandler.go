package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/underbeers/AdvertService/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) addFavorites(c *gin.Context) {

	var input models.Favorites
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	switch input.Type {
	case "advert":
		input.Type = "advert_id"
		input.AdvertId = input.Id
	case "organization":
		input.Type = "organization_id"
		input.AdvertId = input.Id
	case "specialist":
		input.Type = "specialist_id"
		input.AdvertId = input.Id
	case "event":
		input.Type = "event_id"
		input.AdvertId = input.Id
	default:
		newErrorResponse(c, http.StatusBadRequest, "incorrect type param")
		return
	}

	userID := c.Request.Header.Get("userID")

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

	if input.Type == "advert_id" {

		/*Проверка, что такой advert pet id существует*/
		advert, err := h.services.AdvertPet.GetFullAdvert(input.Id)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
			return
		}

	}

	allFavorites, err := h.services.Favorites.GetFavorites(models.FavoritesFilter{UserId: input.UserId})
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for i := 0; i < len(allFavorites); i++ {
		if input.Type == "advert_id" && allFavorites[i].AdvertId == input.Id {
			newErrorResponse(c, http.StatusBadRequest, "record has already been added")
			return
		} else if input.Type == "organization_id" && allFavorites[i].OrganizationId == input.Id {
			newErrorResponse(c, http.StatusBadRequest, "record has already been added")
			return
		} else if input.Type == "specialist_id" && allFavorites[i].SpecialistId == input.Id {
			newErrorResponse(c, http.StatusBadRequest, "record has already been added")
			return
		} else if input.Type == "event_id" && allFavorites[i].EventId == input.Id {
			newErrorResponse(c, http.StatusBadRequest, "record has already been added")
			return
		}
	}

	err = h.services.Favorites.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) getFavorites(c *gin.Context) {

	type FavoritesAdvertsResponse struct {
		FavoritesId int       `json:"favoritesID"`
		Id          int       `json:"id"`
		PetCardId   int       `json:"petCardID"`
		UserId      uuid.UUID `json:"userID"`
		PetName     string    `json:"petName"`
		MainPhoto   string    `json:"mainPhoto"`
		Price       int       `json:"price"`
		Description string    `json:"description"`
		City        string    `json:"city"`
		District    string    `json:"district"`
		Publication time.Time `json:"publication"`
		Gender      string    `json:"gender"`
		BirthDate   time.Time `json:"birthDate"`
		PetTypeName string    `json:"petType"`
		BreedName   string    `json:"breed"`
	}

	type FavoritesInfo struct {
		FavoritesId int `json:"favoritesID"`
		ServiceData int `json:"id"`
	}

	type FavoritesResponse struct {
		Adverts       []FavoritesAdvertsResponse `json:"adverts"`
		Organizations []FavoritesInfo            `json:"organizations"`
		Specialists   []FavoritesInfo            `json:"specialists"`
		Events        []FavoritesInfo            `json:"events"`
	}

	filter := models.FavoritesFilter{}

	userID := c.Request.Header.Get("userID")

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

	favoritesAdverts, err := h.services.Favorites.GetFavoritesAdverts(filter)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	allFavorites, err := h.services.Favorites.GetFavorites(filter)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var resp FavoritesResponse

	if len(favoritesAdverts) == 0 {
		resp.Adverts = make([]FavoritesAdvertsResponse, 0)
	}

	for i := 0; i < len(allFavorites); i++ {
		if allFavorites[i].OrganizationId != 0 {
			resp.Organizations = append(resp.Organizations, FavoritesInfo{FavoritesId: allFavorites[i].Id,
				ServiceData: allFavorites[i].OrganizationId})
		} else if allFavorites[i].EventId != 0 {
			resp.Events = append(resp.Events, FavoritesInfo{FavoritesId: allFavorites[i].Id,
				ServiceData: allFavorites[i].EventId})
		} else if allFavorites[i].SpecialistId != 0 {
			resp.Specialists = append(resp.Specialists, FavoritesInfo{FavoritesId: allFavorites[i].Id,
				ServiceData: allFavorites[i].SpecialistId})
		}
	}

	if len(resp.Organizations) == 0 {
		resp.Organizations = make([]FavoritesInfo, 0)
	}
	if len(resp.Events) == 0 {
		resp.Events = make([]FavoritesInfo, 0)
	}
	if len(resp.Specialists) == 0 {
		resp.Specialists = make([]FavoritesInfo, 0)
	}

	for i := 0; i < len(favoritesAdverts); i++ {
		resp.Adverts = append(resp.Adverts,
			FavoritesAdvertsResponse{
				FavoritesId: favoritesAdverts[i].FavoritesId,
				Id:          favoritesAdverts[i].Id,
				PetCardId:   favoritesAdverts[i].PetCardId,
				UserId:      favoritesAdverts[i].UserId,
				PetName:     favoritesAdverts[i].PetName,
				MainPhoto:   strings.Split(favoritesAdverts[i].MainPhoto[1:len(favoritesAdverts[i].MainPhoto)-1], ",")[0],
				Price:       favoritesAdverts[i].Price,
				Description: favoritesAdverts[i].Description,
				City:        favoritesAdverts[i].City,
				District:    favoritesAdverts[i].District,
				Publication: favoritesAdverts[i].Publication,
				Gender:      favoritesAdverts[i].Gender,
				BirthDate:   favoritesAdverts[i].BirthDate,
				PetTypeName: favoritesAdverts[i].PetTypeName,
				BreedName:   favoritesAdverts[i].BreedName,
			})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) deleteFavorites(c *gin.Context) {
	var id int
	query := c.Request.URL.Query()
	if query.Has("id") {
		cardID, err := strconv.Atoi(query.Get("id"))
		id = cardID
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid id param")
			return
		}
	} else {
		newErrorResponse(c, http.StatusBadRequest, "id not provided")
		return
	}

	userID := c.Request.Header.Get("userID")

	if len(userID) == 0 {
		c.JSON(http.StatusBadRequest, statusResponse{"invalid access token"})
		return
	}

	parseUserID, err := uuid.Parse(userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	filter := models.FavoritesFilter{FavoritesId: id}
	allFavorites, err := h.services.Favorites.GetFavorites(filter)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(allFavorites) != 1 {
		newErrorResponse(c, http.StatusBadRequest, "incorrect favorites id")
		return
	}
	if allFavorites[0].UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to delete")
		return
	}

	err = h.services.Favorites.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}
