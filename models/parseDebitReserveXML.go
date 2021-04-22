package models

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type DebitBets struct {
	XMLName    xml.Name `xml:"Bets"`
	Cust_id    string   `xml:"cust_id,attr"`
	Reserve_id string   `xml:"reserve_id,attr"`
	Amount     string   `xml:"amount,attr"`
	Bet        DebitBet `xml:"Bet"`
}

type DebitBet struct {
	XMLName     xml.Name `xml:"Bet"`
	BetID       string   `xml:"BetID,attr"`
	BetTypeID   string   `xml:"BetTypeID,attr"`
	BetTypeName string   `xml:"BetTypeName,attr"`
	LineID      string   `xml:"LineID,attr"`
}

func ParseDebitReserve(xmlBody string) {

	var bets DebitBets

	xml.Unmarshal([]byte(xmlBody), &bets)

	jsonData, _ := json.Marshal(bets)
	fmt.Println(string(jsonData))

}
