package main

import (
	"context"
	"log"
	"protohush"
	"protohush/api"
	"protohush/firebase"
	"protohush/openai"
)

const DatabaseUri = "https://protohush-xqqu-default-rtdb.firebaseio.com/"

const CredentialsPath = "./firebase_credentials.json"

func main() {
	// Initialize the Chat object with database
	var chat openai.Chat

	database, err := firebase.NewIgDatabase(context.Background(), DatabaseUri, CredentialsPath)

	chat.IGDatabase = database // Assuming you have a method to initialize the database.

	if err != nil {
		return
	}

	followers, err := protohush.ReadFollowersFromFile("./takeout/ig/followers.json")
	likes, err := protohush.ReadLikesFromFile("./takeout/ig/likes.json")

	log.Println(len(followers), err)
	log.Println(len(likes), err)

	//	database.DropDatabase()
	//	database.SaveFollowers(followers)
	//	database.SaveLikes(likes)

	api.NewApi(chat).Run()

}
