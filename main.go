package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akselrivera/go-shortener/models"
	"github.com/akselrivera/go-shortener/utils"
	"github.com/akselrivera/go-shortener/views"
	"github.com/gin-gonic/gin"
)

var db = make(map[string]models.Url, 0)
var errors = make(map[string]string, 0)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", views.IndexView("", db, nil))
	})

	router.POST("/", func(c *gin.Context) {
		url := c.PostForm("url")
		var statusResponse int
		var shortenedUrl string

		log.Println(db[url], db[url] != (models.Url{}))
		if len(url) <= 3 {
			errors["url"] = "Invalid url"
			log.Println(errors)
			statusResponse = http.StatusBadRequest

		} else if _, ok := db[url]; ok {
			errors["url"] = "Url already exists"
			log.Println("this url already exists", errors)
			statusResponse = http.StatusConflict

		} else {
			statusResponse = http.StatusFound
			var hash string = utils.GetHash()
			shortenedUrl = fmt.Sprintf("%s/%s", c.Request.Host, hash)

			var newUrl models.Url

			newUrl.Url = shortenedUrl
			newUrl.Clicks = 0
			newUrl.Expiration = time.Now().Add(48 * time.Hour)
			newUrl.Hash = hash

			db[url] = newUrl

		}

		if statusResponse != http.StatusFound {
			log.Println("DEFAULT", errors)
			c.HTML(http.StatusConflict, "index.html", views.IndexView("", db, errors))
		} else {
			c.Redirect(http.StatusFound, "/tracking?url="+shortenedUrl)

		}

	})

	router.GET("/tracking", func(c *gin.Context) {
		url := c.Query("url")
		log.Println(url)
		log.Println()
		if url != "" {
			var originalUrl string

			for k, v := range db {
				if v.Url == url {
					originalUrl = k
				}
			}

			if originalUrl == "" {
				c.Redirect(http.StatusFound, "/tracking")
			} else {

				formatedDate := db[originalUrl].Expiration.Format("02/01/2006 15:04:05 MST")

				c.HTML(http.StatusOK, "tracking-query.html", gin.H{
					"title":       "GO - Shortener",
					"url":         url,
					"originalUrl": originalUrl,
					"clicks":      db[originalUrl].Clicks,
					"expiration":  formatedDate,
				})
			}

		} else {

			c.HTML(http.StatusOK, "tracking.html", gin.H{
				"title": "GO - Shortener",
			})
		}
	})

	router.POST("/tracking", func(c *gin.Context) {
		url := c.PostForm("url")

		if len(url) <= 3 {
			c.HTML(http.StatusBadRequest, "tracking.html", gin.H{
				"title": "GO - Shortener",
				"db":    db,
			})
		} else {
			c.Redirect(http.StatusFound, "/tracking?url="+url)
		}

	})

	router.GET("/:hash", func(c *gin.Context) {
		hash := c.Param("hash")

		for originalUrl, v := range db {
			if v.Hash == hash {
				v.Clicks = v.Clicks + 1
				db[originalUrl] = v
				log.Println(db[originalUrl].Clicks)
				c.Redirect(http.StatusFound, originalUrl)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
			"db":    db,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
        log.Panicf("error: %s", err)
	}
}
