package api

import (
	"db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"structs"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
)

type admin struct {
	adminGroup *gin.RouterGroup
	adminUsers map[string]string
}

func (admin *admin) GetUsers() {

	f, _ := os.OpenFile("adminuser.json", os.O_RDONLY|os.O_CREATE, 0660)
	data, _ := io.ReadAll(f)

	if err := json.Unmarshal(data, &admin.adminUsers); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func (admin *admin) setOrderApi() {

	admin.adminGroup.GET("/order", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "order.html", nil)
	})

	admin.adminGroup.POST("/order", func(ctx *gin.Context) {
		orders := db.MainDB.Orders.Get()
		if orders == nil {
			ctx.JSON(http.StatusOK, "")
			return
		}
		ctx.JSON(http.StatusOK, orders)
	})
	admin.adminGroup.GET("/outstock", func(ctx *gin.Context) {
		dataOfoutStock := db.MainDB.OutStock.Get()
		if dataOfoutStock == nil {
			ctx.JSON(http.StatusOK, "no date")
		} else {
			ctx.JSON(http.StatusOK, dataOfoutStock)
		}
	})

	admin.adminGroup.DELETE("/outstock", func(ctx *gin.Context) {
		Id := ctx.GetInt("id")
		db.MainDB.OutStock.Delete(Id)
		ctx.JSON(http.StatusOK, fmt.Sprintf("delete %d", Id))
	})
}

func (admin *admin) setProductApi() {

	admin.adminGroup.PATCH("/product/model/", func(ctx *gin.Context) {
		type UpateModel struct {
			Model     structs.Model `json:"model"`
			Id        int           `json:"id"`
			Container string        `json:"container"`
			Kind      string        `json:"kind"`
		}
		model := UpateModel{}
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}
		if err := db.MainDB.Stock.UpateModel(model.Id, model.Container, model.Kind, &model.Model); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "Updata Id")

	})
	admin.adminGroup.POST("/product/model/updataImage", func(ctx *gin.Context) {

		type LinkImage struct {
			Id        int      `json:"id"`
			Container string   `json:"container"`
			Kind      string   `json:"kind"`
			Image     []string `json:"image"`
		}
		linkImage := LinkImage{}

		if err := ctx.ShouldBindJSON(&linkImage); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}
		if err := db.MainDB.Stock.UpateLinkesImage(linkImage.Id, linkImage.Container, linkImage.Kind, linkImage.Image); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return

		}
		ctx.JSON(http.StatusOK, "UPDATA")

	})

	admin.adminGroup.POST("/product/model/sizes", func(ctx *gin.Context) {
		type Size struct {
			Id        int                       `json:"id"`
			SizeName  string                    `json:"sizeName"`
			Container string                    `json:"container"`
			Kind      string                    `json:"kind"`
			C         map[string]map[string]int `json:"Size"`
		}
		newSize := Size{}
		if err := db.MainDB.Stock.UpdataSizeFromModel(newSize.Id, newSize.Container, newSize.Kind, newSize.SizeName, newSize.C[newSize.SizeName]); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "done")
	})

	admin.adminGroup.GET("/product", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "manage.html", nil)
	})
	admin.adminGroup.POST("/product/container", func(ctx *gin.Context) {
		newContainer := ctx.PostForm("Container")
		if len(newContainer) == 0 {
			ctx.JSON(http.StatusOK, "name container is empty")
		} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
			ctx.JSON(http.StatusOK, "has space")
		} else if err := db.MainDB.Stock.AddNewContainer(newContainer); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
		} else {
			ctx.JSON(http.StatusOK, "Created")
		}

	})

	admin.adminGroup.GET("/product/container", func(ctx *gin.Context) {
		Containers := db.MainDB.Stock.GetAllContainerAndKind()
		ctx.JSON(http.StatusOK, &Containers)

	})

	admin.adminGroup.POST("/product/kind", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var newKind kind
		if err := ctx.ShouldBindJSON(&newKind); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		if err := db.MainDB.Stock.NewKindtoContainer(newKind.Container, newKind.Kindname); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "Created")

	})

	admin.adminGroup.DELETE("/product/kind", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var deleteKind kind
		if err := ctx.ShouldBindJSON(&deleteKind); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		if err := db.MainDB.Stock.DeleteKind(deleteKind.Container, deleteKind.Kindname); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "delete")

	})
	admin.adminGroup.DELETE("/product/container", func(ctx *gin.Context) {
		Container, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusOK, "try agein")
		}
		newContainer := string(Container)
		if len(newContainer) == 0 {
			ctx.JSON(http.StatusOK, "name container is empty")
		} else if regexp.MustCompile(`\s`).Match([]byte(newContainer)) {
			ctx.JSON(http.StatusOK, "has space")
		} else if err := db.MainDB.Stock.DeleteContainer(newContainer); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
		} else {
			ctx.JSON(http.StatusOK, "delete")
		}

	})

	admin.adminGroup.PUT("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
		}
		var ki kind
		if err := ctx.ShouldBindJSON(&ki); err != nil {
			if io.EOF == err {
				ctx.JSON(http.StatusOK, "empty Body")
			} else {
				ctx.JSON(http.StatusOK, err.Error())
			}
			return
		}
		fg, err := db.MainDB.Stock.GetAllModelsInKind(ki.Container, ki.Kindname)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"err": err.Error()})
			return
		}
		if fg == nil {
			ctx.JSON(http.StatusOK, gin.H{"err": "no data"})
			return
		}
		ctx.JSON(http.StatusOK, &fg)
	})
	admin.adminGroup.POST("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string        `json:"container"`
			Kindname  string        `json:"kind"`
			Model     structs.Model `json:"model"`
		}
		var ki kind
		err := ctx.ShouldBindJSON(&ki)
		if err != nil {
			ctx.JSON(http.StatusOK, "check json ")
			return
		}
		if err := db.MainDB.Stock.AddModelToKind(ki.Container, ki.Kindname, &ki.Model); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, "create")

	})
	admin.adminGroup.DELETE("/product/model", func(ctx *gin.Context) {
		type kind struct {
			Container string `json:"container"`
			Kindname  string `json:"kind"`
			Id        int    `josn:"id"`
		}
		deleteModel := kind{}
		if err := ctx.ShouldBindJSON(&deleteModel); err != nil {
			ctx.JSON(http.StatusOK, "check json")
			return
		}
		if err := db.MainDB.Stock.DeleteModelFromKind(deleteModel.Container, deleteModel.Kindname, deleteModel.Id); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, deleteModel.Id)
	})

}

func (admin *admin) setLogoutApi() {

	admin.adminGroup.GET("/logout", func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		s.Clear()
		s.Save()

		ctx.HTML(http.StatusOK, "loginadmin.html", nil)
	})
	admin.adminGroup.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "adminMainPage.html", nil)
	})
}
