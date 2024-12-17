

// ** ResponseBody 结构体定义 **
type ResponseBody struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type URLData struct {
	ID  int
	URL string
}

// ** testFetchContentFromURLs - 使用 gock 模拟请求内容 **
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

// ** 使用 gock 模拟 URL 请求的单元测试 **
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

	// 3️⃣ **调用 getURLsFromDB 函数从 mock 数据库中查询 URL**
	urls, err := getURLsFromDB(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) != 3 {
		t.Fatalf("expected 3 urls, but got %d", len(urls))
	}

	// 4️⃣ **使用 gock 模拟外部请求**
	defer gock.Off() // 清除 gock 的所有拦截器

	gock.New("https://www.google.com").
		Reply(200).
		BodyString("This is a mock Google page")

	gock.New("https://jsonplaceholder.typicode.com").
		MatchPath("/posts/1").
		Reply(200).
		BodyString(`{"title": "mock title", "body": "mock body"}`)

	gock.New("https://jsonplaceholder.typicode.com").
		MatchPath("/todos/1").
		Reply(200).
		BodyString(`{"task": "mock task", "status": "completed"}`)

	// 5️⃣ **调用 testFetchContentFromURLs 函数，并模拟返回数据**
	client := resty.New()
	content, err := testFetchContentFromURLs(client, urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(content) != 3 {
		t.Fatalf("expected 3 urls to be fetched, but got %d", len(content))
	}

	// 6️⃣ **验证返回的内容**
	if content[0].Content != "This is a mock Google page" {
		t.Fatalf("expected 'This is a mock Google page', but got '%s'", content[0].Content)
	}

	if !gock.IsDone() {
		t.Fatalf("not all gock requests were made")
	}
}
