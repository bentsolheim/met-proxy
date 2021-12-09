package app

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/palantir/stacktrace"
	"net/http"
)

func CreateGinEngine(cache MetCache) *gin.Engine {
	r := gin.Default()

	r.GET("/weatherapi/*path", func(c *gin.Context) {
		requestURI := c.Request.RequestURI
		data, err := cache.GetFromCacheOrLoad(requestURI)
		if err != nil {
			c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, stacktrace.Propagate(err, "error while getting location data from cache")))
			return
		}

		c.Data(http.StatusOK, "application/json", data)
	})

	return r
}
