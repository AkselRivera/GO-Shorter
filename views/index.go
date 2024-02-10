package views

import (
	"log"

	"github.com/akselrivera/go-shortener/models"
	"github.com/gin-gonic/gin"
)

func IndexView(title string, db map[string]models.Url, errors map[string]string) gin.H {
	if title == "" {
		title = "GO - Shortener"
	}

	log.Println("ERRORS VIEW", errors)
	index := gin.H{
		"title":  title,
		"db":     db,
		"Errors": errors,
	}

	return index
}
