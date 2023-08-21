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

const openaiURL = "https://api.openai.com/v1/chat/completions"

const openaiToken = "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"

type Response struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Chat struct{}

func (c *Chat) Search(instruction string) (*protohush.ChatIntention, error) {

	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	const model = "gpt-3.5-turbo"
	const systemInstruction = "Analyze the user instruction to identify the intention. The intention can be 'FindAllLikes', 'FindLikesByUsername', 'FindAllFollowers', 'FindAllFollowings', 'FindLikesByDate', 'FindFollowersByUsername', or 'FindFollowingsByUsername'. If a specified collection like 'users' is not identified, provide suggestions for other valid collections such as 'followers', 'followings', and 'likes'. Offer guidance on how to access the relevant user data in these alternatives. Return the response as a JSON object comprising the fields: intention, collection, alternative_collections, and value_to_search."

	const userPromptTemplate = "Instruction: '%s'. If the mentioned collection doesn't exist, suggest where the information might be found within valid Firebase collections, and provide the response in the specified JSON format."

	userPrompt := fmt.Sprintf(userPromptTemplate, instruction)

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
		return nil, err
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(jsonData))
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

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	var openaiResp Response
	err = json.Unmarshal(body, &openaiResp)
	if err != nil {
		return nil, err
	}

	var details protohush.ChatIntention
	err = json.Unmarshal([]byte(openaiResp.Choices[0].Message.Content), &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}
