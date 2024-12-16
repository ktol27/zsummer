package tongbu

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type RequestBody struct {
	Urls []string `json:"urls" binding:"required,dive,required,url"`
}

type ResponseBody struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type ResponseTime struct {
	Url  string  `json:"url"`
	Time float64 `json:"time"` // save response time in seconds
}

func FetchFromWebsites(c *gin.Context) {
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//"error":   "请求体无效，请提供 {\"urls\": [\"url1\", \"url2\"]} 格式的 JSON 数据",
			//"details": err.Error(),
		})
		return
	}

	results, responseTimes, totalDuration := FetchContentFromURLs(requestBody.Urls)
	c.JSON(http.StatusOK, gin.H{
		"results":        results,
		"response_times": responseTimes,
		"total_duration": totalDuration,
	})
}

// FetchContentFromURLs fetches content from given URLs and returns the results along with response times
func FetchContentFromURLs(urls []string) ([]ResponseBody, []ResponseTime, float64) {
	client := resty.New().SetTimeout(10 * time.Second) //set timeout to 10 seconds
	var results []ResponseBody
	var responseTimes []ResponseTime
	totalDuration := 0.0 //record total duration

	for _, url := range urls {
		startTime := time.Now()

		resp, err := client.R().Get(url)

		duration := time.Since(startTime).Seconds()
		totalDuration += duration

		if err != nil {
			log.Printf("请求 %s 失败: %v", url, err)
			results = append(results, ResponseBody{Url: url, Content: "无法获取内容"})
			responseTimes = append(responseTimes, ResponseTime{Url: url, Time: duration})
			continue
		}

		content := string(resp.Body())
		if len(content) > 1000 {
			content = content[:1000] // 只返回前1000个字符
		}
		results = append(results, ResponseBody{Url: url, Content: content})
		responseTimes = append(responseTimes, ResponseTime{Url: url, Time: duration})
	}

	return results, responseTimes, totalDuration
}