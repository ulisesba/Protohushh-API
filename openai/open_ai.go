package openai

// Import necessary packages.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"protohush"
)

// Define constants for the OpenAI API endpoint and authorization token.
const (
	openaiURL   = "https://api.openai.com/v1/chat/completions"
	openaiToken = "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"
)

// Response structure for OpenAI API responses.
type Response struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Chat structure to store the database object.
type Chat struct {
	IGDatabase protohush.IGDatabase
}

// Search method processes the user instruction and interacts with OpenAI and database.
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

// generateUserPrompt creates a user prompt from the given instruction.
func (c *Chat) generateUserPrompt(instruction string) string {
	const userPromptTemplate = "Instruction: '%s'. If the mentioned collection doesn't exist, suggest where the information might be found within valid Firebase collections, and provide the response in the specified JSON format."
	return fmt.Sprintf(userPromptTemplate, instruction)
}

// prepareOpenAIRequestData prepares data for sending to OpenAI.
func (c *Chat) prepareOpenAIRequestData(userPrompt string) []byte {
	// Define the message structure for the OpenAI API request.
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	const model = "gpt-3.5-turbo"
	// System instruction for OpenAI to process the user prompt correctly.
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

// sendRequestToOpenAI sends a POST request to the OpenAI API.
func (c *Chat) sendRequestToOpenAI(data []byte) ([]byte, error) {
	// Create a new HTTP request.
	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Set necessary headers for the request.
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

// parseOpenAIResponse parses the response received from the OpenAI API.
func (c *Chat) parseOpenAIResponse(body []byte) (*Response, error) {
	var resp Response
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// handleDatabaseTasks handles database operations based on the intention received from OpenAI.
func (c *Chat) handleDatabaseTasks(resp *Response) (*protohush.ChatQuery, error) {
	var details protohush.ChatQuery
	err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &details)
	if err != nil {
		return nil, err
	}

	// Handle database tasks based on the user's intention.
	switch details.Intention {
	// ... (as per your switch cases)
	}

	return &details, nil
}
