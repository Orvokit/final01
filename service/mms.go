package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	. "skillbox/go/basic/diplom2/reference"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func MMSReport(uri string) ([]MMSData, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("get %s request failed with error: %s\n", uri, err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mms response failed with status code: %d\n", resp.StatusCode)
	}

	mmsData := make([]MMSData, 0)

	err = json.NewDecoder(resp.Body).Decode(&mmsData)
	if err != nil {
		return nil, fmt.Errorf("mms decode failed with error: %s\n", err.Error())
	}

	var result []MMSData

	for _, mms := range mmsData {
		if _, found := CountryByAlpha2(mms.Country); !found {
			continue
		}
		if found := IsMmsProvider(mms.Provider); !found {
			continue
		}
		result = append(result, mms)
	}

	return result, nil
}
