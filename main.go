package main

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ru-rocker/mock-rest-api/parser"
)

func CORSMiddleware(o parser.Options) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", o.AccessControlAllowOrigin)
		c.Header("Access-Control-Allow-Credentials", o.AccessControlAllowCredentials)
		c.Header("Access-Control-Allow-Headers", o.AccessControlAllowHeaders)
		c.Header("Access-Control-Allow-Methods", o.AccessControlAllowMethods)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handler(r parser.Route) gin.HandlerFunc {

	return func(c *gin.Context) {
		regex, _ := regexp.Compile("\\/:[A-Za-z_]+")
		arr := regex.FindAllString(r.Endpoint, -1)
		body := r.Response.Body
		for _, a := range arr {
			val := c.Param(a[2:])
			body = strings.Replace(body, a[1:], val, -1)
		}

		var raw map[string]interface{}
		err := json.Unmarshal([]byte(body), &raw)
		if err != nil {
			panic(err)
		}

		for _, h := range r.Response.Headers {
			c.Header(h.Key, h.Value)
		}
		c.IndentedJSON(r.Response.StatusCode, raw)
	}
}

func main() {
	router := gin.Default()
	c, e := parser.Parse_YAML(getenv("MOCK_CONFIG_FILE", "config/mock.yaml"))
	if e != nil {
		panic(e)
	}
	router.Use(CORSMiddleware(c.Options))

	for _, r := range c.Route {
		method := r.Method
		if method == "GET" {
			router.GET(r.Endpoint, handler(r))
		} else if method == "POST" {
			router.POST(r.Endpoint, handler(r))
		} else if method == "PUT" {
			router.PUT(r.Endpoint, handler(r))
		} else if method == "DELETE" {
			router.DELETE(r.Endpoint, handler(r))
		} else if method == "HEAD" {
			router.HEAD(r.Endpoint, handler(r))
		} else if method == "PATCH" {
			router.PATCH(r.Endpoint, handler(r))
		}
	}
	router.Run(c.Hostname + ":" + c.Port)
}
