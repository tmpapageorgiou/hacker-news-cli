package main

import (
	"fmt"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	apiVersion = "v0"
	topStories = "topstories.json"
	timeout    = 5 * time.Second
)

type HackerNewsAPI struct {
	host    string
	client  JSONGetter
	version string
}

func NewHackerNewsAPI(baseURL string) *HackerNewsAPI {
	return &HackerNewsAPI{
		host:    baseURL,
		client:  NewHTTPJSONClient(),
		version: apiVersion,
	}
}

func (h *HackerNewsAPI) TopStories(limit uint8) ([]Story, error) {
	ids, err := h.topStoriesIDs()
	if err != nil {
		return nil, err
	}

	if uint8(len(ids)) < limit {
		limit = uint8(len(ids))
	}

	return h.stories(ids[:limit])
}

func (h *HackerNewsAPI) stories(ids []int64) ([]Story, error) {

	storyChan, errChan := h.storiesBatch(ids)
	defer close(storyChan)
	defer close(errChan)

	stories := make([]Story, 0, len(ids))
	errs := []error{}
	for i := 0; i < len(ids); i++ {
		select {
		case story := <-storyChan:
			stories = append(stories, story)
		case err := <-errChan:
			errs = append(errs, err)
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("failed to get stories: %+v", errs)
	}

	return stories, err
}

func (h *HackerNewsAPI) storiesBatch(ids []int64) (chan Story, chan error) {
	storyChan := make(chan Story, len(ids))
	errChan := make(chan error)

	for _, id := range ids {
		go func(id int64) {
			story, err := h.Story(id)
			if err != nil {
				errChan <- err
			} else {
				storyChan <- story
			}
		}(id)
	}

	return storyChan, errChan
}

func (h *HackerNewsAPI) Story(id int64) (Story, error) {

	uri := fmt.Sprintf("https://%s", path.Join(h.host, h.version, "item", fmt.Sprintf("%d.json", id)))

	log.WithField("uri", uri).Debug("Request story.")

	story := Story{}
	err := h.client.JSONGet(uri, &story)
	return story, err
}

func (h *HackerNewsAPI) topStoriesIDs() ([]int64, error) {

	uri := fmt.Sprintf("https://%s", path.Join(h.host, h.version, topStories))

	log.WithField("uri", uri).Debug("Request top stories ids.")
	ids := []int64{}
	err := h.client.JSONGet(uri, &ids)

	return ids, err
}

type Slicer interface {
	Slice() []string
}

type Story struct {
	ID          int64   `json:"id"`
	Deleted     string  `json:"deleted"`
	Type        string  `json:"type"`
	By          string  `json:"by"`
	Time        int64   `json:"time"`
	Text        string  `json:"text"`
	Dead        string  `json:"dead"`
	Parent      string  `json:"parent"`
	Poll        string  `json:"poll"`
	Kids        []int64 `json:"kids"`
	URL         string  `json:"url"`
	Score       int64   `json:"score"`
	Title       string  `json:"title"`
	Parts       string  `json:"parts"`
	Descendants int64   `json:"descendants"`
}

func Headers() []string {
	return []string{"title", "text", "url"}
}

func (s Story) Slice() []string {
	return []string{s.Title, s.Text, s.URL}
}
