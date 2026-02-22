package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Domain string `json:"domain"`
}

type RedditAbout struct {
	Data struct {
		Name         string  `json:"name"`
		ID           string  `json:"id"`
		TotalKarma   int     `json:"total_karma"`
		CommentKarma int     `json:"comment_karma"`
		AcceptPMs    bool    `json:"accept_pms"`
		Created      float64 `json:"created_utc"`
		Subreddit    struct {
			Title       string `json:"title"`
			Description string `json:"public_description"`
			Adult       bool   `json:"over_18"`
		} `json:"subreddit"`
	} `json:"data"`
}

type RedditPosts struct {
	Data struct {
		Children []struct {
			Post struct {
				UserID      string  `json:"author_fullname"`
				PostID      string  `json:"id"`
				SubReddit   string  `json:"subreddit"`
				Body        string  `json:"selftext"`
				Title       string  `json:"title"`
				Upvotes     int     `json:"ups"`
				UpvoteRatio float64 `json:"upvote_ratio"`
				Adult       bool    `json:"over_18"`
				Created     float64 `json:"created_utc"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditComments struct {
	Data struct {
		Children []struct {
			Comment struct {
				UserID    string  `json:"author_fullname"`
				PostID    string  `json:"id"`
				SubReddit string  `json:"subreddit"`
				ParentID  string  `json:"parent_id"`
				Body      string  `json:"body"`
				Upvotes   int     `json:"ups"`
				Created   float64 `json:"created_utc"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type ErisProfilePayload struct {
	Type         string `json:"type"`
	Name         string `json:"user_name"`
	ID           string `json:"user_id"`
	Title        string `json:"title_name"`
	Description  string `json:"description"`
	Adult        bool   `json:"over_18"`
	TotalKarma   int    `json:"total_karma"`
	CommentKarma int    `json:"comment_karma"`
	AcceptDMs    bool   `json:"accept_dms"`
	Created      int    `json:"created"`
}

type ErisPostPayload struct {
	Type      string `json:"type"`
	Name      string `json:"user_name"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	SubReddit string `json:"sub_reddit"`
	Body      string `json:"body"`
	Title     string `json:"title"`
	Upvotes   int    `json:"upvotes"`
	Ratio     int    `json:"vote_ratio"`
	Adult     bool   `json:"over_18"`
	Created   int    `json:"created"`
}

type ErisCommentPayload struct {
	Type      string `json:"type"`
	Name      string `json:"user_name"`
	UserID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	SubReddit string `json:"sub_reddit"`
	ParentID  string `json:"parent_id"`
	Body      string `json:"body"`
	Upvotes   int    `json:"upvotes"`
	Created   int    `json:"created"`
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

	config := getConfig()

	client := &http.Client{}

	body, err := scrapeReddit(client, user, "about")
	if err != nil {
		fmt.Println("[!] failed to scrape profile -", err)
		return
	}

	var about RedditAbout
	json.Unmarshal(body, &about)
	data := about.Data

	res, err := sendServer(client, config.Domain, ErisProfilePayload{
		Type:         "profile",
		Name:         data.Name,
		ID:           data.ID,
		Title:        data.Subreddit.Title,
		Description:  data.Subreddit.Description,
		Adult:        data.Subreddit.Adult,
		TotalKarma:   data.TotalKarma,
		CommentKarma: data.CommentKarma,
		AcceptDMs:    data.AcceptPMs,
		Created:      int(data.Created),
	})

	if err != nil || !res.Ok {
		fmt.Println("[!] failed to save account -", res.Err)
		return
	}

	body, err = scrapeReddit(client, user, "submitted")
	if err != nil {
		fmt.Println("[!] failed to scrape posts -", err)
		return
	}

	var posts RedditPosts
	json.Unmarshal(body, &posts)

	for _, child := range posts.Data.Children {
		post := child.Post

		res, err = sendServer(client, config.Domain, ErisPostPayload{
			Type:      "posts",
			Name:      data.Name,
			UserID:    data.ID,
			PostID:    post.PostID,
			SubReddit: post.SubReddit,
			Body:      post.Body,
			Title:     post.Title,
			Upvotes:   post.Upvotes,
			Ratio:     int(post.UpvoteRatio * 100),
			Adult:     post.Adult,
			Created:   int(post.Created),
		})

		if err != nil || !res.Ok {
			fmt.Println("[!] failed to save post -", res.Err)
			continue
		}
	}

	body, err = scrapeReddit(client, user, "comments")
	if err != nil {
		fmt.Println("[!] failed to scrape comments -", err)
		return
	}

	var comments RedditComments
	json.Unmarshal(body, &comments)

	for _, child := range comments.Data.Children {
		comment := child.Comment

		res, err = sendServer(client, config.Domain, ErisCommentPayload{
			Type:      "comments",
			Name:      data.Name,
			UserID:    data.ID,
			PostID:    comment.PostID,
			SubReddit: comment.SubReddit,
			ParentID:  comment.ParentID,
			Body:      comment.Body,
			Upvotes:   comment.Upvotes,
			Created:   int(comment.Created),
		})

		if err != nil || !res.Ok {
			fmt.Println("[!] failed to save comment -", res.Err)
			continue
		}
	}

	fmt.Println("[-] finished with -", user)
}

func scrapeReddit(client *http.Client, user string, what string) ([]byte, error) {
	if what == "about" {
		url := fmt.Sprintf("https://www.reddit.com/user/%s/about.json", user)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		return io.ReadAll(res.Body)
	}

	type Listing struct {
		Data struct {
			After    any   `json:"after"`
			Children []any `json:"children"`
		} `json:"data"`
	}

	var allChildren []any
	after := ""

	for {
		url := fmt.Sprintf("https://www.reddit.com/user/%s/%s.json?limit=100", user, what)
		if after != "" {
			url += "&after=" + after
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		var listing Listing
		json.Unmarshal(body, &listing)

		allChildren = append(allChildren, listing.Data.Children...)

		if listing.Data.After == nil {
			break
		}

		after = listing.Data.After.(string)
		time.Sleep(1 * time.Second)
	}

	result := map[string]any{
		"data": map[string]any{
			"children": allChildren,
			"after":    nil,
		},
	}

	return json.Marshal(result)
}

func sendServer(client *http.Client, domain string, payload any) (*ErisResponse, error) {
	body, _ := json.Marshal(payload)

	res, err := client.Post(domain+"/api/reddit.php", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, _ = io.ReadAll(res.Body)

	var response ErisResponse
	json.Unmarshal(body, &response)

	return &response, nil
}

func getConfig() Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("[!] failed to get config -", err)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	return config
}
