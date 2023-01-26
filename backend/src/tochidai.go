package scsave

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	tochidaiBaseUrl = "https://tochidai.info"
)

func FetchAreaPrice(area string) int64 {
	form := url.Values{}
	form.Add("Word", area)
	form.Add("Process", "Search")

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest(
		"POST",
		tochidaiBaseUrl + "/search/",
		body,
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var res *http.Response
	for (true) {
		res, err = http.DefaultClient.Do(req)
		time.Sleep(time.Second * 20)
		if err == nil {
			break;
		}
		print(err.Error())
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	priceElement := doc.Find("#summary > div.price-box > p.land-price.cm")
	if priceElement.Length() < 1 {
		priceElement = doc.Find("#summary > div.price-box > p.land-price.c")
	}
	if priceElement.Length() < 1 {
		return -1
	}

	stringPrice := priceElement.Text()
	stringPriceList := strings.Split(stringPrice, "万")
	var price int64 = -1
	if len(stringPriceList) > 1 {
		price1, _ := strconv.ParseInt(stringPriceList[0], 10, 64)
		price2, _ := strconv.ParseInt(stringPriceList[1], 10, 64)
		price = price1 * 10000 + price2
		return price
	}
	price, _ = strconv.ParseInt(stringPriceList[0], 10, 64)
	return price
}
