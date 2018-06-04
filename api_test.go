package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTopStoriesIDs(t *testing.T) {

	tests := []struct {
		name            string
		ids             string
		stories         []string
		expectedStories []Story
		limit           uint8
	}{
		{
			name:            "Canonical",
			ids:             "[123, 456, 789]",
			stories:         []string{`{"id": 123, "title": "test 1"}`, `{"id": 456, "title": "test 2"}`, `{"id": 789, "title": "test 3"}`},
			expectedStories: []Story{Story{ID: 123, Title: "test 1"}, Story{ID: 456, Title: "test 2"}, Story{ID: 789, Title: "test 3"}},
			limit:           3,
		},
		{
			name:            "Empty list",
			ids:             "[]",
			stories:         []string{},
			expectedStories: []Story{},
			limit:           3,
		},
		{
			name:            "Limit less than available ids",
			ids:             "[123, 456, 789]",
			stories:         []string{`{"id": 123, "title": "test 1"}`, `{"id": 456, "title": "test 2"}`, `{"id": 789, "title": "test 3"}`},
			expectedStories: []Story{Story{ID: 123, Title: "test 1"}},
			limit:           1,
		},
		{
			name:            "Limit more than available ids",
			ids:             "[123, 456, 789]",
			stories:         []string{`{"id": 123, "title": "test 1"}`, `{"id": 456, "title": "test 2"}`, `{"id": 789, "title": "test 3"}`},
			expectedStories: []Story{Story{ID: 123, Title: "test 1"}, Story{ID: 456, Title: "test 2"}, Story{ID: 789, Title: "test 3"}},
			limit:           4,
		},
	}

	for _, test := range tests {
		Convey(fmt.Sprintf("Given the %s test case", test.name), t, func() {
			mock := NewMockJSONGet()
			mock.Add("https://unittest.com/v0/topstories.json", test.ids, nil)
			for i, expected := range test.expectedStories {
				mock.Add(fmt.Sprintf("https://unittest.com/v0/item/%d.json", expected.ID), test.stories[i], nil)
			}

			client := HackerNewsAPI{
				host:    "unittest.com",
				client:  mock,
				version: "v0",
			}

			stories, err := client.TopStories(test.limit)
			So(err, ShouldBeNil)

			for _, expected := range test.expectedStories {
				So(stories, ShouldContain, expected)
			}
		})
	}
}
