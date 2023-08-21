package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"protohush"
)

const (
	openaiURL   = "https://api.openai.com/v1/chat/completions"
	openaiToken = "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"
)

type Response struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Chat struct {
	IGDatabase protohush.IGDatabase
}

func (c *Chat) Search(instruction string) (*protohush.ChatQuery, error) {
	userPrompt := c.generateUserPrompt(instruction)
	data := c.prepareOpenAIRequestData(userPrompt)

	body, err := c.sendRequestToOpenAI(data)
	if err != nil {
		return nil, err
	}

	openaiResp, err := c.parseOpenAIResponse(body)
	if err != nil {
		return nil, err
	}

	return c.handleDatabaseTasks(openaiResp)
}

func (c *Chat) generateUserPrompt(instruction string) string {
	const userPromptTemplate = "Instruction: '%s'. If the mentioned collection doesn't exist, suggest where the information might be found within valid Firebase collections, and provide the response in the specified JSON format."
	return fmt.Sprintf(userPromptTemplate, instruction)
}

func (c *Chat) prepareOpenAIRequestData(userPrompt string) []byte {
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	const model = "gpt-3.5-turbo"
	const systemInstruction = "Analyze the user instruction to identify the intention and any potential record limit. The intention can be 'FindAllLikes', 'FindLikeByUsername', 'FindAllFollowers', 'FindAllFollowings', 'FindLikesSortedByDate', 'FindFollowersByUsername', or 'FindFollowingsByUsername'. If a specified collection like 'users' is not identified, provide suggestions for other valid collections such as 'followers', 'followings', and 'likes'. Offer guidance on how to access the relevant user data in these alternatives. If a limit is specified, return up to that number of records. Return the response as a JSON object comprising the fields: intention, collection, alternative_collections, value_to_search, and limit."

	messages := []Message{
		{Role: "system", Content: systemInstruction},
		{Role: "user", Content: userPrompt},
	}

	data := map[string]interface{}{
		"model":    model,
		"messages": messages,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error preparing request data: %v", err)
	}

	return jsonData
}

func (c *Chat) sendRequestToOpenAI(data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))

	return body, nil
}

func (c *Chat) parseOpenAIResponse(body []byte) (*Response, error) {
	var resp Response
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Chat) handleDatabaseTasks(resp *Response) (*protohush.ChatQuery, error) {
	var details protohush.ChatQuery
	err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &details)
	if err != nil {
		return nil, err
	}
	switch details.Intention {

	case "FindAllLikes":
		// Code to fetch all likes
		c.IGDatabase.FindAllLikes()
	case "FindLikesByUsername":
		// Code to fetch a "like" from a specific user
		c.IGDatabase.FindLikesByUsername(details.Value)
	case "FindAllFollowers":
		// Code to fetch all followers
		c.IGDatabase.FindAllFollowers()
	case "FindAllFollowings":
		// Code to fetch all the people you're following
		c.IGDatabase.FindAllFollowings()
	case "FindLikesSortedByDate":
		// Code to find 'likes' by date
		c.IGDatabase.FindLikesSortedByDate(details.Limit)
	case "FindFollowersByUsername":
		// Code to find followers by username
		c.IGDatabase.FindFollowersByUsername(details.Value)
	case "FindFollowingsByUsername":
		// Code to find people you're following by username
		c.IGDatabase.FindFollowingsByUsername(details.Value)
	default:
		// Handle unrecognized intentions or provide a default response
	}

	return &details, nil
}
