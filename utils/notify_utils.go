package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func MakeRESTRequest(method, url string, body interface{}, headers map[string]string) (int, map[string]interface{}, error) {
	var requestBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			log.Println("failed to marshal body:", err)
			return 0, nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Notification sending request failed : ", err)
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return resp.StatusCode, nil, fmt.Errorf("failed to decode response: %w", err)
	}
	prettyJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		log.Println("failed to format response JSON:", err)
	} else {
		log.Println("Response Pretty:\n", string(prettyJSON))
	}
	if errMsg, ok := responseData["error"].(string); ok && errMsg != "" {
		return resp.StatusCode, responseData, fmt.Errorf("API error: %s", errMsg)
	}
	return resp.StatusCode, responseData, nil
}

func ExtractRecipientID(body map[string]interface{}) (uuid.UUID, error) {
	content, ok := body["content"].(map[string]interface{})
	if !ok {
		return uuid.Nil, errors.New("invalid or missing 'content' field")
	}
	notifIDStr, ok := content["recipient_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("'recipient_id' not found or not a string")
	}
	notifUUID, err := uuid.Parse(notifIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format: " + err.Error())
	}
	return notifUUID, nil
}
