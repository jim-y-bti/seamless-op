package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"seamlessop/models"

	"github.com/gin-gonic/gin"
)

func UserValidToken(c *gin.Context) {

	fmt.Println("UserValidToken")
	authToken := c.Query("auth_token")
	fmt.Println(string(authToken))

	where := []string{
		authToken,
	}

	userData := models.GetUserToken(where)
	jsonRes, _ := json.Marshal(userData)
	fmt.Println(string(jsonRes))
	strRes := ""
	if len(userData) == 0 {
		strRes += "error_code=-3\nerror_message=InvalidToken"
	} else {
		sEnc := "rm-rfbois|" + userData["c_id"]
		base64CustId := base64.RawStdEncoding.EncodeToString([]byte(sEnc))

		strRes += "error_code=0\nerror_message=NO_ERRORS\n"
		strRes += "balance=" + userData["cw_balance"] + "\n"
		strRes += "data=\n"
		strRes += "city=" + userData["ck_city"] + "\n"
		strRes += "country=" + userData["ck_country"] + "\n"
		strRes += "cust_id=" + base64CustId + "\n"
		strRes += "extSessionID=\n"
		strRes += "currency_code=" + userData["ck_currency_code"]
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, strRes)
	return
}

func UserCreateToken(c *gin.Context) {

	fmt.Println("UserCreateToken")
	custId := c.Query("cust_id")
	where := []string{
		custId,
	}

	userData := models.GetUserData(where)
	jsonRes, _ := json.Marshal(userData)
	fmt.Println(string(jsonRes))
	c.JSON(http.StatusOK, gin.H{"data": jsonRes})
}
