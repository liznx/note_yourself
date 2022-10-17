package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"noteyouself/api"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"

	"noteyouself/db"
)

func TokenAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := c.GetRawData()
		if err != nil {
			fmt.Println(err.Error())
		}
		m := make(map[string]interface{}) //注意该结构接受的内容
		json.Unmarshal(data, &m)
		log.Printf("%v", &m)
		token := m["token"].(string)
		result, _ := db.RedisClient.Get(token).Result()
		if result == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "访问未授权",
			})
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		c.Next()
	}
}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	user := router.Group("/api/v1/wxminiapp/")
	{
		user.POST("/getWxOpenidByCode", api.GetWxOpenidByCode)
	}
	//router.Use(TokenAuthorize())
	//router.POST("/api/v1/wxminiapp/decryptWxData", api.DecryptWxData)
	user.Use(TokenAuthorize())
	user.POST("/decryptWxData", api.DecryptWxData)
	return router
}
