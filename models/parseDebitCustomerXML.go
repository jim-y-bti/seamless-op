package models

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type DebitCredits struct {
	XMLName              xml.Name       `xml:"Credit"`
	CustomerID           string         `xml:"CustomerID,attr"`
	CustomerName         string         `xml:"CustomerName,attr"`
	MerchantCustomerCode string         `xml:"MerchantCustomerCode,attr"`
	Amount               string         `xml:"Amount,attr"`
	DomainID             string         `xml:"DomainID,attr"`
	Purchases            DebitPurchases `xml:"Purchases"`
}

type DebitPurchases struct {
	XMLName  xml.Name      `xml:"Purchases"`
	Purchase DebitPurchase `xml:"Purchase"`
}

type DebitPurchase struct {
	XMLName    xml.Name        `xml:"Purchase"`
	ReserveID  string          `xml:"ReserveID,attr"`
	PurchaseID string          `xml:"PurchaseID,attr"`
	SeqNum     string          `xml:"seq_num,attr"`
	Purchase   DebitSelections `xml:"Selections"`
}

type DebitSelections struct {
	XMLName  xml.Name       `xml:"Selections"`
	Purchase DebitSelection `xml:"Selection"`
}

type DebitSelection struct {
	XMLName xml.Name `xml:"Selection"`
	LineID  string   `xml:"LineID,attr"`
}

func ParseDebitCustomer(xmlBody string) {

	var bets DebitCredits

	xml.Unmarshal([]byte(xmlBody), &bets)

	jsonData, _ := json.Marshal(bets)
	fmt.Println(string(jsonData))

}
