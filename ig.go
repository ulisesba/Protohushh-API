package protohush

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// IGTakeOutFollowerDetail represents an Instagram follower's details extracted from a takeout.
type IGTakeOutFollowerDetail struct {
	Href      string `json:"href"`      // URL to the follower's Instagram profile.
	Value     string `json:"value"`     // IGTakeOutFollowerDetail's Instagram username.
	Timestamp int64  `json:"timestamp"` // Timestamp when the follower was added.
}

type IGTakeOutLikeDetail struct {
	Href      string `json:"href"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

type IGTestLikes struct {
	Htest string
}

// IGTakeOutFollower contains structured information from an Instagram takeout.
type IGTakeOutFollower struct {
	Title          string                    `json:"title"`            // Title of the takeout block, if provided.
	MediaListData  []string                  `json:"media_list_data"`  // List of media associated with the user (e.g., photos, videos).
	StringListData []IGTakeOutFollowerDetail `json:"string_list_data"` // List of follower details.
}

type IGTakeOutLike struct {
	Title          string                `json:"title"`
	StringListData []IGTakeOutLikeDetail `json:"string_list_data"`
}

// IGTakeOutFollowing represents Instagram users that the primary user follows.
type IGTakeOutFollowing IGTakeOutFollower

// IGTakeOutFollowers represents Instagram users who follow the primary user.
type IGTakeOutFollowers IGTakeOutFollower

type IGTakeOutLikes IGTakeOutLike

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
	Href     string `json:"href"`     // Like href
}

type IGDatabase interface {
	// DropDatabase Store follower details.
	DropDatabase() error

	// SaveFollowers Store follower details.
	SaveFollowers(followers []IGTakeOutFollowers) error

	// SaveFollowings Store follower details.
	SaveFollowings(following []IGTakeOutFollowing) error

	SaveLikes(likes []IGTakeOutLikes) error

	// FindAllLikes Retrieve all likes.
	FindAllLikes() ([]Like, error)

	// FindLikesByUsername Retrieve a "like" from a specific user.
	FindLikesByUsername(username string) (*Like, error)

	// FindAllFollowers Retrieve all followers.
	FindAllFollowers() ([]Follower, error)

	// FindAllFollowings Retrieve all the people the primary user is following.
	FindAllFollowings() ([]Following, error)

	// FindLikesSortedByDate Retrieve 'likes' sorted by date.
	FindLikesSortedByDate(limit int) ([]Like, error)

	// FindFollowersByUsername Retrieve followers by username.
	FindFollowersByUsername(username string) ([]Follower, error)

	// FindFollowingsByUsername Retrieve people you're following by username.
	FindFollowingsByUsername(username string) ([]Following, error)
}

// IGTakeOutProvider is an interface that provides methods to retrieve Instagram follower and following data.
type IGTakeOutProvider interface {
	Followers() []IGTakeOutFollowers  // Retrieve follower details.
	Followings() []IGTakeOutFollowing // Retrieve details of users the primary user is following.
}

// ReadFollowersFromFile reads a file containing Instagram follower details and returns a slice of IGTakeOutFollowers.
func ReadFollowersFromFile(filename string) ([]IGTakeOutFollowers, error) {
	// Read file content.
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data []IGTakeOutFollowers
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadFollowingsFromFile reads a file containing details of users the primary user follows
// and returns a slice of IGTakeOutFollowing.
func ReadFollowingsFromFile(filename string) ([]IGTakeOutFollowing, error) {
	// Read file content.
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Define a wrapper for the following.json format.
	type DataWrapper struct {
		RelationshipsFollowing []IGTakeOutFollowing `json:"relationships_following"`
	}

	var wrappedData DataWrapper
	err = json.Unmarshal(fileBytes, &wrappedData)
	if err != nil {
		return nil, err
	}

	return wrappedData.RelationshipsFollowing, nil
}

func ReadLikesFromFile(filename string) ([]IGTakeOutLikes, error) {
	// Read file content.
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	type LikesDataWrapper struct {
		Likes []IGTakeOutLikes `json:"likes_media_likes"`
	}

	var likesDataWrapper LikesDataWrapper
	err = json.Unmarshal(fileBytes, &likesDataWrapper)
	if err != nil {
		return nil, err
	}

	log.Println(likesDataWrapper.Likes)

	return likesDataWrapper.Likes, nil
}
