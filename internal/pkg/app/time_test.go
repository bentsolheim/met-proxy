package app

import "testing"

func TestParseExpiresHeader(t *testing.T) {

	parsedTime, _ := ParseExpiresHeader("Wed, 19 Aug 2020 17:37:48 GMT")
	println(parsedTime.String())
}

func TestParseExpiresHeader2(t *testing.T) {

	_, err := ParseExpiresHeader("Wed, 19 August 2020 17:37:48 GMT")
	println(err.Error())
}
