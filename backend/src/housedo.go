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
	houseDoBaseUrl = "https://www.housedo.com"
)

func FetchPropertyFromUrlHouseDo(url string) (Property, error) {
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

	breadcrumbs := doc.Find("#topic_path > li")
	p.City = strings.Replace(breadcrumbs.Eq(3).Text() + breadcrumbs.Eq(4).Text(), ">", "", -1)

	trs := doc.Find("#main > table.tbl01 > tbody > tr")
	trs.Each(func(i int, s *goquery.Selection) {
		ths := s.Find("th")
		tds := s.Find("td")
		ths.Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "価格") {
				rep := regexp.MustCompile(`[^0-9]`)
				price, _ := strconv.ParseUint(rep.ReplaceAllString(tds.Eq(i).Text(), ""), 10, 32)
				if (price != 0) {
					p.Price = uint32(price)
					return
				}
			}
			if strings.Contains(s.Text(), "アクセス") {
				p.Access = strings.TrimSpace(tds.Eq(i).Text())
				rep := regexp.MustCompile(`『([^』]*)』`)
				subMatch := rep.FindStringSubmatch(tds.Eq(i).Text())
				if len(subMatch) > 1 {
					p.Station = strings.TrimSpace(strings.Replace(rep.FindStringSubmatch(tds.Eq(i).Text())[1], "駅", "", -1))
					return
				}
				rep = regexp.MustCompile(`「([^」]*)」`)
				subMatch = rep.FindStringSubmatch(tds.Eq(i).Text())
				if len(subMatch) > 1 {
					p.Station = strings.TrimSpace(strings.Replace(rep.FindStringSubmatch(tds.Eq(i).Text())[1], "駅", "", -1))
					return
				}
			}
			if strings.Contains(s.Text(), "間取り") {
				p.Layout = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "築年月") {
				buildYear, _ := strconv.ParseUint(strings.Split(tds.Eq(i).Text(), "年")[0], 10, 16)
				p.BuildYear = uint16(buildYear)
				return
			}
			if strings.Contains(s.Text(), "土地面積") {
				landAreaString := strings.Split(tds.Eq(i).Text(), "m")[0]
				rep := regexp.MustCompile(`[^0-9\.]`)
				landArea, _ := strconv.ParseFloat(rep.ReplaceAllString(landAreaString, ""), 32)
				p.LandArea = float32(landArea)
				return
			}
			if strings.Contains(s.Text(), "建物構造") {
				p.Structure = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "建物面積") {
				buildingAreaString := strings.Split(tds.Eq(i).Text(), "m")[0]
				rep := regexp.MustCompile(`[^0-9\.]`)
				buildingArea, _ := strconv.ParseFloat(rep.ReplaceAllString(buildingAreaString, ""), 32)
				p.BuildingArea = float32(buildingArea)
				return
			}
			if strings.Contains(s.Text(), "容積率") {
				p.CoverageRatio = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
				return
			}
			if strings.Contains(s.Text(), "用途地域") {
				p.AreaPurpose = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "地目") {
				p.LandKind = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "接道状況") {
				p.Road = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
				return
			}
			if strings.Contains(s.Text(), "引き渡し時期") {
				p.Timing = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "権利") {
				p.Rights = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "リフォーム") {
				p.Reform = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(s.Text(), "備考") {
				p.OtherNotice = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
				return
			}
		})
	})
	
	return p, nil
}

func FetchPropertyFromSearchPageHouseDo(url string, phpSessionId string) []Property {
	print(url)
	ps := []Property{}
	for {
		req, _ := http.NewRequest("GET", url, nil)
		req.AddCookie(&http.Cookie{
			Name: "PHPSESSID", Value: phpSessionId, Path: "/",
		})

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
      time.Sleep(time.Second * 1)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		propertyTables := doc.Find("#exp > table")
		propertyTables.Each(func(i int, s *goquery.Selection) {
			propertyUrl := s.Find("tbody > tr:nth-child(1) > th > a").AttrOr("href", "")
			if propertyUrl == "" {
				return
			}
			print(houseDoBaseUrl + propertyUrl)
			property, err := FetchPropertyFromUrlHouseDo(houseDoBaseUrl + propertyUrl)
			if err != nil {
				return
			}
      fmt.Printf("p: %#v\n", property)
      ps = append(ps, property)
      time.Sleep(time.Second * 10)
		})
    
    nextUrl := doc.Find("#estateList > div.result-peager > ul > li.next > a").First().AttrOr("href", "")
    if nextUrl == "" {
      break
    }
    print("next page", nextUrl)
    url = nextUrl
	}
  return ps
}
