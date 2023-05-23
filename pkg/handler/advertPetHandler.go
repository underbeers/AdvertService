package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/underbeers/AdvertService/pkg/models"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	moderationFailed = "Не прошло модерацию"
	archived         = "В архиве"
	published        = "Опубликовано"
	pageSize         = 10
	defaultPage      = 1
)

func descriptionFilter(s string) string {
	s = strings.ToLower(s)
	words := strings.Fields(s)
	for i := 0; i < len(words); i++ {
		for j := 0; j < len(banWords); j++ {
			if words[i] == banWords[j] {
				return moderationFailed
			}
		}
	}
	return published
}

func (h *Handler) createNewAdvert(c *gin.Context) {

	var input models.AdvertPet
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
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

	/*Проверка на то, что еще нет объявлений с такой карточкой*/
	filter := models.AdvertPetFilter{}
	filter.PetCardId = input.PetCardId
	_, total, err := h.services.AdvertPet.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if total == 1 {
		newErrorResponse(c, http.StatusBadRequest, "it is forbidden to create two ads for one pet card")
		return
	}

	input.Status = descriptionFilter(input.Description)

	/*Проверка на то, что район относится к нужному городу*/
	district, err := h.services.Location.GetDistricts(models.DistrictFilter{DistrictId: input.DistrictId})
	cityId := input.CityId
	if err != nil || len(district) == 0 || district[0].CityId != cityId {
		c.JSON(http.StatusBadRequest, statusResponse{"incorrect city or district id"})
		return
	}

	err = h.services.AdvertPet.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) getAllAdverts(c *gin.Context) {

	type AdvertsResponse struct {
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

	type Response struct {
		NextPage        string            `json:"nextPage"`
		TotalPage       int64             `json:"totalPage"`
		TotalCount      int64             `json:"totalCount"`
		AdvertsResponse []AdvertsResponse `json:"records"`
	}

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

	if query.Has("cityID") {
		cityID, err := strconv.Atoi(query.Get("cityID"))
		if err != nil || cityID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid city id param")
			return
		}
		filter.CityId = cityID
	}

	if query.Has("districtID") {
		districtID, err := strconv.Atoi(query.Get("districtID"))
		if err != nil || districtID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid district id param")
			return
		}
		filter.DistrictId = districtID
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

	if query.Has("page") {
		page, err := strconv.Atoi(query.Get("page"))
		if err != nil || page <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid page number")
			return
		}
		filter.Page = page
	}

	if query.Has("perPage") {
		perPage, err := strconv.Atoi(query.Get("perPage"))
		if err != nil || perPage <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "incorrect number of elements per page")
			return
		}
		filter.PerPage = perPage
	}

	if query.Has("petCardID") {
		petCardID, err := strconv.Atoi(query.Get("petCardID"))
		if err != nil || petCardID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid pet card id param")
			return
		}
		filter.PetCardId = petCardID
	}

	if query.Has("petTypeID") {
		PetTypeId, err := strconv.Atoi(query.Get("petTypeID"))
		if err != nil || PetTypeId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid pet type id param")
			return
		}
		filter.PetTypeId = PetTypeId
	}

	if query.Has("breedID") {
		BreedId, err := strconv.Atoi(query.Get("breedID"))
		if err != nil || BreedId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid breed id param")
			return
		}
		filter.BreedId = BreedId
	}

	if query.Has("gender") {
		gender := query.Get("gender")
		if gender == "male" {
			filter.Gender = "male"
		} else if gender == "female" {
			filter.Gender = "female"
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "incorrect gender format")
			return
		}
	}

	if filter.Page == 0 {
		filter.Page = defaultPage
	}

	if filter.PerPage == 0 {
		filter.PerPage = pageSize
	}

	advertPetList, total, err := h.services.AdvertPet.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var resp Response

	if len(advertPetList) == 0 {
		resp.AdvertsResponse = make([]AdvertsResponse, 0)
	}

	resp.TotalCount = total
	resp.TotalPage = int64(math.Ceil(float64(total) / float64(filter.PerPage)))

	query.Set("page", strconv.Itoa(filter.Page+1))
	if int64(filter.Page) < resp.TotalPage {
		resp.NextPage = "/api/v1/adverts?" + query.Encode()
	}

	for i := 0; i < len(advertPetList); i++ {
		resp.AdvertsResponse = append(resp.AdvertsResponse,
			AdvertsResponse{
				Id:          advertPetList[i].Id,
				PetCardId:   advertPetList[i].PetCardId,
				UserId:      advertPetList[i].UserId,
				PetName:     advertPetList[i].PetName,
				MainPhoto:   strings.Split(advertPetList[i].MainPhoto[1:len(advertPetList[i].MainPhoto)-1], ",")[0],
				Price:       advertPetList[i].Price,
				Description: advertPetList[i].Description,
				City:        advertPetList[i].City,
				District:    advertPetList[i].District,
				Publication: advertPetList[i].Publication,
				Gender:      advertPetList[i].Gender,
				BirthDate:   advertPetList[i].BirthDate,
				PetTypeName: advertPetList[i].PetTypeName,
				BreedName:   advertPetList[i].BreedName,
			})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getAuthAllAdverts(c *gin.Context) {

	type AdvertsResponse struct {
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
		InFavorites bool      `json:"inFavorites"`
		FavoritesId int       `json:"favoritesID"`
	}

	type Response struct {
		NextPage        string            `json:"nextPage"`
		TotalPage       int64             `json:"totalPage"`
		TotalCount      int64             `json:"totalCount"`
		AdvertsResponse []AdvertsResponse `json:"records"`
	}

	query := c.Request.URL.Query()
	filter := models.AdvertPetFilter{}

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

	if query.Has("cityID") {
		cityID, err := strconv.Atoi(query.Get("cityID"))
		if err != nil || cityID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid city id param")
			return
		}
		filter.CityId = cityID
	}

	if query.Has("districtID") {
		districtID, err := strconv.Atoi(query.Get("districtID"))
		if err != nil || districtID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid district id param")
			return
		}
		filter.DistrictId = districtID
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

	if query.Has("page") {
		page, err := strconv.Atoi(query.Get("page"))
		if err != nil || page <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid page number")
			return
		}
		filter.Page = page
	}

	if query.Has("perPage") {
		perPage, err := strconv.Atoi(query.Get("perPage"))
		if err != nil || perPage <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "incorrect number of elements per page")
			return
		}
		filter.PerPage = perPage
	}

	if query.Has("petCardID") {
		petCardID, err := strconv.Atoi(query.Get("petCardID"))
		if err != nil || petCardID <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid pet card id param")
			return
		}
		filter.PetCardId = petCardID
	}

	if query.Has("petTypeID") {
		PetTypeId, err := strconv.Atoi(query.Get("petTypeID"))
		if err != nil || PetTypeId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid pet type id param")
			return
		}
		filter.PetTypeId = PetTypeId
	}

	if query.Has("breedID") {
		BreedId, err := strconv.Atoi(query.Get("breedID"))
		if err != nil || BreedId <= 0 {
			newErrorResponse(c, http.StatusBadRequest, "invalid breed id param")
			return
		}
		filter.BreedId = BreedId
	}

	if query.Has("gender") {
		gender := query.Get("gender")
		if gender == "male" {
			filter.Gender = "male"
		} else if gender == "female" {
			filter.Gender = "female"
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "incorrect gender format")
			return
		}
	}

	if filter.Page == 0 {
		filter.Page = defaultPage
	}

	if filter.PerPage == 0 {
		filter.PerPage = pageSize
	}

	advertPetList, total, err := h.services.AdvertPet.GetAll(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	favoritesFilter := models.FavoritesFilter{UserId: parseUserID}
	allFavorites, err := h.services.Favorites.GetFavorites(favoritesFilter)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var resp Response

	if len(advertPetList) == 0 {
		resp.AdvertsResponse = make([]AdvertsResponse, 0)
	}

	resp.TotalCount = total
	resp.TotalPage = int64(math.Ceil(float64(total) / float64(filter.PerPage)))

	query.Set("page", strconv.Itoa(filter.Page+1))
	if int64(filter.Page) < resp.TotalPage {
		resp.NextPage = "/api/v1/adverts?" + query.Encode()
	}

	for i := 0; i < len(advertPetList); i++ {
		var favoritesID int
		inFavorites := false
		for j := 0; j < len(allFavorites); j++ {
			if allFavorites[j].AdvertId == advertPetList[i].Id {
				favoritesID = allFavorites[j].Id
				inFavorites = true
				break
			}
		}

		resp.AdvertsResponse = append(resp.AdvertsResponse,
			AdvertsResponse{
				Id:          advertPetList[i].Id,
				PetCardId:   advertPetList[i].PetCardId,
				UserId:      advertPetList[i].UserId,
				PetName:     advertPetList[i].PetName,
				MainPhoto:   strings.Split(advertPetList[i].MainPhoto[1:len(advertPetList[i].MainPhoto)-1], ",")[0],
				Price:       advertPetList[i].Price,
				Description: advertPetList[i].Description,
				City:        advertPetList[i].City,
				District:    advertPetList[i].District,
				Publication: advertPetList[i].Publication,
				Gender:      advertPetList[i].Gender,
				BirthDate:   advertPetList[i].BirthDate,
				PetTypeName: advertPetList[i].PetTypeName,
				BreedName:   advertPetList[i].BreedName,
				InFavorites: inFavorites,
				FavoritesId: favoritesID,
			})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getFullAdvert(c *gin.Context) {

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

	type PhotoResponse struct {
		ThumbnailPhoto string `json:"thumbnail"`
		Photo          string `json:"original"`
	}

	type AdvertResponse struct {
		Id            int             `json:"id" db:"id"`
		PetCardId     int             `json:"petCardID"`
		UserId        uuid.UUID       `json:"userID" db:"user_id"`
		Price         int             `json:"price" db:"price"`
		Description   string          `json:"description" db:"description"`
		City          string          `json:"city"`
		District      string          `json:"district"`
		Chat          bool            `json:"chat" db:"chat"`
		Phone         string          `json:"phone" db:"phone"`
		Status        string          `json:"status" db:"status"`
		Publication   time.Time       `json:"publication" db:"publication"`
		PetName       string          `json:"petName" db:"pet_name"`
		PetTypeName   string          `json:"petType" db:"pet_type"`
		BreedName     string          `json:"breed" db:"breed_name"`
		Photo         []PhotoResponse `json:"photos"`
		BirthDate     time.Time       `json:"birthDate" db:"birth_date"`
		Gender        string          `json:"gender" db:"gender"`
		Color         string          `json:"color" db:"color"`
		Care          string          `json:"care" db:"care"`
		Character     string          `json:"petCharacter" db:"pet_character"`
		Pedigree      string          `json:"pedigree" db:"pedigree"`
		Sterilization bool            `json:"sterilization" db:"sterilization"`
		Vaccinations  bool            `json:"vaccinations" db:"vaccinations"`
	}

	advert, err := h.services.AdvertPet.GetFullAdvert(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
		return
	}

	var photos []PhotoResponse
	if len(advert.Photo) > 5 && len(advert.ThumbnailPhoto) > 5 {
		originalPhoto := strings.Split(advert.Photo[1:len(advert.Photo)-1], ",")
		thumbnailPhoto := strings.Split(advert.ThumbnailPhoto[1:len(advert.ThumbnailPhoto)-1], ",")
		for j := 0; j < len(originalPhoto) || j < len(thumbnailPhoto); j++ {
			photos = append(photos, PhotoResponse{})
			if j < len(thumbnailPhoto) {
				photos[j].ThumbnailPhoto = strings.TrimSpace(thumbnailPhoto[j])
			}
			if j < len(originalPhoto) {
				photos[j].Photo = strings.TrimSpace(originalPhoto[j])
			}
		}
	}
	resp := AdvertResponse{
		Id:            advert.Id,
		PetCardId:     advert.PetCardId,
		UserId:        advert.UserId,
		Price:         advert.Price,
		Description:   advert.Description,
		City:          advert.City,
		District:      advert.District,
		Chat:          advert.Chat,
		Phone:         advert.Phone,
		Status:        advert.Status,
		Publication:   advert.Publication,
		PetTypeName:   advert.PetTypeName,
		PetName:       advert.PetName,
		BreedName:     advert.BreedName,
		Photo:         photos,
		BirthDate:     advert.BirthDate,
		Gender:        advert.Gender,
		Color:         advert.Color,
		Care:          advert.Care,
		Character:     advert.Character,
		Pedigree:      advert.Pedigree,
		Sterilization: advert.Sterilization,
		Vaccinations:  advert.Vaccinations,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getAuthFullAdvert(c *gin.Context) {

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

	type PhotoResponse struct {
		ThumbnailPhoto string `json:"thumbnail"`
		Photo          string `json:"original"`
	}

	type AdvertResponse struct {
		Id            int             `json:"id" db:"id"`
		PetCardId     int             `json:"petCardID"`
		UserId        uuid.UUID       `json:"userID" db:"user_id"`
		Price         int             `json:"price" db:"price"`
		Description   string          `json:"description" db:"description"`
		City          string          `json:"city"`
		District      string          `json:"district"`
		Chat          bool            `json:"chat" db:"chat"`
		Phone         string          `json:"phone" db:"phone"`
		Status        string          `json:"status" db:"status"`
		Publication   time.Time       `json:"publication" db:"publication"`
		PetName       string          `json:"petName" db:"pet_name"`
		PetTypeName   string          `json:"petType" db:"pet_type"`
		BreedName     string          `json:"breed" db:"breed_name"`
		Photo         []PhotoResponse `json:"photos"`
		BirthDate     time.Time       `json:"birthDate" db:"birth_date"`
		Gender        string          `json:"gender" db:"gender"`
		Color         string          `json:"color" db:"color"`
		Care          string          `json:"care" db:"care"`
		Character     string          `json:"petCharacter" db:"pet_character"`
		Pedigree      string          `json:"pedigree" db:"pedigree"`
		Sterilization bool            `json:"sterilization" db:"sterilization"`
		Vaccinations  bool            `json:"vaccinations" db:"vaccinations"`
		InFavorites   bool            `json:"inFavorites"`
		FavoritesId   int             `json:"favoritesID"`
	}

	advert, err := h.services.AdvertPet.GetFullAdvert(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
		return
	}

	favoritesFilter := models.FavoritesFilter{UserId: parseUserID}
	allFavorites, err := h.services.Favorites.GetFavorites(favoritesFilter)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var photos []PhotoResponse
	if len(advert.Photo) > 5 && len(advert.ThumbnailPhoto) > 5 {
		originalPhoto := strings.Split(advert.Photo[1:len(advert.Photo)-1], ",")
		thumbnailPhoto := strings.Split(advert.ThumbnailPhoto[1:len(advert.ThumbnailPhoto)-1], ",")
		for j := 0; j < len(originalPhoto) || j < len(thumbnailPhoto); j++ {
			photos = append(photos, PhotoResponse{})
			if j < len(thumbnailPhoto) {
				photos[j].ThumbnailPhoto = strings.TrimSpace(thumbnailPhoto[j])
			}
			if j < len(originalPhoto) {
				photos[j].Photo = strings.TrimSpace(originalPhoto[j])
			}
		}
	}

	var favoritesID int
	inFavorites := false
	for j := 0; j < len(allFavorites); j++ {
		if allFavorites[j].AdvertId == advert.Id {
			favoritesID = allFavorites[j].Id
			inFavorites = true
			break
		}
	}

	resp := AdvertResponse{
		Id:            advert.Id,
		PetCardId:     advert.PetCardId,
		UserId:        advert.UserId,
		Price:         advert.Price,
		Description:   advert.Description,
		City:          advert.City,
		District:      advert.District,
		Chat:          advert.Chat,
		Phone:         advert.Phone,
		Status:        advert.Status,
		Publication:   advert.Publication,
		PetTypeName:   advert.PetTypeName,
		PetName:       advert.PetName,
		BreedName:     advert.BreedName,
		Photo:         photos,
		BirthDate:     advert.BirthDate,
		Gender:        advert.Gender,
		Color:         advert.Color,
		Care:          advert.Care,
		Character:     advert.Character,
		Pedigree:      advert.Pedigree,
		Sterilization: advert.Sterilization,
		Vaccinations:  advert.Vaccinations,
		InFavorites:   inFavorites,
		FavoritesId:   favoritesID,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) changeStatus(c *gin.Context) {

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

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetFullAdvert(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
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

	/*Проверка на то, что id из токена совпадает с id владельца объявления*/
	if advertPet.UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to change status")
		return
	}

	status := ""

	if advertPet.Status == archived {
		status = published
	} else if advertPet.Status == published {
		status = archived
	} else {
		newErrorResponse(c, http.StatusBadRequest, "can't change status because ad moderation failed")
		return
	}

	if err := h.services.AdvertPet.ChangeStatus(id, status); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) updateAdvert(c *gin.Context) {
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

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetFullAdvert(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
		return
	}

	var input models.UpdateAdvertInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
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

	/*Проверка на то, что id из токена совпадает с id владельца объявления*/
	if advertPet.UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to update")
		return
	}

	if input.Description != nil {
		getStatus := descriptionFilter(*input.Description)
		input.Status = &getStatus
	}

	/*Проверка на то, что район относится к нужному городу*/
	if input.DistrictId != nil {
		district, err := h.services.Location.GetDistricts(models.DistrictFilter{DistrictId: *input.DistrictId})
		cityId := advertPet.CityId
		if input.CityId != nil {
			cityId = *input.CityId
		}
		if err != nil || len(district) == 0 || district[0].CityId != cityId {
			c.JSON(http.StatusBadRequest, statusResponse{"incorrect city or district id"})
			return
		}
	}

	if err := h.services.AdvertPet.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

func (h *Handler) deleteAdvert(c *gin.Context) {
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

	/*Проверка, что такой advert pet id существует*/
	advertPet, err := h.services.AdvertPet.GetFullAdvert(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "incorrect advert pet id")
		return
	}

	/*Проверка на то, что id из токена совпадает с id владельца объявления*/
	if advertPet.UserId != parseUserID {
		newErrorResponse(c, http.StatusBadRequest, "not enough permissions to delete")
		return
	}

	err = h.services.AdvertPet.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}
