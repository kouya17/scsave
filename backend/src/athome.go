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
	atHomeBaseUrl = "https://suumo.jp"
)

func FetchPropertyFromUrlAtHome(url string, args string) (Property, error) {
	p := Property{}
	p.Url = url

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.AddCookie(&http.Cookie{
		Name: "reese84", Value: args, Path: "/",
	})

	res, err := http.DefaultClient.Do(req)
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

	breadcrumbs := doc.Find("#breadcrumbs > li")
	p.City = strings.Fields(breadcrumbs.Eq(4).Text())[0]

	tables := doc.Find("#item-detail_data > div > div > section > div > div.left > table")
	tables.Each(func(i int, s *goquery.Selection) {
		trs := s.Find("tr")
		trs.Each(func(i int, s *goquery.Selection) {
			ths := s.Find("th")
			tds := s.Find("td")
			ths.Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.Text(), "交通") {
					p.Access = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
					accessArray := strings.Fields(tds.Eq(i).Text())
					for _, str := range accessArray {
						if strings.Contains(str, "駅") {
							p.Station = str[:strings.Index(str, "駅")]
							return
						}
					}
				}
				if strings.Contains(s.Text(), "価格") {
					rep := regexp.MustCompile(`[^0-9]`)
					price, _ := strconv.ParseUint(rep.ReplaceAllString(tds.Eq(i).Text(), ""), 10, 32)
					if (price != 0) {
						p.Price = uint32(price)
						return
					}
				}
				if strings.Contains(s.Text(), "間取り") {
					p.Layout = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if strings.Contains(s.Text(), "建物面積") {
					buildingAreaString := strings.Split(tds.Eq(i).Text(), "m")[0]
					rep := regexp.MustCompile(`[^0-9\.]`)
					buildingArea, _ := strconv.ParseFloat(rep.ReplaceAllString(buildingAreaString, ""), 32)
					p.BuildingArea = float32(buildingArea)
					return
				}
				if strings.Contains(s.Text(), "土地面積") {
					landAreaString := strings.Split(tds.Eq(i).Text(), "m")[0]
					rep := regexp.MustCompile(`[^0-9\.]`)
					landArea, _ := strconv.ParseFloat(rep.ReplaceAllString(landAreaString, ""), 32)
					p.LandArea = float32(landArea)
					return
				}
				if strings.Contains(s.Text(), "築年月") {
					buildYear, _ := strconv.ParseUint(strings.Split(tds.Eq(i).Text(), "年")[0], 10, 16)
					p.BuildYear = uint16(buildYear)
					return
				}
				if strings.Contains(s.Text(), "建物構造") {
					p.Structure = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if strings.Contains(s.Text(), "土地権利") {
					p.Rights = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if strings.Contains(s.Text(), "用途地域") {
					p.AreaPurpose = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if strings.Contains(s.Text(), "接道状況") {
					p.Road = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
					return
				}
				if strings.Contains(s.Text(), "容積率") { // TODO: 建ぺい率もいれる
					p.CoverageRatio = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
					return
				}
				if strings.Contains(s.Text(), "地目") {
					p.LandKind = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if strings.Contains(s.Text(), "引渡可能時期") {
					p.Timing = strings.TrimSpace(tds.Eq(i).Text())
					return
				}
				if s.Text() == "備考" {
					p.OtherNotice = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
					return
				}
			})
		})
	})

	return p, nil
}

func FetchPropertyFromSearchPageAtHome(url string) []Property {
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
