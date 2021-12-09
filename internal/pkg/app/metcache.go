package app

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewMetCache() MetCache { return MetCache{cache: make(map[string]cacheValue)} }

type MetCache struct {
	cache map[string]cacheValue
}

func (c MetCache) GetFromCacheOrLoad(requestURI string) ([]byte, error) {

	cacheValue, isStale := c.getValueFromCache(requestURI)

	if isStale {
		loaded, err := loadDataAndCreateCacheValue(requestURI)
		if err != nil {
			return nil, stacktrace.Propagate(err, "error while loading data for %s", requestURI)
		}
		c.cache[requestURI] = *loaded
		cacheValue = *loaded
	}
	return cacheValue.data, nil
}

func (c MetCache) getValueFromCache(requestURI string) (cacheValue, bool) {

	cacheValue, isValueInCache := c.cache[requestURI]
	isStale := false
	if isValueInCache {
		isCachedValueExpired := cacheValue.expires.Before(time.Now())
		if isCachedValueExpired {
			log.Debugf("cached value for %s expired - must reload", requestURI)
			isStale = true
		} else {
			log.Debugf("cache hit for %s - expires in %s", requestURI, MakeReadable(time.Until(cacheValue.expires)))
		}
	} else {
		log.Debugf("no cached value for %s - must load", requestURI)
		isStale = true
	}
	return cacheValue, isStale
}

func loadDataAndCreateCacheValue(requestURI string) (*cacheValue, error) {

	log.Debugf("loading %s", requestURI)
	url := fmt.Sprintf("https://api.met.no%s", requestURI)
	resp, body, err := HttpGetWithResponse(url)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while getting url %s", url)
	}
	lf := locationForecast{}
	err = json.Unmarshal(body, &lf)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while unmarshalling forecast response")
	}
	updatedAt, err := ParseZuluTime(lf.Properties.Meta.UpdatedAt)
	if err != nil {
		return nil, err
	}
	expires, err := ParseExpiresHeader(resp.Header.Get("Expires"))
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to extract Expires header from response")
	}
	log.Debugf("loaded data - updated at %v (%v ago), expires in %s",
		updatedAt.Format("15:04:05"), MakeReadable(time.Since(updatedAt)), MakeReadable(time.Until(expires)),
	)

	return &cacheValue{data: body, created: time.Now(), updated: updatedAt, expires: expires}, nil
}

type cacheValue struct {
	data    []byte
	created time.Time
	updated time.Time
	expires time.Time
}

type locationForecast struct {
	Properties properties
}

type properties struct {
	Meta meta
}

type meta struct {
	UpdatedAt string `json:"updated_at"`
}
