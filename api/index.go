package api

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"

	"stargazer/pkg/cache"
	"stargazer/pkg/feed"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	router       *mux.Router
	FeedCache    = cache.NewFeedCache()
	MaxUsernames = 100
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	router = mux.NewRouter()
	router.HandleFunc("/feed/{username}", HandleRSSFeed).Methods("GET")
	router.HandleFunc("/feeds/{usernames}", HandleMultiUserRSSFeed).Methods("GET")
	router.HandleFunc("/feeds/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No usernames provided", http.StatusBadRequest)
	}).Methods("GET")
	router.HandleFunc("/", HandleRoot).Methods("GET")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func sendFeedResponse(w http.ResponseWriter, feedData interface{}, err error) {
	if err != nil {
		if err.Error() == "GITHUB_TOKEN environment variable is not set" {
			http.Error(w, "Server configuration error: GitHub token not set", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(feedData); err != nil {
		http.Error(w, "Failed to encode feed", http.StatusInternalServerError)
		return
	}
}

func HandleRSSFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	feedData, err := feed.GenerateRSSFeed(username, FeedCache)
	sendFeedResponse(w, feedData, err)
}

func HandleMultiUserRSSFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	usernames := vars["usernames"]
	if usernames == "" {
		http.Error(w, "No usernames provided", http.StatusBadRequest)
		return
	}

	usernameList := strings.Split(usernames, "+")
	if len(usernameList) > MaxUsernames {
		http.Error(w, fmt.Sprintf("Too many usernames. Maximum allowed is %d", MaxUsernames), http.StatusBadRequest)
		return
	}

	feedData, err := feed.GenerateMultiUserRSSFeed(usernameList, FeedCache)
	sendFeedResponse(w, feedData, err)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`
	{
	  "description": "Stargazer API: Follow what other people are starring on GitHub without stalking their profiles.",
	  "schema": {
	    "/feed/{username}": {
	      "method": "GET",
	      "description": "Get the RSS feed for a single GitHub user.",
	      "params": {
	        "username": "string (GitHub username)"
	      }
	    },
	    "/feeds/{usernames}": {
	      "method": "GET",
	      "description": "Get the RSS feed for multiple GitHub users (separate usernames with '+').",
	      "params": {
	        "usernames": "string (e.g. user1+user2+user3)"
	      }
	    }
	  }
	}`)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
