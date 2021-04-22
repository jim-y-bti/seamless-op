package models

import (
	"database/sql"
	"fmt"
	dbTables "seamlessop/models/dbTable"
	"seamlessop/statics"
	"strconv"
)

func InsertBets(where statics.Bets) int {
	fmt.Println("InsertBets")
	bInfo := dbTables.Bets()

	sqlStatement := "insert into " + bInfo["tableName"]
	sqlStatement += " (" + bInfo["cw_id"] + ","
	sqlStatement += bInfo["b_reserve_id"] + ","
	sqlStatement += bInfo["b_balance"] + ","
	sqlStatement += bInfo["b_status"] + ","
	sqlStatement += bInfo["b_created_at"] + ","
	sqlStatement += bInfo["b_updated_at"] + ") values "
	sqlStatement += "($1, $2, $3, $4, now(), now() ) returning " + string(bInfo["b_id"])

	fmt.Println(sqlStatement)
	id := 0
	err := DB.QueryRow(sqlStatement, where.CwId, where.ReserveId, where.Balance, where.Status).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func InsertBetLog(where statics.BetLog) string {
	fmt.Println("InsertBetLog")
	blInfo := dbTables.BetLog()

	sqlStatement := "insert into " + blInfo["tableName"]
	sqlStatement += " (" + blInfo["b_id"] + ","
	sqlStatement += blInfo["c_id"] + ","
	sqlStatement += blInfo["bl_reserve_id"] + ","
	sqlStatement += blInfo["bl_req_id"] + ","
	sqlStatement += blInfo["bl_balance"] + ","
	sqlStatement += blInfo["bl_amount"] + ","
	sqlStatement += blInfo["bl_endpoint"] + ","
	sqlStatement += blInfo["bl_xml_body"] + ","
	sqlStatement += blInfo["bl_res_body"] + ","
	sqlStatement += blInfo["bl_type"] + ","
	sqlStatement += blInfo["bl_created_at"] + ") values "
	sqlStatement += "($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now()) returning " + blInfo["bl_id"]
	// DB.QueryRow()
	id := 0
	fmt.Println(sqlStatement)
	err := DB.QueryRow(sqlStatement,
		where.BId,
		where.CId,
		where.ReserveId,
		where.ReqId,
		where.Balance,
		where.Amount,
		where.Endpoint,
		where.XmlBody,
		where.ResBody,
		where.BlType,
	).Scan(&id)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(id)
}

func GetBetLogResBody(where statics.BetLog) string {

	fmt.Println("GetBetLogResBody")
	// fmt.Println(where)
	var resBody string
	var inAry = []interface{}{where.CId, where.ReserveId}

	blInfo := dbTables.BetLog()
	sqlStatement := "select "
	sqlStatement += blInfo["bl_res_body"]
	sqlStatement += " from " + blInfo["tableName"]
	sqlStatement += " where " + blInfo["c_id"] + "= $1"
	sqlStatement += " and " + blInfo["bl_reserve_id"] + "= $2"
	if len(where.Endpoint) > 0 {
		sqlStatement += " and " + blInfo["bl_endpoint"] + "= $3"
		inAry = append(inAry, where.Endpoint)
	}
	if len(where.ReqId) > 0 {
		sqlStatement += " and " + blInfo["bl_req_id"] + "= $4"
		inAry = append(inAry, where.Endpoint)
	}

	fmt.Println(sqlStatement)
	row := DB.QueryRow(sqlStatement, inAry...)

	switch err := row.Scan(&resBody); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(resBody)
	default:
		panic(err)
	}
	return resBody
}

func GetBetLog(where statics.BetLog) []statics.BetLog {
	fmt.Println("GetBetLog")

	var inAry = []interface{}{where.CId, where.ReserveId}

	blInfo := dbTables.BetLog()
	sqlStatement := "select "
	sqlStatement += blInfo["bl_id"] + ","
	sqlStatement += blInfo["b_id"] + ","
	sqlStatement += blInfo["c_id"] + ","
	sqlStatement += blInfo["bl_reserve_id"] + ","
	sqlStatement += blInfo["bl_req_id"] + ","
	sqlStatement += blInfo["bl_balance"] + ","
	sqlStatement += blInfo["bl_amount"] + ","
	sqlStatement += blInfo["bl_res_body"] + ","
	sqlStatement += blInfo["bl_endpoint"]
	sqlStatement += " from " + blInfo["tableName"]
	sqlStatement += " where " + blInfo["c_id"] + "= $1"
	sqlStatement += " and " + blInfo["bl_reserve_id"] + "= $2"
	if len(where.Endpoint) > 0 {
		sqlStatement += " and " + blInfo["bl_endpoint"] + "= $3"
		inAry = append(inAry, where.Endpoint)
	}
	if len(where.ReqId) > 0 {
		sqlStatement += " and " + blInfo["bl_req_id"] + "= $3"
		inAry = append(inAry, where.ReqId)
	}

	fmt.Println(sqlStatement)
	rows, err := DB.Query(sqlStatement, inAry...)
	if err != nil {
		panic(err)
	}

	// res := map[int]map[string]string{}
	var res []statics.BetLog
	for rows.Next() {
		var blId int
		var bId int
		var cId int
		var reserveId string
		var reqId string
		var balance float64
		var amount float64
		var resBody string
		var endpoint string

		err = rows.Scan(&blId, &bId, &cId, &reserveId, &reqId, &balance, &amount, &resBody, &endpoint)
		if err != nil {
			panic(err)
		}
		fmt.Println("----- GetBetLog -----")
		fmt.Println(blId, cId, reserveId, reqId, amount)
		res = append(res, statics.BetLog{
			CId:       cId,
			BId:       bId,
			ReserveId: reserveId,
			ReqId:     reqId,
			Balance:   balance,
			Amount:    amount,
			ResBody:   resBody,
			Endpoint:  endpoint,
		})

	}

	// fmt.Println(res)
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return res
}

func UpdateBet(where statics.Bets) string {
	fmt.Println("UpdateBet")
	bInfo := dbTables.Bets()

	sqlStatement := "update " + bInfo["tableName"]
	sqlStatement += " set " + bInfo["b_balance"] + " = $1, "
	sqlStatement += bInfo["b_updated_at"] + " = now()"
	sqlStatement += " where " + bInfo["b_reserve_id"] + " = $2"
	fmt.Println(sqlStatement)

	res, err := DB.Exec(sqlStatement, where.Balance, where.ReserveId)
	if err != nil {
		panic(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Row effected: %d \n", count)

	return "updated"
}
