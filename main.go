package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akselrivera/go-shortener/utils"
	"github.com/gin-gonic/gin"
)

type Url struct {
	Url        string
	Clicks     int
	Expiration time.Time
	hash       string
}

var db = make(map[string]Url, 0)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8081"
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "GO - Shortener",
			"db":    db,
		})
	})

	router.POST("/", func(c *gin.Context) {
		url := c.PostForm("url")

		if len(url) <= 3 {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"title": "Main website",
				"db":    db,
			})
		} else if db[url] != (Url{}) {
			c.HTML(http.StatusConflict, "index.html", gin.H{
				"title": "GO - Shortener",
				"db":    db,
			})
		} else {

			var hash string = utils.GetHash()
			var shortenedUrl string = fmt.Sprintf("%s/%s", c.Request.Host, hash)

			var newUrl Url

			newUrl.Url = shortenedUrl
			newUrl.Clicks = 0
			newUrl.Expiration = time.Now().Add(48 * time.Hour)
			newUrl.hash = hash

			db[url] = newUrl

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
			if v.hash == hash {
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
	router.Run(port) // listen and serve on 0.0.0.0:8080
}
