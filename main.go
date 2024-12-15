package main

import (
	"angular/Db2"
	"angular/LOG"
	"angular/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//type RequestBody struct {
//	Urls []string `json:"urls"`
//}
//
//type ResponseBody struct {
//	Url     string `json:"url"`
//	Content string `json:"content"`
//}
//
//func fetchContentFromURL(url string, client *resty.Client) (string, error) {
//	// 使用 Resty 发起 GET 请求
//	resp, err := client.R().
//		SetHeader("User-Agent", "Go-Resty-Client").
//		Get(url)
//	if err != nil {
//		return "", fmt.Errorf("请求 %s 失败: %v", url, err)
//	}
//
//	// 只返回部分内容，防止响应过大
//	body := resp.String()
//	if len(body) > 1000 {
//		body = body[:1000]
//	}
//
//	return body, nil
//}

//func fetchFromWebsites(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
//		return
//	}
//
//	var requestBody RequestBody
//	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//		http.Error(w, "无效的请求数据", http.StatusBadRequest)
//		return
//	}
//
//	if len(requestBody.Urls) == 0 {
//		http.Error(w, "请提供至少一个 URL", http.StatusBadRequest)
//		return
//	}
//
//	// 初始化 Resty 客户端
//	client := resty.New().
//		SetTimeout(10 * time.Second) // 设置请求超时时间
//
//	var results []ResponseBody
//	for _, url := range requestBody.Urls {
//		content, err := fetchContentFromURL(url, client)
//		if err != nil {
//			log.Printf("无法获取 %s: %v", url, err)
//			results = append(results, ResponseBody{Url: url, Content: "无法获取内容"})
//			continue
//		}
//		results = append(results, ResponseBody{Url: url, Content: content})
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	if err := json.NewEncoder(w).Encode(results); err != nil {
//		http.Error(w, "编码响应失败", http.StatusInternalServerError)
//		return
//	}
//}

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
	r.GET("/getUrl", api.GetUrls)
	//http.HandleFunc("/fetch-websites", fetchFromWebsites)

	if err := r.Run(":8080"); err != nil {
		LOG.ErrorLogger.Error("Server error:", err)
	}
}
