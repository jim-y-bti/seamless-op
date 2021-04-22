package models

import (
	"encoding/xml"
)

type ReserveBets struct {
	XMLName       xml.Name   `xml:"Bets"`
	Cust_id       string     `xml:"cust_id,attr"`
	Reserve_id    string     `xml:"reserve_id,attr"`
	Amount        string     `xml:"amount,attr"`
	Currency_code string     `xml:"currency_code,attr"`
	Bet           ReserveBet `xml:"Bet"`
}

type ReserveBet struct {
	XMLName     xml.Name `xml:"Bet"`
	BetID       string   `xml:"BetID,attr"`
	BetTypeID   string   `xml:"BetTypeID,attr"`
	BetTypeName string   `xml:"BetTypeName,attr"`
}

func ParseReserve(xmlBody string) ReserveBets {

	var bets ReserveBets

	xml.Unmarshal([]byte(xmlBody), &bets)

	// jsonData, _ := json.Marshal(bets)
	// fmt.Println(string(jsonData))
	// fmt.Println(string(bets.Cust_id))
	return bets
}
