package app

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

type CacheValue struct {
	data    []byte
	created time.Time
}

func NewMetCache(skipTlsVerification bool) MetCache {

	var client *http.Client
	if skipTlsVerification {
		log.Printf("TLS Verification has been disabled for requests to api.met.no")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}
	return MetCache{cache: make(map[string]CacheValue), client: client}
}

type MetCache struct {
	cache  map[string]CacheValue
	client *http.Client
}

func (c MetCache) GetFromCacheOrLoad(requestURI string) ([]byte, error) {
	cacheValue, isValueInCache := c.cache[requestURI]

	mustLoadData := false
	if isValueInCache {
		valueAge := time.Since(cacheValue.created).Minutes()
		log.Printf("Cached value age: %f min", valueAge)
		isCachedValueExpired := valueAge >= 15.0
		if isCachedValueExpired {
			mustLoadData = true
		}
	} else {
		mustLoadData = true
	}
	if mustLoadData {
		log.Printf("Getting data for %s", requestURI)
		url := fmt.Sprintf("https://api.met.no%s", requestURI)
		resp, err := c.client.Get(url)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error while getting %s", url))
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		updatedAt, err := getUpdatedAtFromResponse(body)
		if err != nil {
			return nil, err
		}

		log.Printf(fmt.Sprintf("Remotely updated at %v (%v ago)", updatedAt.Format("2006-01-02 15:04:05"), humanizeDuration(time.Since(updatedAt))))

		cacheValue = CacheValue{data: body, created: time.Now()}
		c.cache[requestURI] = cacheValue
	}
	return cacheValue.data, nil
}

func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"d", days},
		{"h", hours},
		{"m", minutes},
		{"s", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		default:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}

func getUpdatedAtFromResponse(body []byte) (time.Time, error) {
	lf := LocationForecast{}
	err := json.Unmarshal(body, &lf)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "error while unmarshalling forecast response")
	}
	layout := "2006-01-02T15:04:05Z"
	updatedAt, err := time.Parse(layout, lf.Properties.Meta.UpdatedAt)
	if err != nil {
		return time.Time{}, errors.Wrap(err, fmt.Sprintf("error while parsing updated_at string (%s)", lf.Properties.Meta.UpdatedAt))
	}
	tz, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return time.Time{}, errors.Wrap(err, "error while getting current timezone")
	}
	return updatedAt.In(tz), nil
}

type LocationForecast struct {
	Properties Properties
}

type Properties struct {
	Meta Meta
}

type Meta struct {
	UpdatedAt string `json:"updated_at"`
}
