package models

import (
	"encoding/json"
	"encoding/xml"
)

type CreditCredits struct {
	XMLName              xml.Name        `xml:"Credit"`
	CustomerID           string          `xml:"CustomerID,attr"`
	CustomerName         string          `xml:"CustomerName,attr"`
	MerchantCustomerCode string          `xml:"MerchantCustomerCode,attr"`
	Amount               string          `xml:"Amount,attr"`
	DomainID             string          `xml:"DomainID,attr"`
	Purchases            CreditPurchases `xml:"Purchases"`
}

type CreditPurchases struct {
	XMLName  xml.Name       `xml:"Purchases"`
	Purchase CreditPurchase `xml:"Purchase"`
}

type CreditPurchase struct {
	XMLName    xml.Name         `xml:"Purchase"`
	ReserveID  string           `xml:"ReserveID,attr"`
	PurchaseID string           `xml:"PurchaseID,attr"`
	SeqNum     string           `xml:"seq_num,attr"`
	Purchase   CreditSelections `xml:"Selections"`
}

type CreditSelections struct {
	XMLName  xml.Name        `xml:"Selections"`
	Purchase CreditSelection `xml:"Selection"`
}

type CreditSelection struct {
	XMLName xml.Name `xml:"Selection"`
	LineID  string   `xml:"LineID,attr"`
}

func ParseCreditCustomer(xmlBody string) []byte {

	var bets CreditCredits

	xml.Unmarshal([]byte(xmlBody), &bets)

	jsonData, _ := json.Marshal(bets)
	// fmt.Println(string(jsonData))
	return jsonData
}
