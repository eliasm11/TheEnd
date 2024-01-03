package api

import (
	"db"
	"net/http"
	"strconv"
	"strings"
	"structs"

	"net/smtp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Api struct {
	admin
	user
}

func InitApi() *Api {
	return &Api{}
}

func (api *Api) Setup(server *gin.Engine) {
	api.setUserApi(server)
	api.setGuestApi(server)
	api.setAdminApi(server)
}

func (api *Api) setUserApi(server *gin.Engine) {
	api.user.userGroup = server.Group("/user",
		func(ctx *gin.Context) {
			session := sessions.Default(ctx)

			if session.Get("user") == nil {
				ctx.HTML(http.StatusOK, "loginAndRegister.html", nil)
				ctx.Abort()
				return
			}

			infoUser := strings.Split(session.Get("user").(string), ",")
			if len(infoUser) != 2 {
				ctx.HTML(http.StatusOK, "loginAndRegister.html", nil)
				ctx.Abort()
				return
			}

			user := structs.User{}
			if err := db.MainDB.Users.GetUser(infoUser[0], &user); err != nil || user.Password != infoUser[1] {
				ctx.HTML(http.StatusOK, "loginAndRegister.html", nil)
				ctx.Abort()
				return
			}

		})

	api.user.setLogoutApi()

	api.user.setCommintApi()
	api.user.setCheckoutApi()
	api.user.setInformationApi()
}

func (api *Api) setAdminApi(server *gin.Engine) {
	api.admin.adminGroup = server.Group("/admin", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		if session.Get("admin") == nil {

			ctx.HTML(http.StatusOK, "loginadmin.html", nil)
			ctx.Abort()
			return
		}

		infoAdmin := strings.Split(session.Get("admin").(string), ",")
		if infoAdmin[1] != api.admin.adminUsers[infoAdmin[0]] {

			ctx.HTML(http.StatusOK, "loginadmin.html", nil)
			ctx.Abort()
			return
		}

	})
	api.admin.setOrderApi()
	api.admin.setLogoutApi()
	api.admin.setProductApi()
	api.admin.GetUsers()

}

func (api *Api) setGuestApi(server *gin.Engine) {

	server.GET("/product", func(ctx *gin.Context) {

		kind := ctx.Query("kind")
		Containter := ctx.Query("container")

		id, err := strconv.Atoi(ctx.Query("id"))
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		data, err := db.MainDB.Stock.GetModelsInKind(id, 10, Containter, kind)
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, data)

	})

	server.GET("/NumberProduct", func(ctx *gin.Context) {
		type InputHTTP struct {
			Container string `json:"container"`
			Kind      string `json:"kind"`
		}

		inputHTTP := InputHTTP{}
		if err := ctx.ShouldBindJSON(&inputHTTP); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}

		data, err := db.MainDB.Stock.GetNumberModelsInKind(inputHTTP.Container, inputHTTP.Kind)
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, data)

	})

	server.GET("/AllContainerAndKind", func(ctx *gin.Context) {

		fg := db.MainDB.Stock.GetAllContainerAndKind()
		if len(fg) == 0 {
			ctx.JSON(http.StatusOK, "no data")
		} else {
			ctx.JSON(http.StatusOK, &fg)
		}

	})

	server.NoRoute(func(ctx *gin.Context) {
		Path := ctx.Request.URL.Path[1:]
		if db.MainDB.Stock.IsExist(Path) {
			containerAndkind := db.MainDB.Stock.GetAllContainerAndKind()
			ctx.HTML(http.StatusOK, "selection.html", gin.H{
				"container": ctx.Request.URL.Path[1:],
				"kinds":     containerAndkind[Path],
			},
			)
			ctx.Abort()
			return
		}
		pathArray := strings.Split(Path, "/")
		if len(pathArray) < 2 {
			ctx.Next()
			return
		}
		Number, _ := db.MainDB.Stock.GetNumberModelsInKind(pathArray[0], pathArray[1])
		if listModels, err := db.MainDB.Stock.GetModelsInKind(Number, 12, pathArray[0], pathArray[1]); err == nil {
			ctx.HTML(http.StatusOK, "products.html", gin.H{

				"container":  pathArray[0],
				"kind":       pathArray[1],
				"listModels": listModels,
			})

			ctx.Abort()
			return
		}

	}, gin.WrapH(http.FileServer(http.Dir("public"))))
	server.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		if session.Get("user") == nil || session.Get("username") == nil {
			ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": true})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{"guest": false})
	})

	server.POST("/login", func(ctx *gin.Context) {

		type loginUser struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		userlogin := loginUser{}
		if err := ctx.ShouldBindJSON(&userlogin); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}

		user := structs.User{}
		if err := db.MainDB.Users.GetUser(userlogin.Username, &user); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		if user.Password != userlogin.Password {
			ctx.JSON(http.StatusOK, "check your username or password")
			return

		}
		session := sessions.Default(ctx)
		session.Set("user", userlogin.Username+","+user.Password)
		session.Set("username", userlogin.Username)
		session.Save()
		ctx.JSON(http.StatusOK, "ok")
	})

	server.POST("/register", func(ctx *gin.Context) {

		session := sessions.Default(ctx)
		session.Clear()
		type User1 struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		var useradd User1
		newuser := structs.User{}

		if err := ctx.ShouldBindJSON(&useradd); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		newuser.Username = useradd.Username
		newuser.Password = useradd.Password
		newuser.UserEmail = useradd.Email

		if err := db.MainDB.Users.AddNew(newuser); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		session.Set("user", newuser.Username+","+newuser.Password)
		session.Set("username" , newuser.Username)
		session.Save()

		ctx.JSON(http.StatusOK, "create")
	})

	server.POST("/loginadmin", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password, ok := api.admin.adminUsers[username]
		if !ok || string(password) != ctx.PostForm("password") {
			ctx.JSON(http.StatusOK, gin.H{"err": "check your password and your username"})
			return
		}
		session := sessions.Default(ctx)
		session.Set("admin", username+","+string(password))
		session.Save()

		ctx.JSON(http.StatusOK, "ok")
	})
	server.GET("/forgetpassword", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "forgetpassword.html", nil)
	})

	server.POST("/forgetpassword", func(ctx *gin.Context) {
		user := structs.User{}
		username := ctx.PostForm("username")
		if err := db.MainDB.Users.GetUser(username, &user); err != nil {
			ctx.HTML(http.StatusOK, "forgetpassword.html", err.Error())
			return
		}
		if err := sendMsg(user.Password, user.UserEmail); err != nil {
			ctx.HTML(http.StatusOK, "forgetpassword.html", "Try in other Time")
			return
		}
		ctx.HTML(http.StatusOK, "forgetpassword.html", "Check your Email")
	})
}

func sendMsg(body, email string) error {
	from := "spprtswearstore@gmail.com"
	pass := "qgkv evar klku nspo "

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Your Password\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{email}, []byte(msg))

	return err
}
