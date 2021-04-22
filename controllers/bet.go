package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"seamlessop/models"
	"seamlessop/statics"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BetReserve(c *gin.Context) {
	fmt.Println("Bet Reserve")

	reserveId := c.DefaultQuery("reserve_id", "{na}")
	amount := c.DefaultQuery("amount", "{na}")
	custId := c.DefaultQuery("cust_id", "{na}")
	strRes := ""

	// Check query parameter (Exists & not empty)
	chkMap := map[string]string{
		"custId":    custId,
		"reserveId": reserveId,
		"amount":    amount,
	}
	strRes = chkURLParams(chkMap)
	if len(strRes) != 0 {
		strRes = genErrResponse("-1", strRes)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// Parse cust_id since it is base64 encode while sending validateToken
	// feed back.
	custDecode, err := base64.RawStdEncoding.DecodeString(custId)
	if err != nil {
		// CustomerNotFound -2
		strRes = genErrResponse("-2", "error_"+string(custDecode))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	tmpAry := strings.Split(string(custDecode), "|")
	cId, err := strconv.Atoi(tmpAry[1])
	fmt.Printf("Cust ID: %s, Reserve ID: %s , Amount: %s\n", tmpAry[1], reserveId, amount)

	// Parse the XML body.
	body, _ := ioutil.ReadAll(c.Request.Body)
	// xmlRes := models.ParseReserve(string(body))

	nAmount, err := strconv.ParseFloat(amount, 64)
	betLogInfo := statics.BetLog{
		// BId:       bId,
		CId:       cId,
		ReserveId: reserveId,
		Endpoint:  "reserve",
		BlType:    "img",
		Amount:    nAmount,
		XmlBody:   string(body),
	}

	// Check if reserve already exist.
	resBody := models.GetBetLogResBody(betLogInfo)
	if len(resBody) > 0 {
		fmt.Println("Reserve Found")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, resBody)
		return
	}

	uRes := models.LockAndUpdateBalance(betLogInfo)
	if len(uRes) == 0 {
		strRes = genErrResponse("-2", "CustomerNotFound")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	tmpFloat, err := strconv.ParseFloat(uRes["newBalance"], 64)
	if tmpFloat < 0 {
		strRes = genErrResponse("-4", "InsufficientFunds")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// error_message=NoErrors
	// trx_id=7d9e39183dc441eea63233baa50103d2
	// balance=30950820.3400
	// BonusUsed=0
	// error_code=0

	// genTrx
	genTrx := betLogInfo.ReserveId + strconv.Itoa(rand.Intn(200))

	// genTrx := base64.RawStdEncoding.EncodeToString([]byte(sEnc))

	genInfo := map[string]string{
		"error_code": "0",
		"error_msg":  "NoError",
		"genTrx":     genTrx,
		"balance":    uRes["newBalance"],
	}
	resStr := genReserveResponse(genInfo)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, resStr)

	// Init betInfo for record
	cwId, err := strconv.Atoi(uRes["cw_id"])
	betInfo := statics.Bets{
		CwId:      cwId,
		ReserveId: reserveId,
		Balance:   tmpFloat,
		Status:    1,
	}

	// fmt.Println("--- Check params here ---")
	// fmt.Println(betInfo)
	// Place a Reserve in bets
	betLogInfo.BId = models.InsertBets(betInfo)

	// Place a Reserve in bet_log
	models.InsertBetLog(betLogInfo)

	return
}

func BetCancelReserve(c *gin.Context) {

	reserveId := c.DefaultQuery("reserve_id", "{na}")
	custId := c.DefaultQuery("cust_id", "{na}")
	strRes := ""

	// Check query parameter (Exists & not empty)
	chkMap := map[string]string{
		"custId":    custId,
		"reserveId": reserveId,
	}
	errMsg := chkURLParams(chkMap)
	if len(errMsg) != 0 {
		strRes = genErrResponse("0", errMsg)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// Parse the XML body.
	body, _ := ioutil.ReadAll(c.Request.Body)
	// xmlRes := models.ParseReserve(string(body))

	// Parse cust_id since it is base64 encode while sending validateToken
	// feed back.
	custDecode, err := base64.RawStdEncoding.DecodeString(custId)
	if err != nil {
		// CustomerNotFound 0
		strRes = genErrResponse("0", "error_"+string(custDecode))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	tmpAry := strings.Split(string(custDecode), "|")
	cId, err := strconv.Atoi(tmpAry[1])

	fmt.Printf("Cust ID: %d, Reserve ID: %s \n", cId, reserveId)

	// Check bet log if there're multiple records
	betInfo := statics.BetLog{CId: cId, ReserveId: reserveId, Endpoint: "cancelReserve"}
	fmt.Println(betInfo)

	// Check if cancel log exist
	cancelRes := models.GetBetLog(betInfo)
	if len(cancelRes) > 0 {
		fmt.Println("---- strRes ----")
		for idx, val := range cancelRes {
			fmt.Printf("Cancel reserve already exist Bid: %d \n", idx)
			fmt.Println(val)
			strRes = val.ResBody
			fmt.Println(strRes)
		}
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	fmt.Println(cancelRes)

	// Check for multiple log
	tmpInfo := statics.BetLog{CId: cId, ReserveId: reserveId}
	blRes := models.GetBetLog(tmpInfo)
	fmt.Println("blRes: --------------------")
	fmt.Println(blRes)

	// Send response
	// ToDo: Put into a function or make a switch
	if len(blRes) == 0 {
		genInfo := map[string]string{
			"error_code": "0",
			"error_msg":  "No record of reserve " + reserveId,
			"balance":    "0",
		}
		strRes = genCancelReserveResponse(genInfo)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
	}
	if len(blRes) > 1 {
		genInfo := map[string]string{
			"error_code": "0",
			"error_msg":  "CriticalErrorMultipleReserveLogFound, reserve " + reserveId,
			"balance":    "0",
		}
		strRes = genCancelReserveResponse(genInfo)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
	}
	if len(blRes) == 1 {
		var fBalance float64
		var bId int
		var sAmt float64

		for idx, val := range blRes {
			fmt.Println(idx)
			sBal := val.Balance
			sAmt := val.Amount
			if err != nil {
				panic(err)
			}
			fBalance = sBal + sAmt
			bId = val.BId
		}

		sBalance := fmt.Sprintf("%f", fBalance)
		genInfo := map[string]string{
			"error_code": "0",
			"error_msg":  "No Error",
			"balance":    sBalance,
		}
		strRes = genCancelReserveResponse(genInfo)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)

		updateInfo := statics.Bets{ReserveId: reserveId, Balance: fBalance, Status: 2}
		updateRes := models.UpdateBet(updateInfo)
		fmt.Println(updateRes)

		betLogInfo := statics.BetLog{
			BId:       bId,
			CId:       cId,
			ReserveId: reserveId,
			ReqId:     "",
			Balance:   fBalance,
			Endpoint:  "cancelReserve",
			XmlBody:   string(body),
			ResBody:   strRes,
			BlType:    "img",
			Amount:    sAmt,
		}
		// Insert a cancel reserve log
		blId := models.InsertBetLog(betLogInfo)
		fmt.Println(blId)

	}

	return
}

func BetDebitReserve(c *gin.Context) {

	fmt.Println("Debit Reserve")
	custId := c.Query("cust_id")
	reserveId := c.Query("reserve_id")
	amount := c.Query("amount")
	reqId := c.Query("req_id")
	strRes := ""

	fmt.Printf("Cust ID: %s, Reserve ID: %s, Amount: %s, RequestID: %s\n", custId, reserveId, amount, reqId)

	// Check query parameter (Exists & not empty)
	chkMap := map[string]string{
		"custId":    custId,
		"reserveId": reserveId,
		"amount":    amount,
		"reqId":     reqId,
	}
	errMsg := chkURLParams(chkMap)
	if len(errMsg) != 0 {
		strRes = genErrResponse("0", errMsg)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	custDecode, err := base64.RawStdEncoding.DecodeString(custId)
	if err != nil {
		strRes = genErrResponse("0", "error_"+string(custDecode))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	tmpAry := strings.Split(string(custDecode), "|")
	cId, err := strconv.Atoi(tmpAry[1])
	fmt.Printf("Client ID: %d, Cust ID: %s, Reserve ID: %s , Amount: %s\n", cId, tmpAry[1], reserveId, amount)

	body, _ := ioutil.ReadAll(c.Request.Body)
	nAmount, err := strconv.ParseFloat(amount, 64)

	// Check for duplicate
	betLogInfo := statics.BetLog{
		// BId:       bId,
		CId:       cId,
		ReserveId: reserveId,
		ReqId:     reqId,
		Endpoint:  "debitReserve",
		BlType:    "img",
		Amount:    nAmount,
		XmlBody:   string(body),
	}
	resBody := models.GetBetLogResBody(betLogInfo)
	if len(resBody) > 0 {
		fmt.Println("Debit Reserve Found")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, resBody)
		return
	}

	// Get reserve amount
	reserveBetLog := statics.BetLog{
		CId:       cId,
		ReserveId: reserveId,
		// Endpoint:  "reserve",
	}
	reserveBody := models.GetBetLog(reserveBetLog)
	if len(reserveBody) < 1 {
		strRes = genErrResponse("0", "No records found! ReserveID:"+reserveId)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	var reserveAmt float64
	var tmpBid int
	var totalAmt float64 = 0

	if len(reserveBody) > 0 {
		for idx, val := range reserveBody {
			fmt.Printf("Idx: %d, Amt: %f, bId: %d, endpoint: %s \n", idx, val.Amount, val.BId, val.Endpoint)
			// if val["endpoint"] == "reserve" {
			if val.Endpoint == "reserve" {
				reserveAmt = val.Balance
				tmpBid = val.BId
			}
			if val.Endpoint == "debitReserve" {
				// tmpDebitAmt, _ := strconv.ParseFloat(val["amount"], 64)
				totalAmt += val.Amount
			}
			// If cancel reserve found then halt this process and return an error message.
			if val.Endpoint == "cancelReserve" {
				strRes = genErrResponse("0", "CancelReserve found Reserve:"+reserveId)
				c.Header("Content-Type", "text/plain; charset=utf-8")
				c.String(200, strRes)
				return
			}
		}
	}

	totalAmt = totalAmt + nAmount
	// fmt.Println("nAmount ---------------")
	// fmt.Println(totalAmt)
	// fmt.Println(reserveAmt)
	// fmt.Println("----------------------")
	// nAmount is a convert from amount in request query.
	if totalAmt > reserveAmt {
		strRes = genErrResponse("0", "Insufficient amount found!")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// genTrx
	genTrx := betLogInfo.ReserveId + strconv.Itoa(rand.Intn(200))
	// genTrx := betLogInfo.ReserveId + reqId

	genInfo := map[string]string{
		"error_code": "0",
		"genTrx":     genTrx,
		"error_msg":  "No Error",
		"balance":    fmt.Sprintf("%f", reserveAmt),
	}
	strRes = genDebitReserve(genInfo)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, strRes)

	debitReserveLog := statics.BetLog{
		BId:       tmpBid,
		CId:       cId,
		ReserveId: reserveId,
		ReqId:     reqId,
		Balance:   reserveAmt,
		Endpoint:  "debitReserve",
		XmlBody:   string(body),
		ResBody:   strRes,
		BlType:    "img",
		Amount:    nAmount,
	}
	// Insert a cancel reserve log
	blId := models.InsertBetLog(debitReserveLog)
	fmt.Println(blId)

	return
}

func BetCommitReserve(c *gin.Context) {

	fmt.Println("Commit Reserve")
	custId := c.Query("cust_id")
	reserveId := c.Query("reserve_id")
	strRes := ""

	fmt.Printf("Cust ID: %s, Reserve ID: %s\n", custId, reserveId)

	// Check query parameter (Exists & not empty)
	chkMap := map[string]string{
		"custId":    custId,
		"reserveId": reserveId,
	}
	errMsg := chkURLParams(chkMap)
	if len(errMsg) != 0 {
		strRes = genErrResponse("0", errMsg)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	custDecode, err := base64.RawStdEncoding.DecodeString(custId)
	if err != nil {
		strRes = genErrResponse("0", "error_"+string(custDecode))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	tmpAry := strings.Split(string(custDecode), "|")
	cId, err := strconv.Atoi(tmpAry[1])
	fmt.Printf("Client ID: %d, Cust ID: %s, Reserve ID: %s\n", cId, tmpAry[1], reserveId)

	// Check for duplicate
	betLogInfo := statics.BetLog{
		// BId:       bId,
		CId:       cId,
		ReserveId: reserveId,
		Endpoint:  "commitReserve",
		BlType:    "img",
	}
	resBody := models.GetBetLogResBody(betLogInfo)
	if len(resBody) > 0 {
		fmt.Println("Commit Reserve Found")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, resBody)
		return
	}

	// Check for cancel
	cancelBlInfo := statics.BetLog{
		CId:       cId,
		ReserveId: reserveId,
		Endpoint:  "cancelReserve",
		BlType:    "img",
	}
	resBody = models.GetBetLogResBody(cancelBlInfo)
	if len(resBody) > 0 {
		fmt.Println("Cancel Reserve Found")
		strRes = genErrResponse("0", "error_"+string(custDecode))
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// Get reserve amount
	reserveBetLog := statics.BetLog{
		CId:       cId,
		ReserveId: reserveId,
		// Endpoint:  "reserve",
	}
	reserveBody := models.GetBetLog(reserveBetLog)
	if len(reserveBody) < 1 {
		strRes = genErrResponse("0", "No records found! ReserveID:"+reserveId)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}
	var reserveAmt float64
	var tmpBid int
	var totalAmt float64 = 0

	if len(reserveBody) > 0 {
		for idx, val := range reserveBody {
			fmt.Printf("Idx: %d, Amt: %f, bId: %d, endpoint: %s \n", idx, val.Amount, val.BId, val.Endpoint)
			// if val["endpoint"] == "reserve" {
			if val.Endpoint == "reserve" {
				reserveAmt = val.Balance
				tmpBid = val.BId
			}
			if val.Endpoint == "debitReserve" {
				// tmpDebitAmt, _ := strconv.ParseFloat(val["amount"], 64)
				totalAmt += val.Amount
			}
			// If cancel reserve found then halt this process and return an error message.
			if val.Endpoint == "cancelReserve" {
				strRes = genErrResponse("0", "CancelReserve found Reserve:"+reserveId)
				c.Header("Content-Type", "text/plain; charset=utf-8")
				c.String(200, strRes)
				return
			}
		}
	}

	totalAmt = totalAmt + nAmount
	fmt.Println("nAmount ---------------")
	fmt.Println(totalAmt)
	fmt.Println(reserveAmt)
	fmt.Println("----------------------")
	// nAmount is a convert from amount in request query.
	if totalAmt > reserveAmt {
		strRes = genErrResponse("0", "Insufficient amount found!")
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, strRes)
		return
	}

	// genTrx
	genTrx := betLogInfo.ReserveId + strconv.Itoa(rand.Intn(200))
	// genTrx := betLogInfo.ReserveId + reqId

	genInfo := map[string]string{
		"error_code": "0",
		"genTrx":     genTrx,
		"error_msg":  "No Error",
		"balance":    fmt.Sprintf("%f", reserveAmt),
	}
	strRes = genDebitReserve(genInfo)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, strRes)

	debitReserveLog := statics.BetLog{
		BId:       tmpBid,
		CId:       cId,
		ReserveId: reserveId,
		ReqId:     reqId,
		Balance:   reserveAmt,
		Endpoint:  "debitReserve",
		XmlBody:   string(body),
		ResBody:   strRes,
		BlType:    "img",
		Amount:    nAmount,
	}
	// Insert a cancel reserve log
	blId := models.InsertBetLog(debitReserveLog)
	fmt.Println(blId)

	return

}

func BetDebitCustomer(c *gin.Context) {
	fmt.Println("Debit Customer")
	custId := c.Query("cust_id")
	reqId := c.Query("req_id")
	amount := c.Query("amount")
	fmt.Printf("Cust ID: %s, Req ID: %s, Amount: %s\n", custId, reqId, amount)

	body, _ := ioutil.ReadAll(c.Request.Body)
	models.ParseDebitCustomer(string(body))

	// c.JSON(http.StatusOK, gin.H{"data": books})
	c.JSON(http.StatusOK, gin.H{"data": "Debit Customer"})
}

func BetCreditCustomer(c *gin.Context) {
	fmt.Println("Credit Customer")
	custId := c.Query("cust_id")
	reqId := c.Query("req_id")
	amount := c.Query("amount")
	fmt.Printf("Cust ID: %s, Req ID: %s, Amount: %s\n", custId, reqId, amount)

	body, _ := ioutil.ReadAll(c.Request.Body)
	models.ParseCreditCustomer(string(body))

	c.JSON(http.StatusOK, gin.H{"data": "Credit Customer"})
}

// func ErrInsertSeq() {
// 	fmt.Println("ErrInsertSeq")
// 	debitReserveLog := statics.BetLog{
// 		BId:       tmpBid,
// 		CId:       cId,
// 		ReserveId: reserveId,
// 		ReqId:     reqId,
// 		Balance:   reserveAmt,
// 		Endpoint:  "debitReserve",
// 		XmlBody:   string(body),
// 		ResBody:   strRes,
// 		BlType:    "img",
// 		Amount:    nAmount,
// 	}
// 	// Insert a cancel reserve log
// 	blId := models.InsertBetLog(debitReserveLog)
// 	fmt.Println("ErrInsertSeq End")
// }
