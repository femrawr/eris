package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	Domain string `json:"domain"`
}

type RedditAbout struct {
	Data struct {
		Name         string  `json:"name"`
		TotalKarma   int     `json:"total_karma"`
		CommentKarma int     `json:"comment_karma"`
		AcceptPMs    bool    `json:"accept_pms"`
		Created      float64 `json:"created"`
		Subreddit    struct {
			Title       string `json:"title"`
			Description string `json:"public_description"`
			Adult       bool   `json:"over_18"`
		} `json:"subreddit"`
	} `json:"data"`
}

type ErisPayload struct {
	Type         string `json:"type"`
	Name         string `json:"user_name"`
	Title        string `json:"title_name"`
	Description  string `json:"description"`
	Adult        bool   `json:"over_18"`
	TotalKarma   int    `json:"total_karma"`
	CommentKarma int    `json:"comment_karma"`
	AcceptDMs    bool   `json:"accept_dms"`
	Created      int    `json:"created"`
}

type ErisResponse struct {
	Ok  bool   `json:"ok"`
	Err string `json:"err"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("[!] you need to provide a user")
		return
	}

	user := args[1]
	fmt.Println("[-] getting for -", user)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://www.reddit.com/user/%s/about.json", user), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[!] failed to fetch -", err)
		return
	}

	defer res.Body.Close()

	var about RedditAbout
	json.NewDecoder(res.Body).Decode(&about)

	data := about.Data

	payload := ErisPayload{
		Type:         "account",
		Name:         data.Name,
		Title:        data.Subreddit.Title,
		Description:  data.Subreddit.Description,
		Adult:        data.Subreddit.Adult,
		TotalKarma:   data.TotalKarma,
		CommentKarma: data.CommentKarma,
		AcceptDMs:    data.AcceptPMs,
		Created:      int(data.Created),
	}

	body, _ := json.Marshal(payload)

	config := config()

	res, err = http.Post(config.Domain+"/api/reddit.php", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[!] failed to send to eris server -", err)
		return
	}

	defer res.Body.Close()

	body, _ = io.ReadAll(res.Body)

	var response ErisResponse
	json.Unmarshal(body, &response)

	if response.Ok {
		fmt.Println("[-] profile saved")
	} else {
		fmt.Println("[!] failed to send to eris server -", response.Err)
	}
}

func config() Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("[!] failed to get config -", err)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	return config
}
