package api

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"stargazer/pkg/cache"
	"stargazer/pkg/feed"
)

func TestFeedCache(t *testing.T) {
	c := cache.NewFeedCache()

	// Test empty cache
	feedVal, exists := c.Get("testuser")
	assert.False(t, exists)
	assert.Nil(t, feedVal)

	// Test setting and getting from cache
	testFeed := &feed.RSS{
		Version: "2.0",
		Channel: feed.Channel{
			Title: "Test Feed",
		},
	}
	c.Set("testuser", testFeed)

	feedVal, exists = c.Get("testuser")
	assert.True(t, exists)
	assert.Equal(t, testFeed, feedVal)

	// Test cache expiration
	cTest := c
	cTestEntry := cache.CacheEntry{
		Feed:      testFeed,
		Timestamp: time.Now().Add(-16 * time.Minute),
	}
	cTestEntries := cTestEntriesField(cTest)
	cTestEntries["testuser"] = cTestEntry
	feedVal, exists = c.Get("testuser")
	assert.False(t, exists)
	assert.Nil(t, feedVal)
}

func cTestEntriesField(c *cache.FeedCache) map[string]cache.CacheEntry {
	return getUnexportedField(c, "entries").(map[string]cache.CacheEntry)
}

func getUnexportedField(obj interface{}, field string) interface{} {
	v := reflect.ValueOf(obj).Elem().FieldByName(field)
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface()
}

func TestHandleRSSFeed(t *testing.T) {
	if err := os.Setenv("GITHUB_TOKEN", "test-token"); err != nil {
		t.Fatalf("Failed to set GITHUB_TOKEN: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("GITHUB_TOKEN"); err != nil {
			t.Errorf("Failed to unset GITHUB_TOKEN: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/feed/{username}", HandleRSSFeed).Methods("GET")

	tests := []struct {
		name       string
		username   string
		wantStatus int
	}{
		{
			name:       "Valid username",
			username:   "testuser",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/feed/"+tt.username, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/xml", rr.Header().Get("Content-Type"))

				var f feed.RSS
				err := xml.Unmarshal(rr.Body.Bytes(), &f)
				assert.NoError(t, err)
				assert.Equal(t, "2.0", f.Version)
			}
		})
	}
}

func TestHandleMultiUserRSSFeed(t *testing.T) {
	if err := os.Setenv("GITHUB_TOKEN", "test-token"); err != nil {
		t.Fatalf("Failed to set GITHUB_TOKEN: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("GITHUB_TOKEN"); err != nil {
			t.Errorf("Failed to unset GITHUB_TOKEN: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/feeds/{usernames}", HandleMultiUserRSSFeed).Methods("GET")
	r.HandleFunc("/feeds/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No usernames provided", http.StatusBadRequest)
	}).Methods("GET")
	r.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No usernames provided", http.StatusBadRequest)
	}).Methods("GET")

	tests := []struct {
		name       string
		usernames  string
		wantStatus int
	}{
		{
			name:       "Single user",
			usernames:  "testuser",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Multiple users",
			usernames:  "user1+user2+user3",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty usernames",
			usernames:  "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Too many usernames",
			usernames:  strings.Repeat("user+", MaxUsernames) + "user",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.usernames == "" {
				req = httptest.NewRequest("GET", "/feeds/", nil)
			} else {
				req = httptest.NewRequest("GET", "/feeds/"+tt.usernames, nil)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, "application/xml", rr.Header().Get("Content-Type"))

				var f feed.RSS
				err := xml.Unmarshal(rr.Body.Bytes(), &f)
				assert.NoError(t, err)
				assert.Equal(t, "2.0", f.Version)
			}
		})
	}
}
