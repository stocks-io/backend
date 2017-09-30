package main

type buyRequest struct {
	UserId int    `form:"userId" json:"userId" binding:"required"`
	Units  int    `form:"units" json:"units" binding:"required"`
	Symbol string `form:"symbol" json:"symbol" binding:"required"`
}

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type logoutRequest struct {
	Token string `form:"token" json:"token" binding:"required"`
}
