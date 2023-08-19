package main

import (
	"context"
	"log"
	"protohush"
	"protohush/firebase"
)

const DatabaseUri = "https://protohush-xqqu-default-rtdb.firebaseio.com/"

const CredentialsPath = "./credentials.json"

func main() {

	database, err := firebase.NewIgDatabase(context.Background(), DatabaseUri, CredentialsPath)
	if err != nil {
		return
	}

	followers, err := protohush.ReadFollowersFromFile("./takeout/ig/followers.json")

	log.Println(len(followers), err)

	database.SaveFollowers(followers)

	log.Println(err)
}
