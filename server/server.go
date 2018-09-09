package main

import (
	"fmt"
	"net/http"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func mockHandler(c *gin.Context) {
	time.Sleep(time.Millisecond * 100 * speed)
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Home",
	})
}

func enterWaitRoomHandler(c *gin.Context) {
	if len(waitQueue) > maxWaitQueue {
		c.JSON(200, gin.H{"message": "Queue to long"})
		return
	}
	client := newClientState()
	fmt.Println(client)
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	token.Claims = jwt_lib.MapClaims{
		"Id":  client.id,
		"exp": time.Now().Add(time.Hour * 10).Unix(),
	}
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(waitSecret))
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not generate token"})
		return
	}
	// waitQueue.
	waitQueue = append(waitQueue, client)
	c.JSON(200, gin.H{"id": client.id, "token": tokenString, "message": "OK"})
}

func exitwaitroomHandler(c *gin.Context) {
	id := c.Param("id")

	if moveOut(id) {
		c.JSON(200, gin.H{"message": "OK"})
		return
	}
	c.JSON(200, gin.H{"message": "Please wait"})
}

func openClose() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !open {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "CLOSE",
			})
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.GET("/", homeHandler)

	w := r.Group("/waitroom")
	w.Use(openClose())
	w.GET("/enter", enterWaitRoomHandler)
	w.GET("/exit/:id", exitwaitroomHandler)

	// l := r.Group("/login")
	// l.Use(jwt.Auth(waitSecret))
	// l.GET("/", mockHandler)
	r.GET("/login", mockHandler)

	r.GET("/booking", mockHandler)
	r.GET("/booking/reserveseat", mockHandler)
	r.GET("/booking/gaseat", mockHandler)
	r.GET("/payment", mockHandler)

	// operation command
	r.GET("/open", openHandler)
	r.GET("/close", closeHandler)
	r.GET("/status", statusHandler)

	go wait2door()
	go purgeWaitQueue()
	go purgeExitDoor()

	r.Run() // listen and serve on 0.0.0.0:8080

}
