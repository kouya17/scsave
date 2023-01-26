package scsave

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseUrl = "https://suumo.jp"
)

func FetchPropertyFromUrlSuumo(url string) (Property, error) {
	p := Property{}
	p.Url = url

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return p, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	topContents := doc.Find("#topContents > div.cf.mt10 > div.fl.w420 > div.mt9 > p.b")
	topContents.Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "「") {
			p.Access = s.Text()
			rep := regexp.MustCompile(`「([^」]*)」`)
			p.Station = rep.FindStringSubmatch(s.Text())[1]
		}
	})
	p.City = doc.Find("#topContents > div.cf.mt10 > div.fl.w420 > div.mt9 > p.mt5.b").First().Text()

	detailList := doc.Find("table.mt10.bdGrayT.bdGrayL.bgWhite.pCell10.bdclps.wf > tbody > tr")
	detailList.Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "価格") {
			tds := s.Find("td")
			rep := regexp.MustCompile(`[^0-9]`)
			price, _ := strconv.ParseUint(rep.ReplaceAllString(tds.First().Text(), ""), 10, 32)
			if (price != 0){
				p.Price = uint32(price)
			}
		}
		if strings.Contains(s.Text(), "私道負担・道路") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			p.Road = strings.TrimSpace(tds.Eq(0).Text())
			p.OtherCost = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "建物面積") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			p.Layout = strings.TrimSpace(tds.Eq(0).Text())
			buildingArea, _ := strconv.ParseFloat(strings.Split(tds.Eq(1).Text(), "m")[0], 32)
			p.BuildingArea = float32(buildingArea)
		}
		if strings.Contains(s.Text(), "土地面積") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			landArea, _ := strconv.ParseFloat(strings.Split(tds.Eq(0).Text(), "m")[0], 32)
			p.LandArea = float32(landArea)
			p.CoverageRatio = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "完成時期") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			buildYear, _ := strconv.ParseUint(strings.Split(tds.Eq(0).Text(), "年")[0], 10, 16)
			p.BuildYear = uint16(buildYear)
			p.Timing = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "土地の権利形態") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			p.Rights = strings.TrimSpace(tds.Eq(0).Text())
			p.Structure = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "施工") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			p.BuildCompany = strings.TrimSpace(tds.Eq(0).Text())
			p.Reform = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "用途地域") {
			tds := s.Find("td")
			if tds.Length() < 2 {
				return
			}
			p.AreaPurpose = strings.TrimSpace(tds.Eq(0).Text())
			p.LandKind = strings.TrimSpace(tds.Eq(1).Text())
		}
		if strings.Contains(s.Text(), "その他制限事項") {
			tds := s.Find("td")
			p.OtherRestriction = strings.TrimSpace(tds.Eq(0).Text())
		}
		if strings.Contains(s.Text(), "その他概要・特記事項") {
			tds := s.Find("td")
			p.OtherNotice = strings.TrimSpace(tds.Eq(0).Text())
		}
	})
	return p, nil
}

func FetchPropertyFromSearchPageSuumo(url string) []Property {
	ps := []Property{}
	for {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

    propertyUrls := doc.Find("h2.property_unit-title > a").Map(func(i int, s *goquery.Selection) string {
      return s.AttrOr("href", "")
    })
    for _, v := range propertyUrls {
      if v == "" {
        continue
      }
      property, err := FetchPropertyFromUrlSuumo(baseUrl + v)
			if err != nil {
				continue
			}
      fmt.Printf("p: %#v\n", property)
      ps = append(ps, property)
      time.Sleep(time.Second * 10)
    }
    
    paginationPartsList := doc.Find("p.pagination-parts")
    nextUrl := ""
    paginationPartsList.Each(func(i int, s *goquery.Selection) {
      aElement := s.Find("a")
      if strings.Contains(aElement.Text(), "次へ") {
        nextUrl = aElement.AttrOr("href", "")
      }
    })
    if nextUrl == "" {
      break
    }
    //print("next page")
    url = baseUrl + nextUrl
	}
  return ps
}
