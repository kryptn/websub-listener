package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type cache struct {
	Leases map[string]int
	Posts  map[string][]string
}

var internalCache cache

func newCache() cache {
	return cache{
		Leases: map[string]int{},
		Posts:  map[string][]string{},
	}

}

func init() {
	internalCache = newCache()
}

// SetLease will set the subscription lease timer for a given slug
func SetLease(slug string, t string) {

	lease, err := strconv.Atoi(t)
	if err != nil {
		lease = 300
	}

	internalCache.Leases[slug] = lease
}

func find(items []string, val string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}
	return false
}

func addPost(slug string, postId string) {
	postIds, _ := internalCache.Posts[slug]

	p := append(postIds, postId)
	if len(p) > 5 {
		p = p[len(p)-5:]
	}
	internalCache.Posts[slug] = p

}

// ShouldAct Adds the given post Id if it's not already in the list, otherwise returns false
func ShouldAct(slug string, postID string) bool {
	result := true

	if _, ok := internalCache.Posts[slug]; !ok {
		internalCache.Posts[slug] = []string{}
	}
	postIDs, _ := internalCache.Posts[slug]

	if found := find(postIDs, postID); found {
		result = false
	}

	if result {
		addPost(slug, postID)
	}

	return result
}

// CacheStatusHandler will return the cache status
func CacheStatusHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(internalCache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
