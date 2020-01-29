package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Desc  string `sql:"type:text"`
}

var DB *gorm.DB

func main() {

	var err error

	DB, err := gorm.Open("mysql", "root:@/learngin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()

	// Migrate the schema
	// db.AutoMigrate(&Article{})

	router := gin.Default()

	// grouping end point
	// v1 := router.Group("/api/v1/")
	// {
	// 	v1.GET("/", getHome)
	// 	v1.GET("/article/:title", getArticle)
	// 	v1.POST("/articles", postArticle)
	// }

	// grouping end point and grouping again
	v1 := router.Group("/api/v1/")
	{
		articles := v1.Group("/article")
		{
			articles.GET("/", getHome)
			articles.GET("/:slug", getArticle)
			articles.POST("/", postArticle)
		}
	}

	// single endpoint
	// router.GET("/", getHome)
	// router.GET("/article/:title", getArticle)
	// router.POST("/articles", postArticle)

	router.Run()
}

func getHome(c *gin.Context) {

	items := []Article{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   items,
	})
}

func getArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item Article

	if DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{"status": "error", "message": "record not found"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "sukses",
		"data":   item,
	})
}

func postArticle(c *gin.Context) {

	item := Article{

		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	DB.Create(&item)

	c.JSON(200, gin.H{
		"status": "berhasil post",
		"data":   item,
	})
}
