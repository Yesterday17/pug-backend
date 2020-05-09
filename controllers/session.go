package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"

	"github.com/Yesterday17/pug-backend/auth"
	"github.com/Yesterday17/pug-backend/config"
	e "github.com/Yesterday17/pug-backend/error"
	"github.com/Yesterday17/pug-backend/models"
	"github.com/Yesterday17/pug-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func userPass(c *gin.Context) (string, string, error) {
	cfg := c.MustGet("config").(*config.Config)

	name, ok := c.GetPostForm("username")
	if !ok || name == "" {
		return "", "", errors.New("username or password is empty")
	}

	pass, ok := c.GetPostForm("password")
	if !ok || pass == "" {
		return "", "", errors.New("username or password is empty")
	}

	password, err := base64.StdEncoding.DecodeString(pass)
	if err != nil {
		return "", "", errors.New("invalid password")
	}

	realPassword, err := rsa.DecryptPKCS1v15(rand.Reader, cfg.KeyPrivate, []byte(password))
	if err != nil {
		return "", "", err
	}

	return name, string(realPassword), nil
}

func SessionCreate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	name, pass, err := userPass(c)
	if err != nil {
		c.AbortWithStatusJSON(400, e.ErrInputInvalid)
		return
	}

	var user models.User
	db.First(&user, "username = ?", name)
	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBRead)
		return
	}

	if user.UUID == "" || user.Password != pass {
		c.AbortWithStatusJSON(400, e.ErrNoUser)
		return
	}

	c.Set("uuid", user.UUID)
	c.Set("level", user.UserLevel)
}

func SessionUpdate(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)
	id := c.MustGet("uuid").(string)
	level := c.MustGet("level").(int)

	token, err := auth.GenerateToken(cfg.KeyPrivate, id, level, time.Hour*12)
	if err != nil {
		c.AbortWithStatusJSON(500, e.ErrFailTokenGen)
		return
	}

	c.SetCookie("pug_session", token, int(time.Hour*12), "", "", false, true)
}

func SessionRevoke(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	token, err := auth.CancelToken(cfg.KeyPrivate)
	if err != nil {
		c.JSON(500, e.ErrFailTokenGen)
		return
	}

	c.SetCookie("pug_session", token, int(time.Minute*1), "", "", false, true)
	c.JSON(200, e.NoError)
}

func UserRegister(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	name, pass, err := userPass(c)
	if err != nil {
		c.JSON(400, e.ErrInputInvalid)
		return
	}

	var user models.User
	db.First(&user, "username = ?", name)
	if user.UUID != "" {
		c.AbortWithStatusJSON(400, e.ErrUserExist)
		return
	}

	id := uuid.NewV4().String()
	c.Set("uuid", id)
	c.Set("level", 0)

	user = models.User{
		UUID:      id,
		Username:  name,
		Password:  pass,
		UserLevel: 0,
		Setting: models.UserSettings{
			Account: models.UserAccountSettings{
				Name:  utils.GetPostForm(c, "name", ""),
				Email: utils.GetPostForm(c, "email", ""),
				Icon:  utils.GetPostForm(c, "icon", ""),
			},
		},
	}
	db.Create(&user)

	if db.Error != nil {
		c.AbortWithStatusJSON(500, e.ErrDBWrite)
		return
	}
	c.JSON(201, e.NoError)
}

func SessionGetKey(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	c.Data(200, "text/plain", []byte(cfg.PublicKeyString))
}
