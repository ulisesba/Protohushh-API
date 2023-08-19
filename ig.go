package protohush

import (
	"encoding/json"
	"io/ioutil"
)

// IGTakeOutFollower represents an Instagram follower's details extracted from a takeout.
type IGTakeOutFollower struct {
	Href      string `json:"href"`      // URL to the follower's Instagram profile.
	Value     string `json:"value"`     // IGTakeOutFollower's Instagram username.
	Timestamp int64  `json:"timestamp"` // Timestamp when the follower takeout was captured.
}

type IGTakeOutData struct {
	Title          string              `json:"title"`            // Title of the takeout block (can be empty based on the takeout source).
	MediaListData  []string            `json:"media_list_data"`  // List of media takeout (photos, videos, etc.).
	StringListData []IGTakeOutFollower `json:"string_list_data"` // List of follower details.
}

// Follower represents an actual ig follower.
type Follower struct {
	Username string `json:"username"`
	URL      string `json:"profile_uri"`
}

type IGDataFollowing IGTakeOutData

type IGDataFollowers IGTakeOutData

type IGDatabase interface {
	SaveFollowers(followers IGDataFollowers)
	SaveFollowings(following IGDataFollowing)
}

type IGDataProvider interface {
	Followers() []IGDataFollowers
	Followings() []IGDataFollowing
}

func ReadFollowersFromFile(filename string) ([]IGDataFollowers, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data []IGDataFollowers
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadFollowingsFromFile reads the following.json file and parses its content into a slice of IGTakeOutData.
func ReadFollowingsFromFile(filename string) ([]IGDataFollowing, error) {
	// Read the file content.
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Define a wrapper for the following.json format.
	type DataWrapper struct {
		RelationshipsFollowing []IGDataFollowing `json:"relationships_following"`
	}

	var wrappedData DataWrapper
	err = json.Unmarshal(fileBytes, &wrappedData)
	if err != nil {
		return nil, err
	}

	return wrappedData.RelationshipsFollowing, nil
}
