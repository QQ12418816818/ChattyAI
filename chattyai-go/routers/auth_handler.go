package routers

import (
	"chattyai-go/models"
	"chattyai-go/utils"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type LoginParam struct {
	Username   string `json:"username"`
	VerifyCode string `json:"verifyCode"`
	Password   string `json:"password"`
}

func LoginByPassword(c *gin.Context) {
	loginParam := LoginParam{}
	c.BindJSON(&loginParam)

	user, err := models.GetUserByUsername(loginParam.Username)
	if err != nil {
		c.JSON(200, NewFailCommonResponse(err.Error()))
	} else if user == nil {
		c.JSON(200, NewFailCommonResponse("Failed to find user"))
	} else if user.PasswordHash != loginParam.Password {
		c.JSON(200, NewFailCommonResponse("Password mismatch"))
	} else {
		tokenId, _ := utils.GenerateNanoId()
		models.CreateAuthToken(tokenId, user.UID)
		c.JSON(200, NewSuccessCommonResponseWithData(tokenId))
	}
}

func AuthCheckToken(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")
	if len(authorizationHeader) == 0 {
		c.JSON(200, NewFailCommonResponse("No Authorization"))
	}

	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		token, err := models.GetAuthTokenByToken(authorizationHeader[7:])
		if err != nil {
			c.JSON(200, NewFailCommonResponse(err.Error()))
		} else if token == nil || token.ExpiredTime.Before(time.Now()) {
			c.JSON(200, NewFailCommonResponse("No Authorization"))
		} else {
			u, _ := models.GetUserByUid(token.UserID)

			c.JSON(200, NewSuccessCommonResponseWithData(&UserDTO{
				Nickname:    u.Nickname,
				Avatar:      u.Avatar,
				Description: u.Description,
			}))
		}
	} else {
		c.JSON(200, NewFailCommonResponse("Mistake Authorization"))
	}
}

type UserDTO struct {
	Nickname    string
	Avatar      string
	Description string
}
