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
	Desc  string `sql:"type:text;"`
}

var DB *gorm.DB

func main() {
	var err error

	// Connection to Database (Mysql)
	DB, err = gorm.Open("mysql", "root:root-password@/learn-gin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()

	// Migrate the schema
	DB.AutoMigrate(&Article{})

	router := gin.Default()

	// Grouping alamat API
	v1 := router.Group("/api/v1")
	{
		article := v1.Group("/article")
		{
			article.GET("/", getHome)
			article.GET("/:slug", getArticle)
			article.POST("/", postArticle)
		}
	}

	router.Run()
}

func getHome(c *gin.Context) {

	items := []Article{}
	DB.Find(&items)

	c.JSON(200, gin.H{
		"status":  "Berhasil",
		"message": "Berhasil akses home",
		"data":    items,
	})
}

func getArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item Article

	if DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "record not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status":  "Berhasil",
		"message": "Berhasil akses home",
		"data":    item,
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
		"status": "Berhasil",
		"data" : item,
	})
}
