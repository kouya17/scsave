package scsave

import "strings"

func GetProperties(url string, args string) ([]Property, string) {
	if strings.Contains(url, "suumo") {
		return FetchPropertyFromSearchPageSuumo(url), "suumo"
	}
	if strings.Contains(url, "housedo") {
		return FetchPropertyFromSearchPageHouseDo(url, args), "housedo"
	}
	if strings.Contains(url, "nifty") {
		return FetchPropertyFromSearchPageNifty(url), "nifty"
	}
	return []Property{}, "unknown"
}