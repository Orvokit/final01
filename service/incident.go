package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

func IncidentReport(uri string) ([]IncidentData, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("get %s request failed with error: %s\n", uri, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incident response failed with status code: %d\n", resp.StatusCode)
	}

	result := make([]IncidentData, 0)

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("incident decode failed with error: %s\n", err.Error())
	}

	return result, nil
}
