package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func getIpfsApiUrl() string {
	url := os.Getenv("IPFS_API")
	if url == "" {
		url = "http://127.0.0.1:5001"
	}
	return url
}

func getIpfsGatewayUrl() string {
	url := os.Getenv("IPFS_GATEWAY")
	if url == "" {
		url = "http://127.0.0.1:8081"
	}
	return url
}

// 上传文件，返回CID
func IpfsUpload(reader io.Reader, filename string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", getIpfsApiUrl()+"/api/v0/add", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type ipfsResp struct {
		Hash string `json:"Hash"`
	}
	var result ipfsResp
	body, _ := io.ReadAll(resp.Body)
	lines := bytes.Split(body, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		_ = json.Unmarshal(line, &result)
	}
	if result.Hash == "" {
		return "", fmt.Errorf("IPFS上传失败: %s", string(body))
	}
	return result.Hash, nil
}

// 下载文件，返回 io.ReadCloser
func IpfsDownload(cid string) (io.ReadCloser, error) {
	url := getIpfsApiUrl() + "/api/v0/cat?arg=" + cid
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("IPFS下载请求创建失败: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("IPFS下载失败: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("IPFS下载失败: %s", resp.Status)
	}
	return resp.Body, nil
}
