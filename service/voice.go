package service

import (
	"bufio"
	"os"
	. "skillbox/go/basic/diplom2/reference"
	"strconv"
	"strings"
)

const (
	VoiceCountry = iota
	VoiceBandwidth
	VoiceResponse
	VoiceProvider
	VoiceConnection
	VoiceTTFB
	VoicePurity
	VoiceMedian
	VoiceStructSize
)

type VoiceCallData struct {
	Country             string  `json:"country"`
	Bandwidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	TTFB                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_calls_time"`
}

func VoiceReport(uri string) ([]VoiceCallData, error) {
	file, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]VoiceCallData, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ";")
		if len(words) != VoiceStructSize {
			continue
		}

		var vcd VoiceCallData

		if country, found := CountryByAlpha2(words[VoiceCountry]); !found {
			continue
		} else {
			vcd.Country = country
		}

		vcd.Bandwidth = words[VoiceBandwidth]

		vcd.ResponseTime = words[VoiceResponse]

		vcd.Provider = words[VoiceProvider]

		connectionStability, err := strconv.ParseFloat(words[VoiceConnection], 32)
		if err != nil {
			continue
		}
		vcd.ConnectionStability = float32(connectionStability)

		ttfb, err := strconv.Atoi(words[VoiceTTFB])
		if err != nil {
			continue
		}
		vcd.TTFB = ttfb

		purity, err := strconv.Atoi(words[VoicePurity])
		if err != nil {
			continue
		}
		vcd.VoicePurity = purity

		nedian, err := strconv.Atoi(words[VoiceMedian])
		if err != nil {
			continue
		}
		vcd.MedianOfCallsTime = nedian

		if found := IsVcdProvider(vcd.Provider); !found {
			continue
		}

		result = append(result, vcd)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
