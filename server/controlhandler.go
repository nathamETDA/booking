package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func openHandler(c *gin.Context) {
	open = true
	c.JSON(http.StatusOK, gin.H{
		"message": "Open at " + time.Now().String(),
	})
}

func closeHandler(c *gin.Context) {
	open = false
	c.JSON(http.StatusOK, gin.H{
		"message": "Close at " + time.Now().String(),
	})
}

func statusHandler(c *gin.Context) {
	c.String(http.StatusOK,
		"Open:%v\nClient at door: %v\nClient in waiting room: %v\nexitDoor[]: %v",
		open, len(exitDoor), len(waitQueue), exitDoor)
}
