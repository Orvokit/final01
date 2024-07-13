package main

import (
	"encoding/json"
	"log"
	"net/http"
	. "skillbox/go/basic/diplom2/reference"
	. "skillbox/go/basic/diplom2/service"
	"sort"

	"github.com/gorilla/mux"
)

const (
	smsUri      = "simulator/sms.data"
	mmsUri      = "http://127.0.0.1:8383/mms"
	voiceUri    = "simulator/voice.data"
	emailUri    = "simulator/email.data"
	billingUri  = "simulator/billing.data"
	supportUri  = "http://127.0.0.1:8383/support"
	incidentUri = "http://127.0.0.1:8383/accendent"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api", handleRequest).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	http.Handle("/", router)

	// http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.HandleFunc("/api", handleRequest)
	log.Fatal(http.ListenAndServe(":8282", nil))
}

type ResultT struct {
	// true, если все этапы сбора данных прошли успешно, false во всех остальных случаях
	Status bool `json:"status"`
	// заполнен, если все этапы сбора данных прошли успешно, nil во всех остальных случаях
	Data *ResultSetT `json:"data"`
	// пустая строка если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки (детали ниже)
	Error string `json:"error"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	var resultT ResultT

	if resultSetT, err := getResultData(); err == nil {
		resultT = ResultT{
			Status: true,
			Data:   resultSetT,
			Error:  "",
		}
	} else {
		log.Printf("Error on collect data: %s", err)
		resultT = ResultT{
			Status: false,
			Data:   nil,
			Error:  "Error on collect data",
		}
	}

	response, _ := json.Marshal(resultT)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(response)
}

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}

func getResultData() (*ResultSetT, error) {
	sms, err := SMSReport(smsUri)
	if err != nil {
		return nil, err
	}

	mms, err := MMSReport(mmsUri)
	if err != nil {
		return nil, err
	}

	voice, err := VoiceReport(voiceUri)
	if err != nil {
		return nil, err
	}

	emails, err := EmailReport(emailUri)
	if err != nil {
		return nil, err
	}

	billing, err := BillingReport(billingUri)
	if err != nil {
		return nil, err
	}

	support, err := SupportReport(supportUri)
	if err != nil {
		return nil, err
	}

	incidents, err := IncidentReport(incidentUri)
	if err != nil {
		return nil, err
	}

	resultSetT := &ResultSetT{
		SMS:       smsCollect(sms),
		MMS:       mmsCollect(mms),
		VoiceCall: voice,
		Email:     emailCollect(emails),
		Billing:   *billing,
		Support:   supportCollect(support),
		Incidents: incidentsCollect(incidents),
	}

	return resultSetT, nil
}

func smsCollect(smsData []SMSData) [][]SMSData {
	var sortByProvider []SMSData
	var sortByCountry []SMSData
	for _, sms := range smsData {
		sms.Country, _ = CountryByAlpha2(sms.Country)
		sortByProvider = append(sortByProvider, sms)
		sortByCountry = append(sortByCountry, sms)
	}

	sort.SliceStable(sortByProvider, func(i, j int) bool {
		return sortByProvider[i].Provider < sortByProvider[j].Provider
	})

	sort.SliceStable(sortByCountry, func(i, j int) bool {
		return sortByCountry[i].Country < sortByCountry[j].Country
	})

	var result [][]SMSData
	result = append(result, sortByProvider)
	result = append(result, sortByCountry)

	return result
}

func mmsCollect(mmsData []MMSData) [][]MMSData {
	var sortByProvider []MMSData
	var sortByCountry []MMSData
	for _, mms := range mmsData {
		mms.Country, _ = CountryByAlpha2(mms.Country)
		sortByProvider = append(sortByProvider, mms)
		sortByCountry = append(sortByCountry, mms)
	}

	sort.SliceStable(sortByProvider, func(i, j int) bool {
		return sortByProvider[i].Provider < sortByProvider[j].Provider
	})

	sort.SliceStable(sortByCountry, func(i, j int) bool {
		return sortByCountry[i].Country < sortByCountry[j].Country
	})

	var result [][]MMSData
	result = append(result, sortByProvider)
	result = append(result, sortByCountry)

	return result
}

func emailCollect(emailsReport []EmailData) map[string][][]EmailData {

	emailsByCounty := make(map[string][]EmailData, 0)
	for _, email := range emailsReport {
		countryName, _ := CountryByAlpha2(email.Country)
		emailsByCounty[countryName] = append(emailsByCounty[countryName], email)
	}

	emailsSamplingByCountry := make(map[string][][]EmailData, 0)

	for country, emails := range emailsByCounty {
		sort.SliceStable(emails, func(i, j int) bool {
			return emails[i].DeliveryTime < emails[j].DeliveryTime
		})

		if len(emails) > 3 {
			emailsSamplingByCountry[country] = append(emailsSamplingByCountry[country], emails[:3])
			emailsSamplingByCountry[country] = append(emailsSamplingByCountry[country], emails[len(emails)-3:])
		} else {
			emailsSamplingByCountry[country] = append(emailsSamplingByCountry[country], emails)
			emailsSamplingByCountry[country] = append(emailsSamplingByCountry[country], emails)
		}
	}

	return emailsSamplingByCountry
}

func supportCollect(supportData []SupportData) []int {
	var count int
	for _, support := range supportData {
		count += support.ActiveTickets

	}

	expect := float64(count) * float64(60/18)

	var workload int
	switch {
	case count < 9:
		workload = 1
	case count < 17:
		workload = 2
	default:
		workload = 3
	}

	return []int{workload, int(expect)}
}

func incidentsCollect(incidentData []IncidentData) []IncidentData {
	sort.SliceStable(incidentData, func(i, j int) bool {
		return incidentData[i].Status < incidentData[j].Status
	})

	return incidentData
}
