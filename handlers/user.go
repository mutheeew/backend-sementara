package handlers

import (
	dto "dumbflix/dto/result"
	usersdto "dumbflix/dto/user"
	"dumbflix/models"
	"dumbflix/repositories"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handler struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handler {
	return &handler{UserRepository}
}

func (h *handler) FindUsers(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userAdmin := userLogin.(jwt.MapClaims)["is_admin"].(bool)
	if userAdmin {
		users, err := h.UserRepository.FindUsers()
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
		}

		if len(users) > 0 {
			return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "Data for all users was successfully obtained", Data: users})
		} else {
			return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "No record found"})
		}
	} else {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusUnauthorized, Message: "Sorry, you're not Admin"})
	}
}

func (h *handler) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.UserRepository.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "User data successfully obtained", Data: user})
}

func (h *handler) UpdateUser(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	request := new(usersdto.UpdateUserRequest)
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user, err := h.UserRepository.GetUser(int(userId))

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user.IsAdmin = request.IsAdmin

	if request.Fullname != "" {
		user.Fullname = request.Fullname
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		user.Password = request.Password
	}

	if request.Gender != "" {
		user.Gender = request.Gender
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}

	if request.Address != "" {
		user.Address = request.Address
	}

	data, err := h.UserRepository.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "User data updated successfully", Data: convertResponse(data)})
}

func (h *handler) DeleteUser(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	user, err := h.UserRepository.GetUser(int(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.UserRepository.DeleteUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Message: "User data deleted successfully", Data: convertResponse(data)})
}

func convertResponse(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		ID:       u.ID,
		Fullname: u.Fullname,
		Email:    u.Email,
		Password: u.Password,
		Gender:   u.Gender,
		Phone:    u.Phone,
		Address:  u.Address,
	}
}
