package tongbu

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	_ "github.com/lib/pq" //
	"log"
	"net/http"
	"time"
)

// URL 响应结果结构体

type ResponseBody struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type ResponseTime struct {
	Url  string  `json:"url"`
	Time float64 `json:"time"` // 响应时间，单位为秒
}

var db *sql.DB

// 连接数据库
func connectDB() *sql.DB {
	dsn := "user=postgres password=982655 dbname=02 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cant connect to the database: %v", err)
	}
	return db
}

func init() {
	db = connectDB()
}

// 从 url_list 表中读取所有 URL

func getURLsFromDB() ([]struct {
	ID  int
	URL string
}, error) {
	query := `SELECT id, url FROM url_list`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []struct {
		ID  int
		URL string
	}

	for rows.Next() {
		var urlData struct {
			ID  int
			URL string
		}
		if err := rows.Scan(&urlData.ID, &urlData.URL); err != nil {
			return nil, err
		}
		urls = append(urls, urlData)
	}

	return urls, nil
}

// 插入 URL 响应内容和时间到 urls_requests 表
func insertResponseData(url string, content string, responseTime float64) error {
	query := `INSERT INTO urls_requests (url, response_body, response_time,status) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, url, content, responseTime)
	if err != nil {
		return err
	}
	return nil
}

// 获取多个 URL 的内容和响应时间（从数据库中读取 URL）
func fetchContentFromURLs() ([]ResponseBody, []ResponseTime, float64) {
	client := resty.New()
	var results []ResponseBody
	var times []ResponseTime

	// 从数据库获取 URL 列表
	urls, err := getURLsFromDB()
	if err != nil {
		log.Printf("failed to get the url from database: %v", err)
		return nil, nil, 0
	}

	startTime := time.Now()
	// 顺序请求每个 URL
	for _, urlData := range urls {
		url := urlData.URL

		startTime := time.Now()
		log.Printf("request URL: %s", url)

		resp, err := client.R().Get(url)
		elapsedTime := time.Since(startTime).Seconds()

		if err != nil {
			log.Printf("request %s failed: %v", url, err)
			// 将失败的请求信息写入数据库
			insertResponseData(url, "cant get the context", elapsedTime)
			results = append(results, ResponseBody{Url: url, Content: "cant get the context"})
			times = append(times, ResponseTime{Url: url, Time: elapsedTime})
			continue
		}

		content := string(resp.Body())
		if len(content) > 1000 {
			content = content[:1000] // 只返回前 1000 个字符
		}

		// 将成功的响应数据存储到数据库
		err = insertResponseData(url, content, elapsedTime)
		if err != nil {
			log.Printf("failed to insert the response message: %v", err)
		}

		results = append(results, ResponseBody{Url: url, Content: content})
		times = append(times, ResponseTime{Url: url, Time: elapsedTime})
	}

	totalDuration := time.Since(startTime).Seconds()
	log.Printf("totalduration: %.2f s", totalDuration)
	return results, times, totalDuration
}

// 处理 /fetch-websites 请求

func FetchFromWebsites(c *gin.Context) {
	results, times, totalDuration := fetchContentFromURLs()
	c.JSON(http.StatusOK, gin.H{
		"results":        results,
		"response_times": times,
		"total_duration": totalDuration,
	})
}
