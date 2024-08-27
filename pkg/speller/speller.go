package speller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Suggestion struct {
	Word        string   `json:"word"`
	Suggestions []string `json:"suggestions"`
}

func CheckSpelling(text string) ([]Suggestion, error) {
	apiURL := "https://speller.yandex.net/services/spellservice.json/checkText"
	resp, err := http.Get(fmt.Sprintf("%s?text=%s", apiURL, url.QueryEscape(text)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []Suggestion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
