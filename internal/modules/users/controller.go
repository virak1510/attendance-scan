package users

import (
	"attendance/pkg"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	UserService *Service
	validate    *validator.Validate
}

func NewController(u *Service) *Controller {
	return &Controller{
		UserService: u,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (u *Controller) Login(c *gin.Context) {
	loginParam := LoginParams{}

	if err := c.ShouldBindJSON(&loginParam); err != nil {
		response := pkg.BuildResponse(http.StatusBadRequest, err.Error(), pkg.Null())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err := u.validate.Struct(loginParam)
	if err != nil {
		errors := pkg.BuildErrorData(err)
		response := pkg.BuildResponse(http.StatusBadRequest, pkg.ValidateErr, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loginPayload, err := u.UserService.Login(c, loginParam)
	if err != nil {
		fmt.Println(err)
		response := pkg.BuildResponse(http.StatusBadRequest, err.Error(), pkg.Null())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, loginPayload)
	c.SetCookie("token", loginPayload.Token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, response)
}

func (u *Controller) Register(c *gin.Context) {
	user := CreateUserParams{}

	if err := c.ShouldBindJSON(&user); err != nil {
		response := pkg.BuildResponse(http.StatusBadRequest, err.Error(), pkg.Null())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := u.validate.Struct(user)
	if err != nil {
		errors := pkg.BuildErrorData(err)
		response := pkg.BuildResponse(http.StatusBadRequest, pkg.ValidateErr, errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createdUser, err := u.UserService.RegisterUser(c, user)
	if err != nil {
		fmt.Println(err)
		response := pkg.BuildResponse(http.StatusInternalServerError, err.Error(), pkg.Null())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, createdUser)
	c.JSON(http.StatusOK, response)
}

func (u *Controller) GetAllUsers(c *gin.Context) {
	users, err := u.UserService.GetAllUsers(c)
	if err != nil {
		response := pkg.BuildResponse(http.StatusInternalServerError, err.Error(), pkg.Null())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, users)
	c.JSON(http.StatusOK, response)
}
