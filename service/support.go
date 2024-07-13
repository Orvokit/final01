package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func SupportReport(uri string) ([]SupportData, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("get %s request failed with error: %s\n", uri, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("support response failed with status code: %d\n", resp.StatusCode)
	}

	result := make([]SupportData, 0)

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("support decode failed with error: %s\n", err.Error())
	}

	return result, nil
}
