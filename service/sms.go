package service

import (
	"bufio"
	"os"
	. "skillbox/go/basic/diplom2/reference"
	"strings"
)

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func SMSReport(uri string) ([]SMSData, error) {
	file, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]SMSData, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")
		if len(words) != 4 {
			continue
		}

		sms := SMSData{
			Country:      words[0],
			Bandwidth:    words[1],
			ResponseTime: words[2],
			Provider:     words[3],
		}

		if _, found := CountryByAlpha2(sms.Country); !found {
			continue
		}

		if found := IsSmsProvider(sms.Provider); !found {
			continue
		}

		result = append(result, sms)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
