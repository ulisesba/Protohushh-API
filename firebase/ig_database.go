package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
	"protohush"
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

// SaveFollowers saves a slice of IGDataFollowers to Firebase Realtime Database under the "followers" node.
// Each IGDataFollowers entry is pushed with a unique key.
func (d *Database) SaveFollowers(data []protohush.IGDataFollowers) error {
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

// SaveFollowings saves a slice of IGDataFollowing to Firebase Realtime Database under the "followings" node.
// Each IGDataFollowing entry is pushed with a unique key.
func (d *Database) SaveFollowings(data []protohush.IGDataFollowing) error {
	ref := d.Client.NewRef("followings") // Reference to the "followings" node in the Firebase Realtime Database.
	for _, igData := range data {
		_, err := ref.Push(d.Ctx, &igData)
		if err != nil {
			return err // Return an error if pushing to Firebase fails.
		}
	}
	return nil // Return nil if all takeout was pushed successfully.
}
