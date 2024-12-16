package api

import (
	"angular/webpage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUrls(c *gin.Context) {
	url := "http://google.com"
	//url := c.Query("url")
	//url = "http://www." + url + ".com"
	content, _ := webpage.GetVisibleContent(url)
	c.JSON(http.StatusOK, gin.H{"content": content})
}
