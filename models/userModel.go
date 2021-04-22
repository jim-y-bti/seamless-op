package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	dbTables "seamlessop/models/dbTable"
	"seamlessop/statics"
	"strconv"
)

func GetUserToken(where []string) map[string]string {

	fmt.Println("GetUserToken")
	var c_id string
	var c_name string
	var c_token string
	var cw_balance string
	var ck_city string
	var ck_country string
	var ck_currency_code string

	cInfo := dbTables.Clients()
	cwInfo := dbTables.ClientsWallets()
	ckInfo := dbTables.ClientsKYC()

	// sql := "select "
	sqlStatement := "select "
	sqlStatement += "c." + cInfo["c_id"] + ","
	sqlStatement += cInfo["c_name"] + ","
	sqlStatement += cInfo["c_token"] + ","
	sqlStatement += cwInfo["cw_balance"] + ","
	sqlStatement += ckInfo["ck_city"] + ","
	sqlStatement += ckInfo["ck_country"] + ","
	sqlStatement += ckInfo["ck_currency_code"]
	sqlStatement += " from " + cInfo["tableName"] + " c "
	sqlStatement += " join " + cwInfo["tableName"] + " cw on c." + cInfo["c_id"] + " = cw." + cwInfo["c_id"]
	sqlStatement += " join " + ckInfo["tableName"] + " ck on c." + cInfo["c_id"] + " = ck." + cwInfo["c_id"]
	sqlStatement += " where " + cInfo["c_token"] + "='" + where[0] + "'"
	fmt.Println(sqlStatement)

	row := DB.QueryRow(sqlStatement)
	mapRes := make(map[string]string)

	switch err := row.Scan(&c_id, &c_name, &c_token, &cw_balance, &ck_city, &ck_country, &ck_currency_code); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		mapRes["c_id"] = c_id
		mapRes["c_name"] = c_name
		mapRes["c_token"] = c_token
		mapRes["cw_balance"] = cw_balance
		mapRes["ck_city"] = ck_city
		mapRes["ck_country"] = ck_country
		mapRes["ck_currency_code"] = ck_currency_code
	default:
		panic(err)
	}
	return mapRes
}

func GetUserData(where []string) map[string]string {

	fmt.Println("GetUserData")
	var cId string
	var cName string
	var cToken string

	tbInfo := dbTables.Clients()
	// sql := "select "
	sqlStatement := "select "
	sqlStatement += tbInfo["c_id"] + "," + tbInfo["c_name"] + "," + tbInfo["c_token"]
	sqlStatement += " from " + tbInfo["tableName"] + " where " + tbInfo["c_name"] + "='" + where[0] + "'"
	fmt.Println(sqlStatement)

	row := DB.QueryRow(sqlStatement)
	mapRes := make(map[string]string)

	switch err := row.Scan(&cId, &cName, &cToken); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		// mapRes := map[string]string{"cId": cId, "cName": cName, "cToken": cToken}
		mapRes["cId"] = cId
		mapRes["cName"] = cName
		mapRes["cToken"] = cToken
	default:
		panic(err)
	}

	return mapRes
}

func LockAndUpdateBalance(where statics.BetLog) map[string]string {

	fmt.Println("LockAndUpdateBalance")
	var cw_id string
	var c_id string
	var cw_balance string

	mapRes := make(map[string]string)
	cwInfo := dbTables.ClientsWallets()

	sqlStatement := "select "
	sqlStatement += cwInfo["cw_id"] + ","
	sqlStatement += cwInfo["c_id"] + ","
	sqlStatement += cwInfo["cw_balance"]
	sqlStatement += " from " + cwInfo["tableName"]
	sqlStatement += " where 1=1 "
	if where.CId != 0 {
		sqlStatement += " and " + cwInfo["c_id"] + " = $1"
	}
	sqlStatement += " for update"
	fmt.Println(sqlStatement, where.CId)

	// Start the transaction
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	row := tx.QueryRow(sqlStatement)
	err = row.Scan(&cw_id, &c_id, &cw_balance)
	if err != nil {
		tx.Rollback()
		return mapRes
	}

	currBalance, err := strconv.ParseFloat(cw_balance, 64)
	if err != nil {
		tx.Rollback()
		return mapRes
	}
	// reserveAmount, err := strconv.ParseFloat(where["amount"], 64)
	// if err != nil {
	// 	tx.Rollback()
	// 	return mapRes
	// }
	newBalance := currBalance - where.Amount
	if newBalance < 0 {
		mapRes["newBalance"] = fmt.Sprintf("%f", newBalance)
		return mapRes
	}

	fmt.Printf("cw_id: %s, Reserve Amount: %f, Curr Amount: %s \n", cw_id, where.Amount, cw_balance)
	// Update Bets with new reserve
	sqlStatement = "update " + cwInfo["tableName"]
	sqlStatement += " set " + cwInfo["cw_balance"] + "= $1"
	fmt.Println(sqlStatement)
	_, err = tx.ExecContext(ctx, sqlStatement, newBalance)
	if err != nil {
		tx.Rollback()
		return mapRes
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	mapRes["cw_id"] = cw_id
	mapRes["c_id"] = c_id
	mapRes["newBalance"] = fmt.Sprintf("%f", newBalance)

	return mapRes
}
