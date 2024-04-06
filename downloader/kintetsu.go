package downloader

import (
	"io"
	"net/http"
)

func GetData() ([]byte, error) {
	resp, err := http.Get("https://tid.kintetsu.co.jp/LocationHtml/trainlocationinfo01.html?innerLink=true")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
