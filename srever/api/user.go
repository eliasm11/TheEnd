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
		username := sessions.Default(ctx).Get("username").(string)
		orders := []db.Orders{}

		type Err struct {
			Id  int    `json:"id"`
			ERR string `json:"err"`
		}
		err := []Err{}
		if err := ctx.ShouldBindJSON(&orders); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}
		for _, v := range orders {
			v.Username = username
			if err_ := db.MainDB.Buy(v); err_ != nil {
				if err_ != nil && err_ != db.ErrDataBase {
					ctx.JSON(http.StatusOK, err_.Error())
					return
				}
				err = append(err, Err{Id: v.IdModel, ERR: err_.Error()})
			}
		}
		if len(err) == 0 {
			ctx.JSON(http.StatusOK, "OK")
		} else {
			ctx.JSON(http.StatusOK, &err)
		}

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
		visa := user.UserVisa.Number[:3] + strings.Repeat("*", 10) + user.UserVisa.Number[13:]
		ctx.JSON(http.StatusOK, visa)
	})

	user.userGroup.POST("/addr", func(ctx *gin.Context) {
		username := sessions.Default(ctx).Get("username").(string)
		addr := structs.Addr{}

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

		ctx.JSON(http.StatusOK, "update")

	})

}

func (user *user) setLogoutApi() {
	user.userGroup.GET("/logout", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Set("user", nil)
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
