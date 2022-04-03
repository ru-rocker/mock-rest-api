package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yalp/jsonpath"

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

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func handler(r parser.Route) gin.HandlerFunc {

	return func(c *gin.Context) {
		regex, _ := regexp.Compile("\\/:[A-Za-z_]+")
		arr := regex.FindAllString(r.Endpoint, -1)

		request_body := c.Request.Body
		var request interface{}

		request_value, err := ioutil.ReadAll(request_body)
		if err != nil {
			panic(err)
		}
		_ = json.Unmarshal(request_value, &request)
		if err != nil {
			panic(err)
		}

		for _, resp := range r.Response {

			body := resp.Body
			for _, a := range arr {
				val := c.Param(a[2:])
				body = strings.Replace(body, a[1:], val, -1)
			}

			var raw map[string]interface{}
			err := json.Unmarshal([]byte(body), &raw)
			if err != nil {
				panic(err)
			}

			for _, h := range resp.Headers {
				c.Header(h.Key, h.Value)
			}

			condition := resp.Condition
			if condition.Type == "request_header" {
				header_key := c.Request.Header.Get(condition.Key)
				header_value := c.Request.Header.Values(condition.Key)
				if condition.State == "equal" {
					if contains(header_value, condition.Value) {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "present" {
					if header_key != "" {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "absent" {
					if header_key == "" {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				}
			} else if condition.Type == "request_param" {
				request_param := c.Param(condition.Key)
				if condition.State == "equal" {
					if request_param == condition.Value {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "present" {
					if request_param != "" {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "absent" {
					if request_param == "" {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				}
			} else if condition.Type == "query_param" {
				query, ok := c.GetQueryArray(condition.Key)
				if condition.State == "equal" {
					if contains(query, condition.Value) {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "present" {
					if ok {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "absent" {
					if !ok {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				}
			} else if condition.Type == "request_body" {
				data, _ := jsonpath.Read(request, condition.Key)
				if condition.State == "equal" {
					if data == condition.Value {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "present" {
					if data != nil {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				} else if condition.State == "absent" {
					if data == nil {
						c.IndentedJSON(resp.StatusCode, raw)
						break
					}
				}
			} else {
				c.IndentedJSON(resp.StatusCode, raw)
				break
			}
		}
	}
}

func generateResponseRequestHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
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
