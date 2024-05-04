package router

import (
	"github.com/gin-gonic/gin"
	"goPro/controller"
	"html/template"
)

func Routers() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"add":      add,
		"subtract": subtract,
	})
	router.Static("/uploads", "./uploads")
	router.LoadHTMLGlob("templates/*")
	router.GET("/register", controller.Register)
	router.POST("/doregister", controller.DoRegister)
	router.POST("/index", controller.Index)
	router.GET("/index", controller.Pages)
	router.GET("/login", controller.Login)

	router.GET("/articles/new", controller.ArticlesNew)
	router.POST("/DoCreatArticles", controller.VerifyLogin, controller.DoCreatArticles)
	router.GET("/articles/delete", controller.ArticlesDelete)
	router.POST("/DoDeleteArticles", controller.VerifyLogin, controller.DoDeleteArticles)
	router.POST("/DoSearchArticles", controller.VerifyLogin, controller.DoSearchArticles)
	router.GET("/articles/edit", controller.ArticlesEdit)
	router.POST("/DoEditArticles", controller.VerifyLogin, controller.DoEditArticles)
	router.GET("/", controller.Pages)
	router.GET("/UserHome", controller.VerifyLogin, controller.UserHome)
	router.GET("/PostDetail", controller.PostDetail)
	router.POST("/upload", controller.VerifyLogin, controller.Upload)
	router.GET("/TopArticles", controller.TopArticles)
	router.Run(":8080")

}
func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}
