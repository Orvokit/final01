package service

import (
	"fmt"
	"os"
	"strconv"
)

const (
	CreateCustomerMask int64 = 1 << iota
	PurchaseMask
	PayoutMask
	RecurringMask
	FraudControlMask
	CheckoutPageMask
)

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"reccuring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

func BillingReport(uri string) (*BillingData, error) {

	data, err := os.ReadFile(uri)
	if err != nil {
		return nil, fmt.Errorf("Error when opening file: %s", err)
	}

	mask, err := strconv.ParseInt(string(data), 2, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parse to int: %s", err)
	}

	billing := &BillingData{
		CreateCustomer: mask&CreateCustomerMask != 0,
		Purchase:       mask&PurchaseMask != 0,
		Payout:         mask&PayoutMask != 0,
		Recurring:      mask&RecurringMask != 0,
		FraudControl:   mask&FraudControlMask != 0,
		CheckoutPage:   mask&CheckoutPageMask != 0,
	}

	return billing, nil
}
