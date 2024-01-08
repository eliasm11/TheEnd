package api

import (
	"db"
	"net/http"
	"strings"
	"structs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type user struct {
	userGroup *gin.RouterGroup
}

func (user *user) setCheckoutApi() {

	user.userGroup.GET("/checkout", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "payment.html", gin.H{})
	})
	user.userGroup.POST("/buy", func(ctx *gin.Context) {


		var Item map[db.Container]map[db.Kind]map[db.Id]map[db.Size]map[db.Color]db.Qty 
		orders := db.Orders{}

		if err := ctx.ShouldBindJSON(&Item); err != nil {
			ctx.JSON(http.StatusOK, err.Error())//"check json")
			return
		}

		orders.Item = Item
		orders.User.Username = sessions.Default(ctx).Get("username").(string)

		if err := db.MainDB.Buy(&orders); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
		}
		ctx.JSON(http.StatusOK, "OK")

	})
}

func (user *user) setInformationApi() {

	user.userGroup.POST("/phone", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		phone := ctx.PostForm("phone")
		if err := db.MainDB.Users.UpdataPhone(username, phone); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "update")
	})

	user.userGroup.GET("/phone", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		user := structs.User{}
		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		if user.Phone == "" {
			ctx.JSON(http.StatusOK, "")
		} else {
			ctx.JSON(http.StatusOK, user.Phone)
		}

	})

	user.userGroup.POST("/visa", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)

		visa := structs.Visa{}
		if err := ctx.ShouldBindJSON(&visa); err != nil {
			ctx.JSON(http.StatusOK, "check the json")
			return
		}
		if err := db.MainDB.Users.AddVisa(username, visa); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "add visa")
	})

	user.userGroup.GET("/visa", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		if len(user.UserVisa.Number) < 16 {
			ctx.JSON(http.StatusOK, "")
			return
		}
		visa := user.UserVisa.Number[:3] + strings.Repeat("*", 10) + user.UserVisa.Number[13:]
		ctx.JSON(http.StatusOK, visa)
	})

	user.userGroup.POST("/addr", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		addr := structs.Addr{}
		addr.City = strings.TrimSpace(addr.City)
		addr.Area = strings.TrimSpace(addr.Area)
		addr.Street = strings.TrimSpace(addr.Street)

		if err := ctx.ShouldBindJSON(&addr); err != nil {
			ctx.JSON(http.StatusOK, "check the json plz")
			return
		}
		if err := db.MainDB.Users.UpdataAddr(username, addr); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "update")
	})

	user.userGroup.GET("/addr", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, user.UserAddr)

	})
	user.userGroup.GET("/username", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		user := structs.User{}

		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, user.Name)
	})

	user.userGroup.POST("/name", func(ctx *gin.Context) {

		username := sessions.Default(ctx).Get("username").(string)

		Name := ctx.PostForm("Name")

		if err := db.MainDB.Users.UpdateName(username, Name); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "update")

	})

	user.userGroup.POST("/password", func(ctx *gin.Context) {

		session := sessions.Default(ctx)
		type _password struct {
			NewPassowrd string `json:"newPassowrd"`
			OldPassowrd string `json:"oldPassowrd"`
		}
		password := _password{}
		if err := ctx.ShouldBindJSON(&password); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}

		if err := db.MainDB.Users.UpdataPassword(session.Get("username").(string), password.OldPassowrd, password.NewPassowrd); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		s := sessions.Default(ctx)
		s.Delete("user")
		s.Save()

		ctx.JSON(http.StatusOK, "update")

	})

}

func (user *user) setLogoutApi() {
	user.userGroup.GET("/logout", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Clear()
		session.Save()

	})
}

func (user *user) setCommintApi() {

	user.userGroup.DELETE("/commint", func(ctx *gin.Context) {

		type UserCommint struct {
			Commint   string `json:"commint"`
			Container string `json:"container"`
			Kind      string `json:"kind"`
			Idmodel   int    `json:"idmodel"`
		}
		commint := UserCommint{}

		if err := ctx.ShouldBindJSON(&commint); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}

		username := sessions.Default(ctx).Get("username").(string)

		if err := db.MainDB.Stock.DeleteCommint(commint.Idmodel, commint.Container, commint.Kind, username, commint.Commint); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "delete")

	})
	user.userGroup.POST("/commint", func(ctx *gin.Context) {

		type UserCommint struct {
			Username  string `josn:"username"`
			Commint   string `json:"commint"`
			Stars     int    `json:"stars"`
			Container string `json:"container"`
			Kind      string `json:"kind"`
			Idmodel   int    `json:"idmodel"`
		}
		commint := UserCommint{}
		if err := ctx.ShouldBindJSON(&commint); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}
		if commint.Commint == "" {
			ctx.JSON(http.StatusOK, "commint empty")
			return
		}
		if commint.Stars <= 0 || commint.Stars > 6 {
			ctx.JSON(http.StatusOK, "stars invaild")
			return
		}

		commint.Username = sessions.Default(ctx).Get("username").(string)

		if err := db.MainDB.AddCommint(commint.Idmodel, commint.Container, commint.Kind, structs.UserCommint{Username: commint.Username, Commint: commint.Commint, Stars: commint.Stars}); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "add")

	})

}
