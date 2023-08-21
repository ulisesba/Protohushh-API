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

type Following struct {
	Username string `json:"username"`    // Following's username.
	URL      string `json:"profile_uri"` // Following's Instagram profile URL.
}

type Like struct {
	Username string `json:"username"` // Like username.
}

// IGDataFollowing represents Instagram users that the primary user follows.
type IGDataFollowing IGTakeOutData

// IGDataFollowers represents Instagram users who follow the primary user.
type IGDataFollowers IGTakeOutData

// IGDatabase is an interface that provides methods to interact with Instagram data.
type IGDatabase interface {
	// SaveFollowers Store follower details.
	SaveFollowers(followers IGDataFollowers)

	// SaveFollowings Store follower details.
	SaveFollowings(following IGDataFollowing) // Store details of users the primary user is following.

	// FindAllLikes Retrieve all likes.
	FindAllLikes() []Like

	// FindLikesByUsername Retrieve a "like" from a specific user.
	FindLikesByUsername(username string) Like

	// FindAllFollowers Retrieve all followers.
	FindAllFollowers() []Follower

	// FindAllFollowings Retrieve all the people the primary user is following.
	FindAllFollowings() []Following

	// FindLikesSortedByDate Retrieve 'likes' sorted by date.
	FindLikesSortedByDate(limit int) []Like

	// FindFollowersByUsername Retrieve followers by username.
	FindFollowersByUsername(username string) []Follower

	// FindFollowingsByUsername Retrieve people you're following by username.
	FindFollowingsByUsername(username string) []Following
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
