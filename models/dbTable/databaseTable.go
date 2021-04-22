package dbTables

func Clients() map[string]string {

	res := make(map[string]string)
	res["tableName"] = "clients"
	res["c_id"] = "c_id"
	res["c_name"] = "c_name"
	res["c_token"] = "c_token"
	res["created_at"] = "created_at"

	return res
}

func ClientsWallets() map[string]string {
	res := make(map[string]string)
	res["tableName"] = "client_wallets"
	res["cw_id"] = "cw_id"
	res["c_id"] = "c_id"
	res["cw_balance"] = "cw_balance"
	res["cw_created_at"] = "cw_created_at"

	return res
}

func ClientsKYC() map[string]string {
	res := make(map[string]string)
	res["tableName"] = "clients_kyc"
	res["ck_id"] = "ck_id"
	res["c_id"] = "c_id"
	res["ck_city"] = "ck_city"
	res["ck_country"] = "ck_country"
	res["ck_currency_code"] = "ck_currency_code"

	return res
}

func Bets() map[string]string {
	res := make(map[string]string)
	res["tableName"] = "bets"
	res["b_id"] = "b_id"
	res["cw_id"] = "cw_id"
	res["b_reserve_id"] = "b_reserve_id"
	res["b_balance"] = "b_balance"
	res["b_status"] = "b_status"
	res["b_created_at"] = "b_created_at"
	res["b_updated_at"] = "b_updated_at"

	return res
}

func BetLog() map[string]string {
	res := make(map[string]string)
	res["tableName"] = "bet_log"
	res["bl_id"] = "bl_id"
	res["b_id"] = "b_id"
	res["c_id"] = "c_id"
	res["bl_reserve_id"] = "bl_reserve_id"
	res["bl_req_id"] = "bl_req_id"
	res["bl_balance"] = "bl_balance"
	res["bl_amount"] = "bl_amount"
	res["bl_endpoint"] = "bl_endpoint"
	res["bl_xml_body"] = "bl_xml_body"
	res["bl_res_body"] = "bl_res_body"
	res["bl_created_at"] = "bl_created_at"
	res["bl_type"] = "bl_type"

	return res
}
