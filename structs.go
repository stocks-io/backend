package main

type orderRequest struct {
	Token  string `form:"token" json:"token" binding:"required"`
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

type registerRequest struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

type leader struct {
	Username string
	Cash     float64
}
