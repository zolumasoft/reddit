// Package reddit implements a basic client for the Reddit API
package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	redditURL = "https://www.reddit.com/r/"
)

// Item describes a Reddit item
type Item struct {
	Title    string
	URL      string
	Comments int `json:"num_comments"`
}

// Response captures the recent items pulled using the Reddit API
type response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

// String - formats Title, URL and total number of comments per article(item)
func (i Item) String() string {
	com := ""
	switch i.Comments {
	case 0:
		// nothing
	case 1:
		com = " (1 comment)"
	default:
		com = fmt.Sprintf(" (%d comments)", i.Comments)
	}
	return fmt.Sprintf("%s%s\n%s", i.Title, com, i.URL)
}

// Get - fetches the most recent Items posted the specified sub-reddit
func Get(reddit string) ([]Item, error) {
	url := fmt.Sprintf("%s%s.json", redditURL, reddit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	r := new(response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	items := make([]Item, len(r.Data.Children))
	for i, child := range r.Data.Children {
		items[i] = child.Data
	}
	return items, nil
}
