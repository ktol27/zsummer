package webpage

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"strings"
)

// GetVisibleContent 获取指定 URL 的可见内容
func GetVisibleContent(url string) (string, error) {
	// 参数校验
	if url == "" {
		return "", errors.New("url cannot be empty")
	}

	// 创建 resty 客户端
	client := resty.New()

	// 发送 GET 请求
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}

	// 检查响应状态
	if resp.StatusCode() != 200 {
		return "", errors.New("request failed with status code: " + resp.Status())
	}

	// 获取响应内容
	content := string(resp.Body())

	// 清理 HTML 标签和多余空白
	content = cleanHTML(content)
	content = strings.TrimSpace(content)

	return content, nil
}

// cleanHTML 清理 HTML 标签
func cleanHTML(html string) string {
	// 移除 script 标签及内容
	html = removeTag(html, "script")

	// 移除 style 标签及内容
	html = removeTag(html, "style")

	// 移除其他 HTML 标签
	html = removeAllTags(html)

	return html
}

// removeTag 移除指定标签及其内容
func removeTag(html string, tag string) string {
	startTag := "<" + tag
	endTag := "</" + tag + ">"

	for {
		start := strings.Index(html, startTag)
		if start == -1 {
			break
		}

		end := strings.Index(html[start:], endTag)
		if end == -1 {
			break
		}

		html = html[:start] + html[start+end+len(endTag):]
	}

	return html
}

// removeAllTags 移除所有 HTML 标签
func removeAllTags(html string) string {
	for {
		start := strings.Index(html, "<")
		if start == -1 {
			break
		}

		end := strings.Index(html[start:], ">")
		if end == -1 {
			break
		}

		html = html[:start] + " " + html[start+end+1:]
	}

	return html
}
