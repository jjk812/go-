package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"goPro/dao"
	"goPro/logic"
	"goPro/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var db *gorm.DB
var ctx context.Context
var rdb *redis.Client

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@/mytest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{})
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // 无密码
		DB:       0,                // 默认数据库
	})
}
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}
func ArticlesNew(c *gin.Context) {
	c.HTML(http.StatusOK, "ArticlesNew.html", nil)
}
func ArticlesDelete(c *gin.Context) {
	c.HTML(http.StatusOK, "ArticlesDelete.html", nil)
}
func ArticlesEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "ArticlesEdit.html", nil)
}
func DoEditArticles(c *gin.Context) {
	var post model.Post
	postIDStr := c.PostForm("PostId")
	post.Content = c.PostForm("Content")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	post.ID = uint(postID)
	username, _ := c.Get("username")
	post.Username = username.(string)
	success := dao.UpdatePost(post)
	if success {
		c.JSON(200, "Update Success!!")
	} else {
		c.JSON(500, "你不拥有或该篇post或该篇post不存在")
	}
}
func DoDeleteArticles(c *gin.Context) {
	var post model.Post
	postIDStr := c.PostForm("PostId")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	fmt.Println("postID", postID)
	post.ID = uint(postID)
	username, _ := c.Get("username")
	post.Username = username.(string)

	rows := dao.DeletePost(post)
	if rows > 0 {
		c.JSON(200, "Delete Success!!")
	} else {
		c.JSON(500, "没有该篇post")
	}
}
func DoCreatArticles(c *gin.Context) {
	var post model.Post
	post.Title = c.PostForm("Title")
	post.Content = c.PostForm("Content")
	username, _ := c.Get("username")
	post.Username = username.(string)
	dao.AddPost(post)
	c.JSON(200, "Add Success!!")
}
func DoSearchArticles(c *gin.Context) {

	Title := c.PostForm("Title")

	posts, _ := dao.SearchPost(Title)
	c.HTML(http.StatusOK, "ArticlesSearch.html", gin.H{
		"posts": posts,
	})
}
func DoRegister(c *gin.Context) {
	var user model.User

	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	fmt.Println(user)
	status := logic.InsertUser(user)
	if status == 1 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"SuccessMessage": "注册成功，请登录",
		})
	} else {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"ErrorMessage": "注册失败，用户名已存在",
		})
	}
}

func Index(c *gin.Context) {
	var user model.User
	user.Username = c.PostForm("username")
	user.Password = c.PostForm("password")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	success := logic.TestPassword(user)
	if success {
		c.SetCookie("login_user", user.Username, 3600, "/", "", false, true)
		currentPagePosts, pagination := logic.Pages(page)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"posts":      currentPagePosts,
			"pagination": pagination,
		})
	} else {
		c.JSON(401, "Login Failed")
	}
}

func UserHome(c *gin.Context) {
	username, _ := c.Get("username")
	alertMessage := c.Query("alert")
	posts, _ := dao.SearchUserPost(username.(string))
	imags, _ := dao.SearchImage(username.(string))
	c.HTML(http.StatusOK, "UserHome.html", gin.H{
		"posts":        posts,
		"images":       imags,
		"AlertMessage": alertMessage,
	})

}
func PostDetail(c *gin.Context) {
	id := c.Query("id")
	postid, _ := strconv.Atoi(id)
	post, _ := dao.SearchPostID(postid)
	rdb.ZIncrBy(ctx, "article_rank", 1, id)
	score, _ := rdb.ZScore(ctx, "article_rank", id).Result()
	c.HTML(http.StatusOK, "PostDetail.html", gin.H{
		"post":  post,
		"score": score,
	})
}

func TopArticles(c *gin.Context) {
	var posts []model.Post

	// 从Redis获取文章ID的有序列表
	ids, err := rdb.ZRevRange(ctx, "article_rank", 0, 9).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var intIDs []int
	// 文章评分映射
	scores := make(map[int]float64)
	for _, idStr := range ids {
		idInt, _ := strconv.Atoi(idStr)
		intIDs = append(intIDs, idInt)

		// 获取每篇文章的评分
		score, err := rdb.ZScore(ctx, "article_rank", idStr).Result()
		if err == nil {
			scores[idInt] = score
		}

	}

	// 查询MySQL数据库
	db.Where("id IN (?)", intIDs).Find(&posts)

	// 创建一个包含文章和评分的结构
	type ArticleWithScore struct {
		Post  model.Post
		Score float64
	}

	// 组合文章和它们的评分
	var articlesWithScores []ArticleWithScore
	for _, post := range posts {
		articlesWithScores = append(articlesWithScores, ArticleWithScore{
			Post:  post,
			Score: scores[int(post.ID)],
		})
	}

	// 传递 ArticlesWithScores 到模板
	c.HTML(http.StatusOK, "TopArticles.html", gin.H{
		"postsScores": articlesWithScores,
	})
}
func Pages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	currentPagePosts, pagination := logic.Pages(page)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":      currentPagePosts,
		"pagination": pagination,
	})
}
func VerifyLogin(c *gin.Context) {
	// 检查是否存在名为"login_user"的cookie
	login_user, err := c.Cookie("login_user")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  "未经认证的访问,请先登录",
		})
		c.Abort()
		return
	}
	// 解析用户ID，例如假设login_user就是用户ID
	username := login_user
	// 将用户ID添加到gin.Context中，可以在后续的处理程序中获取
	c.Set("username", username)

	c.Next()
}
func Upload(c *gin.Context) {
	var image model.Image
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取当前目录"})
		return
	}

	// 创建上传目录路径
	savePath := filepath.Join(wd, "uploads", filepath.Base(file.Filename))
	// 检查uploads目录是否存在
	if _, err := os.Stat(filepath.Join(wd, "uploads")); os.IsNotExist(err) {
		// 如果不存在，创建目录
		os.MkdirAll(filepath.Join(wd, "uploads"), os.ModePerm)
	}
	// 保存文件到之前构建的路径
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	username, _ := c.Get("username")
	fmt.Println(username)
	image.ImgURL = filepath.Base(file.Filename)
	image.Username = username.(string)

	dao.AddImage(image)
	c.Redirect(http.StatusFound, "/UserHome?alert=上传成功")
}
