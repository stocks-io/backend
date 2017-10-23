package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type leader struct {
	Email string
	Cash  float64
}

type loginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type logoutRequest struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type registerRequest struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

func setupUserRoutes() {
	users := app.Group("/users")
	{
		users.POST("/login", func(c *gin.Context) {
			var req loginRequest
			c.ShouldBindWith(&req, binding.Form)
			if !userExists(req.Email) {
				c.JSON(http.StatusNotFound, gin.H{"message": "user does not exist"}) // 404
				return
			}
			var hash string
			err := db.QueryRow("SELECT password FROM userinfo WHERE email=?", req.Email).Scan(&hash)
			checkErr(err)
			if !checkPasswordHash(req.Password, hash) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "incorrect password"}) // 401
				return
			}
			id, err := getUserId(req.Email)
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
			c.JSON(http.StatusAccepted, gin.H{ // 201
				"userId":  id,
				"token":   fmt.Sprintf("%s", token),
				"expires": then,
			})
		})
		users.POST("/logout", func(c *gin.Context) {
			var req logoutRequest
			c.ShouldBindWith(&req, binding.Form)
			var exists bool
			err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM sessions WHERE token=?", req.Token).Scan(&exists)
			checkErr(err)
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{ // 404
					"message": "token does not exist",
				})
				return
			}
			_, err = db.Exec("DELETE FROM sessions WHERE token=?", req.Token)
			checkErr(err)
			c.JSON(http.StatusOK, gin.H{ // 200
				"message": "successfully logged out",
			})
		})
		users.POST("/register", func(c *gin.Context) {
			var req registerRequest
			c.ShouldBindWith(&req, binding.Form)
			if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
				c.JSON(http.StatusBadRequest, gin.H{ // 400
					"message": "all fields are required",
				})
				return
			}
			if userExists(req.Email) {
				c.JSON(http.StatusNotAcceptable, gin.H{ // 406
					"message": "email already taken",
				})
				return
			}
			if emailExists(req.Email) {
				c.JSON(http.StatusNotAcceptable, gin.H{ // 406
					"message": "email already taken",
				})
				return
			}
			now := time.Now().Unix()
			stmt, err := db.Prepare("INSERT userinfo SET first_name=?, last_name=?, email=?, password=?, added=?")
			checkErr(err)
			hashed, err := hashPassword(req.Password)
			checkErr(err)
			_, err = stmt.Exec(req.FirstName, req.LastName, req.Email, hashed, now)
			checkErr(err)
			id, err := getUserId(req.Email)
			checkErr(err)
			stmt, err = db.Prepare("INSERT portfolio SET user_id=?, cash=?")
			checkErr(err)
			_, err = stmt.Exec(id, 10000)
			checkErr(err)
			c.JSON(http.StatusCreated, gin.H{ // 201
				"message": "successfully registered",
			})
		})
		users.GET("/leaderboard", func(c *gin.Context) {
			rows, err := db.Query("SELECT email, cash FROM userinfo JOIN portfolio ON portfolio.user_id=userinfo.id ORDER BY cash DESC")
			checkErr(err)
			var leaderboard []leader
			for rows.Next() {
				var email string
				var cash float64
				err = rows.Scan(&email, &cash)
				checkErr(err)
				user := leader{
					Email: email,
					Cash:  cash,
				}
				leaderboard = append(leaderboard, user)
			}
			c.JSON(http.StatusOK, leaderboard) // 200
		})
	}
}
