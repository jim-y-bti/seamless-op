package controllers

import "fmt"

func genErrResponse(errCode string, errMsg string) string {
	fmt.Println("genErrResponse")
	response := "error_code=" + errCode + "\nerror_message=" + errMsg
	return response
}

func genReserveResponse(obj map[string]string) string {
	fmt.Println("genReserveResponse")
	response := "error_code=" + obj["error_code"] + "\n"
	response += "error_message=" + obj["error_msg"] + "\n"
	response += "trx_id=" + obj["genTrx"] + "\n"
	response += "balance=" + obj["balance"] + "\n"
	response += "BonusUsed=0\n"
	return response
}

func genCancelReserveResponse(obj map[string]string) string {
	fmt.Println("genCancelReserveResponse")
	response := "error_code=" + obj["error_code"] + "\n"
	response += "error_message=" + obj["error_msg"] + "\n"
	response += "balance=" + obj["balance"] + "\n"
	return response
}

func genDebitReserve(obj map[string]string) string {
	fmt.Println("genDebitReserve")
	response := "error_code=" + obj["error_code"] + "\n"
	response += "trx_id=" + obj["genTrx"] + "\n"
	response += "error_message=" + obj["error_msg"] + "\n"
	response += "balance=" + obj["balance"] + "\n"
	return response
}

// Check for exist / not empty.
func chkURLParams(obj map[string]string) string {
	fmt.Println("chkURLParams")
	response := ""
	for idx, val := range obj {
		if val == "{na}" {
			response += idx + " is missing. "
		}
		if len(val) == 0 {
			response += idx + " is empty. "
		}
	}
	return response
}
