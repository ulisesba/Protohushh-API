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
func (c Chat) Search(instruction string) (*protohush.ChatResponse, error) {
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

func (c *Chat) handleDatabaseTasks(resp *Response) (*protohush.ChatResponse, error) {
	var details protohush.ChatQuery
	err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &details)
	if err != nil {
		return nil, err
	}

	var chatResp protohush.ChatResponse
	chatResp.Intention = details.Intention

	switch details.Intention {
	case "FindAllLikes":
		likes, err := c.IGDatabase.FindAllLikes()
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving all likes: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", likes)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindLikesByUsername":
		likes, err := c.IGDatabase.FindLikesByUsername(details.Value)
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving likes by username: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", likes)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindLikesSortedByDate":
		likes, err := c.IGDatabase.FindLikesSortedByDate(details.Limit)
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving likes sorted by date: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", likes)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindAllFollowers":
		followers, err := c.IGDatabase.FindAllFollowers()
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving all followers: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", followers)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindFollowersByUsername":
		followers, err := c.IGDatabase.FindFollowersByUsername(details.Value)
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving followers by username: " + err.Error()
		} else {
			log.Println(followers)
			chatResp.Data = fmt.Sprintf("%v", followers)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindAllFollowings":
		followings, err := c.IGDatabase.FindAllFollowings()
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving all followings: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", followings)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	case "FindFollowingsByUsername":
		followings, err := c.IGDatabase.FindFollowingsByUsername(details.Value)
		if err != nil {
			chatResp.Data = ""
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply + " Error retrieving followings by username: " + err.Error()
		} else {
			chatResp.Data = fmt.Sprintf("%v", followings)
			composedReply, _ := c.composeReply(&chatResp)
			chatResp.Data = composedReply
		}

	default:
		chatResp.Data = ""
		composedReply, _ := c.composeReply(&chatResp)
		chatResp.Data = composedReply + " Unknown intention."
	}

	return &chatResp, nil
}

func (c *Chat) composeReply(chatResp *protohush.ChatResponse) (string, error) {
	userPrompt := fmt.Sprintf("As a customer support bot, I've been tasked with answering a user's question based on the intention '%s' and the provided data '%s'. What's the most concise and user-friendly way to respond?", chatResp.Intention, chatResp.Data)

	switch chatResp.Intention {
	case "FindAllLikes":
		userPrompt = fmt.Sprintf("How many likes have been recorded in the data: '%s'?", chatResp.Data)

	case "FindLikesByUsername":
		userPrompt = fmt.Sprintf("Is there a 'like' from the username specified in the data: '%s'?", chatResp.Data)

	case "FindAllFollowers":
		userPrompt = fmt.Sprintf("Who are the followers listed in the provided data: '%s'?", chatResp.Data)

	case "FindAllFollowings":
		userPrompt = fmt.Sprintf("Who is the primary user following according to the data: '%s'?", chatResp.Data)

	case "FindLikesSortedByDate":
		userPrompt = fmt.Sprintf("What are the 'likes' sorted by date in the data: '%s'?", chatResp.Data)

	case "FindFollowersByUsername":
		userPrompt = fmt.Sprintf("Does exists the username on the data?: '%s'?", chatResp.Data)

	case "FindFollowingsByUsername":
		userPrompt = fmt.Sprintf("Is the username specified following anyone in the data: '%s'?", chatResp.Data)

	default:
		userPrompt = fmt.Sprintf("Based on the intention '%s' and the data '%s', what would be the most appropriate response?", chatResp.Intention, chatResp.Data)
	}
	// System instruction specifically designed for the compose reply task.
	const systemInstructionCompose = "Compose a user-friendly response based on the provided intention and data. The response should be clear, concise, and informative, adhering to best practices for user engagement."
	// Define the message structure for the OpenAI API request.
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	const model = "gpt-3.5-turbo"

	// Separately prepare OpenAI request data for composeReply
	messages := []Message{
		{Role: "system", Content: systemInstructionCompose},
		{Role: "user", Content: userPrompt},
	}
	data := map[string]interface{}{
		"model":    model, // assuming model is defined elsewhere or passed in
		"messages": messages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error preparing request data: %v", err)
	}

	body, err := c.sendRequestToOpenAI(jsonData)
	if err != nil {
		return "", err
	}

	openaiResp, err := c.parseOpenAIResponse(body)
	if err != nil {
		return "", err
	}

	// We're making the assumption here that OpenAI will always return at least one choice in the response.
	return openaiResp.Choices[0].Message.Content, nil
}
