package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v45/github"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
	PubDate     string `xml:"pubDate"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func getGitHubClient() (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}
	ts := github.BasicAuthTransport{
		Username: token,
	}
	return github.NewClient(ts.Client()), nil
}

func GenerateRSSFeed(username string, cache interface {
	Get(string) (*RSS, bool)
	Set(string, *RSS)
}) (*RSS, error) {
	if cachedFeed, found := cache.Get(username); found {
		return cachedFeed, nil
	}

	client, err := getGitHubClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	stars, _, err := client.Activity.ListStarred(ctx, username, opt)
	if err != nil {
		return nil, fmt.Errorf("error fetching starred repositories: %v", err)
	}

	feed := &RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       fmt.Sprintf("GitHub Stars - %s", username),
			Link:        fmt.Sprintf("https://github.com/%s?tab=stars", username),
			Description: fmt.Sprintf("RSS feed of repositories starred by %s", username),
			PubDate:     time.Now().Format(time.RFC1123Z),
		},
	}

	for _, star := range stars {
		repo := star.GetRepository()
		item := Item{
			Title:       repo.GetFullName(),
			Link:        repo.GetHTMLURL(),
			Description: repo.GetDescription(),
			PubDate:     star.GetStarredAt().Format(time.RFC1123Z),
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	cache.Set(username, feed)

	return feed, nil
}

func GenerateMultiUserRSSFeed(usernames []string, cache interface {
	Get(string) (*RSS, bool)
	Set(string, *RSS)
}) (*RSS, error) {
	var allItems []Item
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(usernames))

	for _, username := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()

			feed, err := GenerateRSSFeed(username, cache)
			if err != nil {
				errChan <- fmt.Errorf("error fetching feed for %s: %v", username, err)
				return
			}

			mu.Lock()

			for _, item := range feed.Channel.Items {
				item.Description = fmt.Sprintf("%s (starred by %s)", item.Description, username)
				allItems = append(allItems, item)
			}
			mu.Unlock()
		}(username)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(allItems, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC1123Z, allItems[i].PubDate)
		timeJ, _ := time.Parse(time.RFC1123Z, allItems[j].PubDate)
		return timeI.After(timeJ)
	})

	feed := &RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "GitHub Stars - Combined Feed",
			Link:        "https://github.com",
			Description: fmt.Sprintf("Combined RSS feed of repositories starred by %s", strings.Join(usernames, ", ")),
			PubDate:     time.Now().Format(time.RFC1123Z),
			Items:       allItems,
		},
	}

	return feed, nil
}
