package routes

import (
	"dumbflix/handlers"
	"dumbflix/pkg/middleware"
	"dumbflix/pkg/mysql"
	"dumbflix/repositories"

	"github.com/labstack/echo/v4"
)

func FilmRoutes(e *echo.Group) {
	filmRepository := repositories.RepositoryFilm(mysql.DB)
	h := handlers.HandlerFilm(filmRepository)

	e.GET("/film/:id", middleware.Auth(h.GetFilm))
	e.GET("/films", h.FindFilms)
	e.POST("/film", middleware.Auth(middleware.UploadFile(h.CreateFilm)))
	// e.PATCH("/film", middleware.Auth(middleware.UploadFile(h.UpdateFilm)))
	e.DELETE("/film", middleware.Auth(h.DeleteFilm))
}
