package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const filesDir = ""
const APIToken = ""
const base = ""

func main() {
	files, err := getFiles()
	if err != nil {
		log.Fatal(err)
	}
	runAPI(files)
}

func getFiles() (files []string, err error) {
	fileLists, err := ioutil.ReadDir(fmt.Sprintf("%s/", filesDir))
	if err != nil {
		return []string{}, err
	}
	for _, file := range fileLists {
		if !file.IsDir() && strings.Contains(file.Name(), "mp3") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func runAPI(files []string) {
	r := gin.Default()

	apiGroup := r.Group("api")
	{
		apiGroup.GET("/", func(c *gin.Context) {
			token := c.GetHeader("token")
			if token == "" {
				c.JSON(http.StatusUnauthorized, "no auth token provided")
				return
			}
			if token != APIToken {
				c.JSON(http.StatusUnauthorized, "invalid auth token")
				return
			}
			query, _ := c.GetQuery("q")
			var urls []string
			for _, file := range files {
				if query != "" || strings.Contains(strings.ToLower(file), strings.ToLower(query)) {
					u, _ := url.Parse(fmt.Sprintf("%s/%s", base, file))
					urls = append(urls, u.String())
				}
			}
			c.JSON(
				http.StatusOK,
				map[string]interface{}{
					"urls":  urls,
					"total": len(urls),
				},
			)
			return
		})
	}

	r.Run(":9009")
}
