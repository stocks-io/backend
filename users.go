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
			var req loginRequest
			c.BindWith(&req, binding.Form)
			if !userExists(req.Username) {
				c.JSON(401, gin.H{"message": "user does not exist"})
				return
			}
			id, err := getUserId(req.Username)
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
		users.POST("/register", func(c *gin.Context) {
			var req registerRequest
			c.BindWith(&req, binding.Form)
			if userExists(req.Username) {
				c.JSON(401, gin.H{
					"message": "username already taken",
				})
				return
			}
			if emailExists(req.Email) {
				c.JSON(401, gin.H{
					"message": "email already taken",
				})
				return
			}
			now := time.Now().Unix()
			stmt, err := db.Prepare("INSERT userinfo SET first_name=?, last_name=?, username=?, email=?, password=?, added=?")
			checkErr(err)
			_, err = stmt.Exec(req.FirstName, req.LastName, req.Username, req.Email, req.Password, now)
			checkErr(err)
			id, err := getUserId(req.Username)
			checkErr(err)
			stmt, err = db.Prepare("INSERT portfolio SET user_id=?, cash=?")
			checkErr(err)
			_, err = stmt.Exec(id, 10000)
			checkErr(err)
			c.JSON(200, gin.H{
				"page": "successfully registered",
			})
		})
	}
}
