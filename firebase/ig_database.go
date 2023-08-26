package firebase

import (
	"context"
	"log"
	"protohush"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

// Database wraps the Firebase client and context, offering methods to interact with the Firebase Realtime Database.
type Database struct {
	Client *db.Client      // Firebase Realtime Database client.
	Ctx    context.Context // Context used for Firebase operations.
}

func NewIgDatabase(ctx context.Context, dbURL, credPath string) (*Database, error) {

	conf := &firebase.Config{
		DatabaseURL: dbURL, // Set the Firebase Realtime Database URL.
	}

	// Use the Firebase admin SDK credentials to authenticate.
	opt := option.WithCredentialsFile(credPath)

	// Initialize a new Firebase app instance.
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return nil, err
	}

	// Get a Firebase Realtime Database client from the app instance.
	client, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}

	// Return a new Database instance with the Firebase client and provided context.
	return &Database{
		Client: client,
		Ctx:    ctx,
	}, nil

}

// DropDatabase delete all existing nodes
func (d *Database) DropDatabase() error {
	// Reference to the root node in the Firebase Realtime Database (want to delete all)
	ref := d.Client.NewRef("/")

	// Set the value at the root node to nil, effectively deleting all children nodes.
	if err := ref.Set(d.Ctx, nil); err != nil {
		return err // Return an error if the drop fails
	}
	// Return nil if the operation was successful.
	return nil
}

// SaveFollowers saves a slice of IGTakeOutFollowers to Firebase Realtime Database under the "followers" node.
// Each IGTakeOutFollowers entry is pushed with a unique key.
func (d *Database) SaveFollowers(data []protohush.IGTakeOutFollowers) error {
	ref := d.Client.NewRef("followers") // Reference to the "followers" node in the Firebase Realtime Database.
	followers := make([]protohush.Follower, 0)

	for _, igData := range data {
		for _, v := range igData.StringListData {
			followers = append(followers, protohush.Follower{
				Username: v.Value,
				URL:      v.Href,
			})
		}
	}

	for _, f := range followers {
		_, err := ref.Push(d.Ctx, &f)
		if err != nil {
			return err // Return an error if pushing to Firebase fails.
		}
	}

	return nil // Return nil if all takeout was pushed successfully.
}

// SaveFollowings saves a slice of IGTakeOutFollowing to Firebase Realtime Database under the "followings" node.
// Each IGTakeOutFollowing entry is pushed with a unique key.
func (d *Database) SaveFollowings(data []protohush.IGTakeOutFollowing) error {
	ref := d.Client.NewRef("followings") // Reference to the "followings" node in the Firebase Realtime Database.
	for _, igData := range data {
		_, err := ref.Push(d.Ctx, &igData)
		if err != nil {
			return err // Return an error if pushing to Firebase fails.
		}
	}
	return nil // Return nil if all takeout was pushed successfully.
}

func (d *Database) FindAllLikes() ([]protohush.Like, error) {
	ref := d.Client.NewRef("likes")
	var likes []protohush.Like

	if err := ref.Get(d.Ctx, &likes); err != nil {
		return nil, err
	}
	return likes, nil
}

func (d *Database) FindLikesByUsername(username string) (*protohush.Like, error) {
	ref := d.Client.NewRef("likes").OrderByChild("username").EqualTo(username)
	var like protohush.Like

	if err := ref.Get(d.Ctx, &like); err != nil {
		return nil, err
	}
	return &like, nil
}

func (d *Database) FindAllFollowers() ([]protohush.Follower, error) {
	ref := d.Client.NewRef("followers")

	// Change the type to map
	followerData := make(map[string]protohush.Follower)
	if err := ref.Get(d.Ctx, &followerData); err != nil {
		log.Println("Error fetching followers:", err)
		return nil, err
	}

	// Convert the map to a slice
	var followers []protohush.Follower
	for _, f := range followerData {
		followers = append(followers, f)
	}

	return followers, nil
}

func (d *Database) FindAllFollowings() ([]protohush.Following, error) {
	ref := d.Client.NewRef("followings")
	var followings []protohush.Following

	if err := ref.Get(d.Ctx, &followings); err != nil {
		return nil, err
	}
	return followings, nil
}

func (d *Database) FindLikesSortedByDate(limit int) ([]protohush.Like, error) {
	ref := d.Client.NewRef("likes").OrderByChild("date").LimitToFirst(limit)
	var likes []protohush.Like

	if err := ref.Get(d.Ctx, &likes); err != nil {
		return nil, err
	}
	return likes, nil
}

func (d *Database) FindFollowersByUsername(username string) ([]protohush.Follower, error) {
	ref := d.Client.NewRef("followers")
	q := ref.OrderByChild("username").EqualTo(username)

	followerData := make(map[string]protohush.Follower)
	if err := q.Get(d.Ctx, &followerData); err != nil {

		log.Println(err)
		return nil, err
	}

	var followers []protohush.Follower
	for _, f := range followerData {
		followers = append(followers, f)
	}

	return followers, nil
}

func (d *Database) FindFollowingsByUsername(username string) ([]protohush.Following, error) {
	ref := d.Client.NewRef("followings").OrderByChild("username").EqualTo(username)
	var followings []protohush.Following

	if err := ref.Get(d.Ctx, &followings); err != nil {
		return nil, err
	}
	return followings, nil
}
