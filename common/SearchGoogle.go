package common

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var SerpAPIKey string = ""

func SearchGoogle(query string) (string, error) {
	SerpAPIKey := os.Getenv("SERP_API_KEY")
	searchURL := "https://serpapi.com/search.json?q=" + url.QueryEscape(query) + "&api_key=" + SerpAPIKey

	// Gửi yêu cầu HTTP GET
	resp, err := http.Get(searchURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Đọc phản hồi
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Phân tích kết quả JSON
	results := gjson.GetBytes(body, "organic_results")

	if !results.Exists() {
		return "Không tìm thấy kết quả nào cho: " + query, nil
	}

	// Tạo chuỗi kết quả
	var resultString strings.Builder
	resultString.WriteString("Kết quả tìm kiếm cho: " + query + "\n\n")

	// Giới hạn chỉ lấy 3 kết quả đầu tiên để tránh tin nhắn quá dài
	count := 0
	results.ForEach(func(_, value gjson.Result) bool {
		// dừng sau 3 kết quả có thể tùy chỉnh để lấy nhiều kết quả hơn
		if count >= 3 {
			return false
		}

		title := value.Get("title").String()
		link := value.Get("link").String()
		snippet := value.Get("snippet").String()

		resultString.WriteString(fmt.Sprintf("%d. **%s**\n%s\n%s\n\n", count+1, title, snippet, link))
		count++
		return true
	})

	return resultString.String(), nil
}
