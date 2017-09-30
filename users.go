package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupUserRoutes() {
	users := app.Group("/users")
	{
		users.POST("/login", func(c *gin.Context) {
			var id int
			var req loginRequest
			c.BindWith(&req, binding.Form)
			if !userExists(req.Username) {
				c.JSON(401, gin.H{"message": "user does not exist"})
				return
			}
			rows, err := db.Query("SELECT id FROM userinfo WHERE username=?", req.Username)
			checkErr(err)
			rows.Next()
			err = rows.Scan(&id)
			checkErr(err)
			token, err := exec.Command("uuidgen").Output()
			token = token[0 : len(token)-1]
			checkErr(err)
			now := time.Now().Unix()
			then := time.Now().Add(time.Hour * 24).Unix()
			stmt, err := db.Prepare("INSERT sessions SET user_id=?, token=?, added=?, expires=?")
			checkErr(err)
			_, err = stmt.Exec(id, token, now, then)
			checkErr(err)
			fmt.Printf("%s\n", token)
			c.JSON(200, gin.H{
				"userId":  id,
				"token":   fmt.Sprintf("%s", token),
				"expires": then,
			})
		})
		users.POST("/logout", func(c *gin.Context) {
			var req logoutRequest
			c.BindWith(&req, binding.Form)
			var exists bool
			err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM sessions WHERE token=?", req.Token).Scan(&exists)
			checkErr(err)
			if !exists {
				c.JSON(401, gin.H{
					"message": "token does not exist",
				})
				return
			}
			_, err = db.Exec("DELETE FROM sessions WHERE token=?", req.Token)
			checkErr(err)
			c.JSON(200, gin.H{
				"message": "successfully logged out",
			})
		})
		users.GET("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": "register",
			})
		})
	}
}
