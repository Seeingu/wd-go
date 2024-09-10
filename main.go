package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/tmaxmax/go-sse"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var htmlFilePath string
var rootCmd = &cobra.Command{
	Use:   "wd <any html file>",
	Short: "watch html file change",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		unknownFile := args[0]
		if !strings.HasSuffix(unknownFile, ".html") {
			cmd.Usage()
			os.Exit(1)
		}
		_, err := os.Stat(unknownFile)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%s does not exist\n", unknownFile)
			os.Exit(1)
		}
		htmlFilePath = unknownFile
	},
}

func registerApi() *gin.Engine {
	apiEngine := gin.New()
	apiG := apiEngine.Group("/api")
	{
		apiG.GET("/reload.js", func(c *gin.Context) {
			content, _ := os.ReadFile("static/reload.js")
			c.Data(http.StatusOK, "application/javascript", content)
		})
		// Publish example response
		apiG.POST("/publish", func(c *gin.Context) {
			m := &sse.Message{}
			m.AppendData("Hello world!")
			err := sseHandler.Publish(m)
			if err != nil {
				c.Status(http.StatusInternalServerError)
			} else {
				c.Status(http.StatusOK)
			}
		})
		// Force page reload
		apiG.GET("/publish/reload", func(c *gin.Context) {
			SSEReload()
			c.Status(http.StatusOK)
		})
		apiG.GET("/events", gin.WrapH(sseHandler))
	}
	return apiEngine
}

const DefaultPort = 3012

func startWatch() {
	err := Watch(htmlFilePath, func() {
		SSEReload()
	})
	if err != nil {
		log.Println("Watch file change failed")
		log.Fatal(err)
	}
}

var corsEnabled = cors.New(cors.Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"Origin"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
})

func mainRouteHandler(c *gin.Context) {
	apiEngine := registerApi()
	p := c.Param("any")
	if strings.HasPrefix(p, "/api") {
		apiEngine.HandleContext(c)
		return
	}

	// default html index
	if p == "/" {
		p = htmlFilePath
	}
	// Read file relative to current dir
	contents, err := os.ReadFile(path.Join(".", p))
	if err != nil {
		msg :=
			fmt.Sprintf("failed to read file %q: %v", p, err)
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(
			msg,
		))
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", contents)
}

func parseArgs() int {
	port := rootCmd.Flags().IntP("port", "p", DefaultPort, "http server port")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	return *port
}

// go run . fixture1.html --port 3013
func main() {
	port := parseArgs()
	go startWatch()

	r := gin.Default()
	r.Use(corsEnabled)
	r.GET("/*any", mainRouteHandler)
	r.Run(":" + strconv.Itoa(port))
}
