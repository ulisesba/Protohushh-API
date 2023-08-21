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
	ask := "give me all the likes that  i made to guillermo"
	details, err := chat.Search(ask)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Intention: %s, Value: %s, Collection: %s\n", details.Intention, details.Value, details.Collection)
	database, err := firebase.NewIgDatabase(context.Background(), DatabaseUri, CredentialsPath)
	if err != nil {
		return
	}

	followers, err := protohush.ReadFollowersFromFile("./takeout/ig/followers.json")

	log.Println(len(followers), err)

	database.SaveFollowers(followers)

	log.Println(err)
}
