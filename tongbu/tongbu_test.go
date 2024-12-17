package tongbu

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-resty/resty/v2"
	"github.com/h2non/gock"
	_ "log"
	"testing"
)

// 模拟的 URLData 结构体
type URLData struct {
	ID  int
	URL string
}

// database failed connection
func TestDBConnectionFailure(t *testing.T) {
	// create database
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("error creating mock database connection: %v", err)
	}
	defer db.Close()

	//
	mock.ExpectPing().WillReturnError(errors.New("database connection failed"))

	//
	err = db.Ping()
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if err.Error() != "database connection failed" {
		t.Fatalf("expected error 'database connection failed', but got %v", err)
	}
}

// ** 模拟从数据库中获取 URL 列表的函数 **
func GetURLsFromDB(db *sql.DB) ([]URLData, error) {
	rows, err := db.Query("SELECT id, url FROM url_list")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URLData
	for rows.Next() {
		var urlData URLData
		if err := rows.Scan(&urlData.ID, &urlData.URL); err != nil {
			return nil, err
		}
		urls = append(urls, urlData)
	}
	return urls, nil
}

// ** 模拟调用多个 URL 请求的函数 **
func testFetchContentFromURLs(client *resty.Client, urls []URLData) ([]ResponseBody, error) {
	var results []ResponseBody
	for _, urlData := range urls {
		resp, err := client.R().Get(urlData.URL)
		if err != nil {
			return nil, err
		}
		content := string(resp.Body())
		results = append(results, ResponseBody{Url: urlData.URL, Content: content})
	}
	return results, nil
}

// ** 使用 sqlmock 和 gock 模拟外部请求的单元测试 **
func TestFetchContentFromURLsWithMockDB(t *testing.T) {
	// 1️⃣ **创建 Mock 数据库连接**
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database connection: %v", err)
	}
	defer db.Close()

	// 2️⃣ **模拟 SQL 查询的返回数据**
	rows := sqlmock.NewRows([]string{"id", "url"}).
		AddRow(1, "https://www.google.com").
		AddRow(2, "https://jsonplaceholder.typicode.com/posts/1").
		AddRow(3, "https://jsonplaceholder.typicode.com/todos/1")

	mock.ExpectQuery("SELECT id, url FROM url_list").WillReturnRows(rows)

	// 3️⃣ **从 mock 数据库中获取 URL 列表**
	urls, err := GetURLsFromDB(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) != 3 {
		t.Fatalf("expected 3 urls, but got %d", len(urls))
	}

	// 4️⃣ **模拟外部请求 (gock)**
	defer gock.Off() // 清除 gock 的所有拦截器

	gock.New("https://www.goo.com").
		Reply(200).
		BodyString("This is a mock Goo page")

	gock.New("https://jsonplaceholder.typicode.com").
		Path("/posts/1").
		Reply(200).
		BodyString(`{"title": "mock title", "body": "mock body"}`)

	gock.New("https://jsonplaceholder.typicode.com").
		Path("/todos/1").
		Reply(200).
		BodyString(`{"task": "mock task", "status": "completed"}`)

	// 5️⃣ **调用 testFetchContentFromURLs 函数，模拟 URL 内容的返回**
	client := resty.New()
	content, err := testFetchContentFromURLs(client, urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(content) != 3 {
		t.Fatalf("expected 3 urls to be fetched, but got %d", len(content))
	}

}
