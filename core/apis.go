package core

import (
	"net/http"
	"time"

	"github.com/claptrap666/go-code-counter/middleware"
	"github.com/claptrap666/go-code-counter/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var APIServer = gin.Default()

var SESSION_SECRET string = "changeme"

func InitRouter() {
	store := cookie.NewStore([]byte(SESSION_SECRET))
	APIServer.Use(sessions.Sessions("mysession", store))
	APIServer.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	APIServer.Static("/ui", "./static")
	apis := APIServer.Group("/api/v1")
	// apis.Use(middleware.OidcMiddleware(
	// 	CurrentConfig.Oauth.Provider,
	// 	CurrentConfig.Oauth.ClientID,
	// 	CurrentConfig.Oauth.ClientSecret,
	// 	CurrentConfig.Oauth.RedirectURL))
	apis.Use(gin.BasicAuth(gin.Accounts{
		"admin": CurrentConfig.Server.Secret,
	}))
	apis.POST("commits", ListCommits)
	apis.POST("count", CountByTime)
	APIServer.GET("/oauth_callback", middleware.HandleOAuth2Callback)
	apis.Use(middleware.CasbinMiddleware())

}

type query struct {
	PageIndex int `json:"pageindex"`
	PageSize  int `json:"pagesize"`
}

func ListCommits(c *gin.Context) {
	commitsObjs := []models.T_commits{}
	var q = query{}
	err := c.ShouldBindJSON(&q)
	if err != nil {
		c.String(http.StatusOK, `请求格式有错`)
		c.JSON(http.StatusInternalServerError, NewMessage(MESSAGE_ERROR, err))
		return
	}
	logrus.Info("query: %v", q)
	tx := CurrentDB.Limit(q.PageSize).Offset(q.PageSize * (q.PageIndex - 1)).Where("addlines < 10000").Find(&commitsObjs)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, NewMessage(MESSAGE_ERROR, tx.Error))
		return
	}
	c.JSON(http.StatusOK, NewMessage(MESSAGE_OK, commitsObjs))
}

type query2 struct {
	Commiters []string  `json:"commiters"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}

type result struct {
	Commiter string
	Total    int
}

func CountByTime(c *gin.Context) {
	var q = query2{}
	err := c.ShouldBindJSON(&q)
	if err != nil {
		c.String(http.StatusOK, `请求格式有错`)
		c.JSON(http.StatusInternalServerError, NewMessage(MESSAGE_ERROR, err))
		return
	}
	logrus.Info("query: %v", q)
	r := []result{}
	tx := CurrentDB.
		Model(&models.T_commits{}).
		Select("commiter, sum(addlines) as total").
		Where("commiter in (?)", q.Commiters).
		Where("cdate > ?", q.StartDate).
		Where("cdate < ?", q.EndDate).
		Where("addlines < 10000").
		Group("commiter").
		Find(&r)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, NewMessage(MESSAGE_ERROR, tx.Error))
		return
	}
	c.JSON(http.StatusOK, NewMessage(MESSAGE_OK, r))
}
