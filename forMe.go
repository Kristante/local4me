package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Структура для всех заявок
type Requests struct {
	ID        int       `json:"id"`
	Subject   string    `json:"subject"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	Member    struct {
		Name string `json:"name"`
	} `json:"member"`
	Team struct {
		Name string `json:"name"`
	} `json:"team"`
}

// Структура для конкретной заявки
type Request struct {
	ID      int    `json:"id"`
	Subject string `json:"subject"`
	Member  struct {
		Name string `json:"name"`
	} `json:"member"`
	CustomFields []CustomField `json:"custom_fields"`
	CreatedBy    CreatedBy     `json:"created_by"`
	ReopenCount  int           `json:"reopen_count"`
}

type CustomField struct {
	ID    string          `json:"id"`
	Value json.RawMessage `json:"value"`
}

type CreatedBy struct {
	Name string `json:"name"`
}

type Note struct {
	ID          int          `json:"id"`
	Person      Person       `json:"person"`
	CreatedAt   string       `json:"created_at"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	Medium      string       `json:"medium"`
	Internal    bool         `json:"internal"`
	Account     struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	NodeID string `json:"nodeID"`
}

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	SourceID string `json:"sourceID,omitempty"`
	NodeID   string `json:"nodeID"`
}

type Attachment struct {
	CreatedAt string `json:"created_at"`
	ID        int    `json:"id"`
	Inline    bool   `json:"inline"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Size      int    `json:"size"`
	URI       string `json:"uri"`
	NodeID    string `json:"nodeID"`
	Note      struct {
		ID     int    `json:"id"`
		NodeID string `json:"nodeID"`
	} `json:"note"`
}

const (
	API_URL_OPEN    = "https://api.itsm.mos.ru/v1/requests/open"
	API_URL_REQUEST = "https://api.itsm.mos.ru/v1/requests/"
)

func getAllRequests(token4Me string) (error, []Requests) {
	if token4Me == "" {
		return errors.New("не был установлен токен авторизации 4me"), nil
	}

	client := http.Client{}

	request, err := http.NewRequest("GET", API_URL_OPEN, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка при создании запроса: %v", err)), nil
	}

	request.Header.Add("Authorization", "Bearer "+token4Me)
	request.Header.Add("x-4me-account", "rpa")

	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка при выполнении запроса: %v", err)), nil
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка при чтении тела ответа: %v", err)), nil
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ошибка %s, тело ответа: %v", response.Status, string(body))), nil
	}

	var requests []Requests
	if err := json.Unmarshal(body, &requests); err != nil {
		return errors.New(fmt.Sprintf("ошибка парсинга JSON: %v", err)), nil
	}

	return nil, requests
}

func getInfoForRequest(requestID int, apiToken string) (error, *Request) {
	apiURL := API_URL_REQUEST + strconv.Itoa(requestID)
	client := http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка создания запроса: %v", err)), nil
	}

	request.Header.Add("Authorization", "Bearer "+apiToken)
	request.Header.Add("x-4me-account", "rpa")

	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка выполнения запроса: %v", err)), nil
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка чтения тела ответа: %v", err)), nil
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ошибка: %s. Тело ответа: %s", response.Status, string(body))), nil
	}

	var requests Request
	if err := json.Unmarshal(body, &requests); err != nil {
		return errors.New(fmt.Sprintf("ошибка парсинга JSON: %v", err)), nil
	}

	return nil, &requests
}

func GetNotesForRequest(requestID int, apiToken string) (error, []Note) {
	apiURL := API_URL_REQUEST + strconv.Itoa(requestID) + "/notes"
	client := http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка создания запроса: %v", err)), nil
	}

	request.Header.Add("Authorization", "Bearer "+apiToken)
	request.Header.Add("x-4me-account", "rpa")

	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка выполнения запроса: %v", err)), nil
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("ошибка чтения тела ответа: %v", err)), nil
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ошибка: %s. Тело ответа: %s", response.Status, string(body))), nil
	}

	var notes []Note

	if err := json.Unmarshal(body, &notes); err != nil {
		return errors.New(fmt.Sprintf("ошибка парсинга JSON: %v", err)), nil
	}

	return nil, notes
}
