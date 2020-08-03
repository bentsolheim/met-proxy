package app

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func CreateGinEngine(cache MetCache) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/met/location-forecast", func(c *gin.Context) {
			lat, _ := c.GetQuery("lat")
			lon, _ := c.GetQuery("lon")
			if !(lat != "" && lon != "") {
				c.JSON(http.StatusBadRequest, rest.WrapResponse(nil, errors.New("both lat and lon request parameters are required")))
				return
			}
			data, err := cache.GetFromCacheOrLoad(lat, lon)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, errors.Wrap(err, "error while getting location data from cache")))
				return
			}
			c.Data(http.StatusOK, "application/json", data)
		})
	}

	return r
}
