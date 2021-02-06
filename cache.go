package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Cache struct {
	Leases map[string]int
	Posts  map[string][]string

	leaseRenewers map[string]func()
}

func NewCache() *Cache {
	return &Cache{
		Leases: map[string]int{},
		Posts:  map[string][]string{},
	}
}

// func init() {
// 	internalCache = NewCache()
// }

// var internalCache Cache

// SetLease will set the subscription lease timer for a given slug
func (c *Cache) SetLease(slug string, t string) {

	lease, err := strconv.Atoi(t)
	if err != nil {
		lease = 300
	}

	leaseExp := int(time.Now().Unix()) + lease - int(lease/5)
	c.Leases[slug] = leaseExp
}

func find(items []string, val string) bool {
	for _, item := range items {
		if item == val {
			return true
		}
	}
	return false
}

func (c *Cache) addPost(slug string, postId string) {
	postIds, _ := c.Posts[slug]

	p := append(postIds, postId)
	if len(p) > 5 {
		p = p[len(p)-5:]
	}
	c.Posts[slug] = p

}

// ShouldAct Adds the given post Id if it's not already in the list, otherwise returns false
func (c *Cache) ShouldAct(slug string, postID string) bool {
	result := true

	if _, ok := c.Posts[slug]; !ok {
		c.Posts[slug] = []string{}
	}
	postIDs, _ := c.Posts[slug]

	if found := find(postIDs, postID); found {
		result = false
	}

	if result {
		c.addPost(slug, postID)
	}

	return result
}

// CacheStatusHandler will return the cache status
func (c *Cache) CacheStatusHandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
