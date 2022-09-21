package util

import (
	"github.com/google/uuid"
	"log"
	"net/url"
	"strings"
)

func ParseHost(articleUrl string) string {
	parse, err := url.Parse(articleUrl)
	if err != nil {
		return ""
	}
	return parse.Host
}

func GenUuid() string {
	u1, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return u1.String()
}

func GetDirection(country string) (newCountry string) {
	countryList := []string{"china_hk", "china_tw", "china_macao", "china_xz"}
	for _, cou := range countryList {
		if strings.Contains(country, cou) {
			newCountry = strings.Replace(country, "china_hk", "hk", -1)
			return strings.ToUpper(newCountry)
		}
	}
	return
}
