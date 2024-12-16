package restypage

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"time"
)

type RequestBody struct {
	Urls []string `json:"urls" binding:"required,dive,required,url"`
}

type ResponseBody struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type TimeResponse struct {
	Times []float64 `json:"times"`
}

// create a request to fetch content from multiple urls
func fetchContentFromURLs(urls []string) ([]ResponseBody, TimeResponse) {
	client := resty.New()
	var results []ResponseBody
	var times []float64
	resultCh := make(chan ResponseBody, len(urls))
	timesCh := make(chan float64, len(urls))

	starTime := time.Now()
	// send requests based on concurrency level /use GORoutine
	for _, url := range urls {
		go fetchSingleURL(client, url, resultCh, timesCh)
	}

	// 收集结果和响应时间
	for i := 0; i < len(urls); i++ {
		results = append(results, <-resultCh)
		times = append(times, <-timesCh)
	}

	totalDuration := time.Since(starTime).Seconds()

	log.Printf("total duration is %.2f s", totalDuration)
	return results, TimeResponse{Times: times}
}

// fetch content from multiple urls

func FetchFromWebsites(c *gin.Context) {
	var requestBody RequestBody

	//Bind the JSON data in the request body to the structure
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//"error":   "请求体无效，请提供 {\"urls\": [\"url1\", \"url2\"]} 格式的 JSON 数据",
			//"details": err.Error(),
		})
		return
	}

	results, times := fetchContentFromURLs(requestBody.Urls)
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"times":   times.Times, // 返回响应时间数组
	})
}

func fetchSingleURL(
	client *resty.Client,
	url string,
	resultCh chan<- ResponseBody,
	timesCh chan<- float64) {

	startTime := time.Now() // 记录开始时间
	log.Printf("开始请求 URL: %s", url)

	resp, err := client.R().Get(url)
	//
	elapsedTime := time.Since(startTime).Seconds() // 计算响应时间（秒）

	if err != nil {
		log.Printf("请求 %s 失败: %v", url, err)
		resultCh <- ResponseBody{Url: url, Content: "无法获取内容"}
		timesCh <- elapsedTime
		return
	}

	content := string(resp.Body())
	if len(content) > 1000 {
		content = content[:1000] // 只返回前1000个字符
	}

	resultCh <- ResponseBody{
		Url:     url,
		Content: content,
	}

	timesCh <- elapsedTime // 将响应时间发送到 channel
}
