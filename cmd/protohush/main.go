package main

import (
	"context"
	"fmt"
	"log"
	"protohush"
	"protohush/firebase"
	"protohush/openai"
)

const DatabaseUri = "https://protohush-xqqu-default-rtdb.firebaseio.com/"

const CredentialsPath = "./firebase_credentials.json"

func main() {

	chat := &openai.Chat{}
	ask := "give me my last 5 likes"
	details, err := chat.Search(ask)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Intention: %s, Value: %s, Collection: %s Limit: %d", details.Intention, details.Value, details.Collection, details.Limit)
	database, err := firebase.NewIgDatabase(context.Background(), DatabaseUri, CredentialsPath)
	if err != nil {
		return
	}

	followers, err := protohush.ReadFollowersFromFile("./takeout/ig/followers.json")

	log.Println(len(followers), err)

	database.SaveFollowers(followers)

	log.Println(err)
}
