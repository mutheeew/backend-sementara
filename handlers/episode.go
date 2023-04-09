package handlers

import (
	episodedto "dumbflix/dto/episode"
	dto "dumbflix/dto/result"
	"dumbflix/models"
	"dumbflix/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerEpisode struct {
	EpisodeRepository repositories.EpisodeRepository
}

func HandlerEpisode(EpisodeRepository repositories.EpisodeRepository) *handlerEpisode {
	return &handlerEpisode{EpisodeRepository}
}

func (h *handlerEpisode) FindEpisodes(c echo.Context) error {
	episodes, err := h.EpisodeRepository.FindEpisode()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: episodes})
}

func (h *handlerEpisode) GetEpisode(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var episode models.Episode
	episode, err := h.EpisodeRepository.GetEpisode(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Episode successfully obtained", Data: convertResponseEpisode(episode)})
}

func (h *handlerEpisode) CreateEpisode(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userAdmin := userLogin.(jwt.MapClaims)["is_admin"].(bool)
	if userAdmin {
		request := new(episodedto.EpisodeRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
		}

		validation := validator.New()
		err := validation.Struct(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		episode := models.Episode{
			Title:         request.Title,
			Thumbnailfilm: request.Thumbnailfilm,
			Linkfilm:      request.Linkfilm,
			FilmID:        request.FilmID,
		}

		episode, err = h.EpisodeRepository.CreateEpisode(episode)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		episode, _ = h.EpisodeRepository.GetEpisode(episode.ID)

		return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Product data created successfully", Data: convertResponseEpisode(episode)})
	} else {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusUnauthorized, Message: "Sorry, you're not Admin"})
	}
}

func (h *handlerEpisode) DeleteEpisode(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userAdmin := userLogin.(jwt.MapClaims)["is_admin"].(bool)
	if userAdmin {
		id, _ := strconv.Atoi(c.Param("id"))

		episode, err := h.EpisodeRepository.GetEpisode(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
		}

		data, err := h.EpisodeRepository.DeleteEpisode(episode)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}
		return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Episode deleted successfully", Data: convertResponseEpisode(data)})
	} else {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusUnauthorized, Message: "Sorry, you're not Admin"})
	}
}

func (h *handlerEpisode) UpdateEpisode(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userAdmin := userLogin.(jwt.MapClaims)["is_admin"].(bool)
	if userAdmin {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		request := new(episodedto.EpisodeRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
		}

		validation := validator.New()
		err = validation.Struct(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		episode, err := h.EpisodeRepository.GetEpisode(int(id))
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
		}

		if request.Title != "" {
			episode.Title = request.Title
		}
		if request.Thumbnailfilm != "" {
			episode.Thumbnailfilm = request.Thumbnailfilm
		}
		if request.Linkfilm != "" {
			episode.Linkfilm = request.Linkfilm
		}

		data, err := h.EpisodeRepository.UpdateEpisode(episode)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Episode updated successfully", Data: convertResponseEpisode(data)})
	} else {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusUnauthorized, Message: "Sorry, you're not Admin"})
	}
}

func convertResponseEpisode(u models.Episode) episodedto.EpisodeResponse {
	return episodedto.EpisodeResponse{
		ID:            u.ID,
		Title:         u.Title,
		Thumbnailfilm: u.Thumbnailfilm,
		Linkfilm:      u.Linkfilm,
		Film:          u.FilmID,
		FilmID:        models.FilmResponse(u.Film),
	}
}
