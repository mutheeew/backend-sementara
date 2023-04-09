package handlers

import (
	filmdto "dumbflix/dto/film"
	dto "dumbflix/dto/result"
	"dumbflix/models"
	"dumbflix/repositories"
	"fmt"
	"os"

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var path_file = "http://localhost:5000/uploads/"

type handlerFilm struct {
	FilmRepository repositories.FilmRepository
}

func HandlerFilm(FilmRepository repositories.FilmRepository) *handlerFilm {
	return &handlerFilm{FilmRepository}
}

func (h *handlerFilm) FindFilms(c echo.Context) error {
	films, err := h.FilmRepository.FindFilms()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}
	for i, f := range films {
		films[i].Thumbnailfilm = path_file + f.Thumbnailfilm
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: films})
}

func (h *handlerFilm) GetFilm(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var film models.Film

	film, err := h.FilmRepository.GetFilm(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	film.Thumbnailfilm = path_file + film.Thumbnailfilm

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertResponseFilm(film)})
}

func (h *handlerFilm) CreateFilm(c echo.Context) error {
	dataFile := c.Get("dataFile").(string)
	fmt.Println("this is data file", dataFile)

	// userLogin := c.Get("userLogin")
	// userId := userLogin.(jwt.MapClaims)["id"].(float64)

	category_id, _ := strconv.Atoi(c.FormValue("category_id"))
	year, _ := strconv.Atoi(c.FormValue("year"))

	request := filmdto.FilmRequest{
		Title:         c.FormValue("title"),
		Year:          year,
		CategoryID:    category_id,
		Description:   c.FormValue("description"),
		Thumbnailfilm: dataFile,
		LinkFilm:      c.FormValue("linkfilm"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	film := models.Film{
		Title:         request.Title,
		Thumbnailfilm: request.Thumbnailfilm,
		Year:          request.Year,
		CategoryID:    request.CategoryID,
		LinkFilm:      request.LinkFilm,
		Description:   request.Description,
		ID:            int(userId),
	}

	film, err = h.FilmRepository.CreateFilm(film)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	film, _ = h.FilmRepository.GetFilm(film.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertResponseFilm(film)})
}

// func (h *handlerFilm) UpdateFilm(c echo.Context) error {
// 	dataFile := c.Get("dataFile").(string)
// 	fmt.Println("this is data file", dataFile)

// 	// userLogin := c.Get("userLogin")
// 	// userId := userLogin.(jwt.MapClaims)["id"].(float64)
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	year, _ := strconv.Atoi(c.FormValue("Year"))
// 	category_id, _ := strconv.Atoi(c.FormValue("Category_id"))

// 	request := filmdto.FilmRequest{
// 		Title:         c.FormValue("title"),
// 		Year:          year,
// 		CategoryID:    category_id,
// 		Description:   c.FormValue("description"),
// 		Thumbnailfilm: dataFile,
// 		LinkFilm:      c.FormValue("linkfilm"),
// 	}

// 	validation := validator.New()
// 	err := validation.Struct(request)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
// 	}

// 	film, _ := h.FilmRepository.GetFilm(id)

// 	film.ID = request.ID

// 	if request.Title != "" {
// 		film.Title = request.Title
// 	}
// 	if request.Thumbnailfilm != "" {
// 		film.Thumbnailfilm = request.Thumbnailfilm
// 	}
// 	if request.Year != 0 {
// 		film.Year = request.Year
// 	}
// 	if request.LinkFilm != "" {
// 		film.LinkFilm = request.LinkFilm
// 	}
// 	if request.Description != "" {
// 		film.Description = request.Description
// 	}

// 	data, err := h.FilmRepository.UpdateFilm(film)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Profile data updated successfully", Data: convertResponseFilm(data)})
// }

func (h *handlerFilm) DeleteFilm(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	film, err := h.FilmRepository.GetFilm(int(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	fileName := film.Thumbnailfilm
	dirPath := "uploads"

	filePath := fmt.Sprintf("%s/%s", dirPath, fileName)

	data, err := h.FilmRepository.DeleteFilm(film)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Failed to delete file"+fileName+":", err)
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	fmt.Println("File " + fileName + " deleted successfully")

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Profile data deleted successfully", Data: convertResponseFilm(data)})
}

func convertResponseFilm(u models.Film) models.FilmResponse {
	return models.FilmResponse{
		ID:            u.ID,
		Title:         u.Title,
		Thumbnailfilm: u.Thumbnailfilm,
		Category:      u.Category,
		// Year:          u.Year,
		LinkFilm:    u.LinkFilm,
		Description: u.Description,
	}
}
