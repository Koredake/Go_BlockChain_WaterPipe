package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/db"
)

func InsertUserInfo(c *gin.Context) {
	address := c.PostForm("address")

	err, strErr := db.ShowUser(address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": strErr,
		})
		//respError(c, err)
		return
	}
	SendEth(address)
	respOK(c, nil)
}
