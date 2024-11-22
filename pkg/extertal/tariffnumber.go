package external

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TariffNumberResponse struct {
	Query       string       `json:"query"`
	Year        string       `json:"year"`
	Lang        string       `json:"lang"`
	Version     string       `json:"version"`
	Total       int          `json:"total"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

var TariffNumberUrl = "https://www.tariffnumber.com/api/v1/cnSuggest"

func GetTariffNumber(code string) (TariffNumberResponse, error) {
	//TODO: Добавить контекст и ретраи
	url := fmt.Sprintf("%s?term=%s&lang=en&year=2024", TariffNumberUrl, code)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("cannot get request to tariffnumber url, err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("tariffnumber response != 200, status code is: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read response from tariffnumber, err: %v", err)
	}

	var result TariffNumberResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("cannot unmarshal response from tariffnumber, err: %v", err)
	}

	return result, nil
}
