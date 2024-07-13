package service

import (
	"bufio"
	"os"
	. "skillbox/go/basic/diplom2/reference"
	"strconv"
	"strings"
)

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

func EmailReport(uri string) ([]EmailData, error) {

	file, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]EmailData, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")
		if len(words) != 3 {
			continue
		}

		var email EmailData

		email.Country = words[0]

		email.Provider = words[1]

		time, err := strconv.Atoi(words[2])
		if err != nil {
			continue
		}
		email.DeliveryTime = time

		if _, found := CountryByAlpha2(email.Country); !found {
			continue
		}

		if found := IsEmailProvider(email.Provider); !found {
			continue
		}

		result = append(result, email)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
