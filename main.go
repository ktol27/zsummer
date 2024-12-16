package main

import (
	"angular/Db2"
	"angular/LOG"
	"angular/api"
	"angular/restypage"
	_ "angular/restypage"
	"angular/tongbu"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	LOG.InitLogger()
	Db2.InitDatabase()

	r := gin.Default()
	r.Use(cors.Default())
	r.Static("/uploads", "./uploads")

	r.POST("/users/:id/avatar", api.UploadAvatar)
	r.GET("/users", api.GetUsers)
	r.POST("/users", api.AddUser)
	r.DELETE("/users/:username", api.DeleteUser)
	//r.GET("/getUrl", api.GetUrls)
	r.POST("/urls/fetch", restypage.FetchFromWebsites) // Concurrent fetch
	r.POST("/urls/fetching", tongbu.FetchFromWebsites) // Serial fetch

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
