package srever

import (
	"srever/api"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Srever int

func (*Srever) Run() {

	//gin.SetMode(gin.ReleaseMode)
	store := cookie.NewStore([]byte("secret"))
	serverEngin := gin.New()

	serverEngin.Use(gin.Logger(), gin.Recovery(), sessions.Sessions("login", store))

	serverEngin.LoadHTMLGlob("templates/*.html")
	api.InitApi().Setup(serverEngin)

	serverEngin.Run(":8080")

}

func InitSever() *Srever {
	return new(Srever)
}
