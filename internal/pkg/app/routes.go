package app

import (
	"fmt"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CacheValue struct {
	data    []byte
	created time.Time
}

func CreateGinEngine() *gin.Engine {
	r := gin.Default()

	cache := map[string]CacheValue{}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/met/location-forecast", func(c *gin.Context) {
			lat, _ := c.GetQuery("lat")
			lon, _ := c.GetQuery("lon")
			if !(lat != "" && lon != "") {
				c.JSON(http.StatusBadRequest, rest.WrapResponse(nil, errors.New("both lat and lon request parameters are required")))
				return
			}
			data, err := getFromCacheOrLoad(cache, lat, lon)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, errors.Wrap(err, "error while getting location data from cache")))
				return
			}
			c.Data(http.StatusOK, "application/json", data)
		})
	}

	return r
}

func getFromCacheOrLoad(cache map[string]CacheValue, lat string, lon string) ([]byte, error) {
	locKey := fmt.Sprintf("lat=%s&lon=%s", lat, lon)
	cacheValue, isValueInCache := cache[locKey]

	mustLoadData := false
	if isValueInCache {
		isCachedValueExpired := time.Since(cacheValue.created).Minutes() >= 15
		if isCachedValueExpired {
			mustLoadData = true
		}
	} else {
		mustLoadData = true
	}
	if mustLoadData {
		log.Printf("Getting data for %s", locKey)
		url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/complete?%s", locKey)
		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error while getting %s", url))
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		cacheValue = CacheValue{data: body, created: time.Now()}
		cache[locKey] = cacheValue
	}
	return cacheValue.data, nil
}
