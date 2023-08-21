package protohush

import (
	"encoding/json"
	"io/ioutil"
)

// IGTakeOutFollower represents an Instagram follower's details extracted from a takeout.
type IGTakeOutFollower struct {
	Href      string `json:"href"`      // URL to the follower's Instagram profile.
	Value     string `json:"value"`     // IGTakeOutFollower's Instagram username.
	Timestamp int64  `json:"timestamp"` // Timestamp when the follower was added.
}

// IGTakeOutData contains structured information from an Instagram takeout.
type IGTakeOutData struct {
	Title          string              `json:"title"`            // Title of the takeout block, if provided.
	MediaListData  []string            `json:"media_list_data"`  // List of media associated with the user (e.g., photos, videos).
	StringListData []IGTakeOutFollower `json:"string_list_data"` // List of follower details.
}

// Follower represents detailed information about an Instagram follower.
type Follower struct {
	Username string `json:"username"`    // Follower's username.
	URL      string `json:"profile_uri"` // Follower's Instagram profile URL.
}

// IGDataFollowing represents Instagram users that the primary user follows.
type IGDataFollowing IGTakeOutData

// IGDataFollowers represents Instagram users who follow the primary user.
type IGDataFollowers IGTakeOutData

// IGDatabase is an interface that provides methods to store Instagram follower and following data.
type IGDatabase interface {
	SaveFollowers(followers IGDataFollowers)  // Store follower details.
	SaveFollowings(following IGDataFollowing) // Store details of users the primary user is following.
}

// IGDataProvider is an interface that provides methods to retrieve Instagram follower and following data.
type IGDataProvider interface {
	Followers() []IGDataFollowers  // Retrieve follower details.
	Followings() []IGDataFollowing // Retrieve details of users the primary user is following.
}

// ReadFollowersFromFile reads a file containing Instagram follower details and returns a slice of IGDataFollowers.
func ReadFollowersFromFile(filename string) ([]IGDataFollowers, error) {
	// Read file content.
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

// ReadFollowingsFromFile reads a file containing details of users the primary user follows
// and returns a slice of IGDataFollowing.
func ReadFollowingsFromFile(filename string) ([]IGDataFollowing, error) {
	// Read file content.
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
