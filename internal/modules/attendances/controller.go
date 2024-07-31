package attendances

import (
	"attendance/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Attendanceservice *Service
}

func NewController(u *Service) *Controller {
	return &Controller{
		Attendanceservice: u,
	}
}

func (a *Controller) GetAttendances(c *gin.Context) {
	att, err := a.Attendanceservice.GetAttendances(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, att)
	c.JSON(http.StatusOK, response)
}

func (a *Controller) CheckIn(c *gin.Context) {

	user, exist := c.Get("user")
	if !exist {
		response := pkg.BuildResponse(http.StatusUnauthorized, pkg.Unauthorized, pkg.Null())
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	authUser := user.(*pkg.AuthUser)

	attendance := AttendanceParams{}

	if err := c.ShouldBindJSON(&attendance); err != nil {
		response := pkg.BuildResponse(http.StatusBadRequest, err.Error(), pkg.Null())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := a.Attendanceservice.CheckIn(c, authUser)
	if err != nil {
		response := pkg.BuildResponse(http.StatusBadRequest, err.Error(), pkg.Null())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, pkg.Null())
	c.JSON(http.StatusOK, response)
}

func (a *Controller) CheckOut(c *gin.Context) {

	user, exist := c.Get("user")
	if !exist {
		response := pkg.BuildResponse(http.StatusUnauthorized, pkg.Unauthorized, pkg.Null())
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	authUser := user.(*pkg.AuthUser)

	err := a.Attendanceservice.CheckOut(c, authUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := pkg.BuildResponse(http.StatusOK, pkg.SUCCESS, pkg.Null())
	c.JSON(http.StatusOK, response)
}
