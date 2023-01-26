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
	niftyBaseUrl = "https://myhome.nifty.com/"
)

func FetchPropertyFromUrlNifty(url string) (Property, error) {
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

	pankuzus := doc.Find("#wrapper > section.section.is-bg-pj2.is-border-bottom.is-padding-top-xs.is-padding-bottom-xs.is-margin-bottom-xs > div > div > div.column.is-mobile-0 > nav > ul > li")
	p.City = strings.TrimSpace(pankuzus.Eq(4).Text()) + strings.TrimSpace(pankuzus.Eq(5).Text())

	infoDls := doc.Find("#main > section.unit.unitDetail > div > div.detailTop.clearfix > div.detailBox > dl.detailList.strikeStyle.clearfix")
	infoDls.Each(func(i int, dl *goquery.Selection) {
		dtText := dl.Find("dt").Eq(0).Text()
		ddText := dl.Find("dd").Eq(0).Text()
		if strings.Contains(dtText, "価格") {
			rep := regexp.MustCompile(`[^0-9]`)
			price, _ := strconv.ParseUint(rep.ReplaceAllString(ddText, ""), 10, 32)
			if (price != 0) {
				p.Price = uint32(price)
				return
			}
		}
	})

	infoUls := doc.Find("#main > section.unit.unitDetail > div > div.detailTop.clearfix > div.detailBox > ul")
	infoUls.Each(func(i int, ul *goquery.Selection) {
		lis := ul.Find("li")
		lis.Each(func(i int, li *goquery.Selection) {
			dls := li.Find("dl")
			dls.Each(func(i int, dl *goquery.Selection) {
				dts := dl.Find("dt")
				dds := dl.Find("dd")
				dts.Each(func(i int, dt *goquery.Selection) {
					if strings.Contains(dt.Text(), "間取り") {
						p.Layout = strings.TrimSpace(dds.Eq(i).Text())
						return
					}
					if strings.Contains(dt.Text(), "建物面積") {
						buildingAreaString := strings.Split(dds.Eq(i).Text(), "㎡")[0]
						rep := regexp.MustCompile(`[^0-9\.]`)
						buildingArea, _ := strconv.ParseFloat(rep.ReplaceAllString(buildingAreaString, ""), 32)
						p.BuildingArea = float32(buildingArea)
						return
					}
					if strings.Contains(dt.Text(), "土地面積") {
						landAreaString := strings.Split(dds.Eq(i).Text(), "㎡")[0]
						rep := regexp.MustCompile(`[^0-9\.]`)
						landArea, _ := strconv.ParseFloat(rep.ReplaceAllString(landAreaString, ""), 32)
						p.LandArea = float32(landArea)
						return
					}
					if strings.Contains(dt.Text(), "築年月") {
						buildYear, _ := strconv.ParseUint(strings.Split(strings.TrimSpace(dds.Eq(i).Text()), "年")[0], 10, 16)
						p.BuildYear = uint16(buildYear)
						return
					}
				})
			})
		})
	})

	detailTableTrs := doc.Find("#detailInfoTable > tbody > tr")
	detailTableTrs.Each(func(i int, tr *goquery.Selection) {
		ths := tr.Find("th")
		tds := tr.Find("td")
		ths.Each(func(i int, th *goquery.Selection) {
			if strings.Contains(th.Text(), "交通機関") {
				p.Access = strings.TrimSpace(tds.Eq(i).Text())
				rep := regexp.MustCompile(`「([^」]*)」`)
				subMatch := rep.FindStringSubmatch(tds.Eq(i).Text())
				if len(subMatch) > 1 {
					p.Station = strings.TrimSpace(rep.FindStringSubmatch(tds.Eq(i).Text())[1])
					return
				}
				if strings.Contains(tds.Eq(i).Text(), "駅") {
					for _, v := range strings.Fields(tds.Eq(i).Text()) {
						if strings.Contains(v, "駅") {
							p.Station = strings.Replace(v, "駅", "", -1)
							return
						}
					}
				}
			}
			if strings.Contains(th.Text(), "建物構造") {
				p.Structure = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "入居時期") {
				p.Timing = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "地目") {
				p.LandKind = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "用途地域") {
				p.AreaPurpose = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "建ぺい率・容積率") {
				p.CoverageRatio = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
				return
			}
			if strings.Contains(th.Text(), "土地権利") {
				p.Rights = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "私道負担・道路") {
				p.Road = TrimWordGaps(strings.TrimSpace(tds.Eq(i).Text()))
				return
			}
			if strings.Contains(th.Text(), "リフォーム") {
				p.Reform = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "施工会社名") {
				p.BuildCompany = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "その他制限事項") {
				p.OtherRestriction = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
			if strings.Contains(th.Text(), "備考") {
				p.OtherNotice = strings.TrimSpace(tds.Eq(i).Text())
				return
			}
		})
	})

	return p, nil
}

func FetchPropertyFromSearchPageNifty(url string) []Property {
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

		propertyDivs := doc.Find("#searchResultItems > div.buy_list > div")
		propertyDivs.Each(func(i int, div *goquery.Selection) {
			println(i)
			atag := div.Find("div.nayose_head > h2 > p > a")
			if atag.Size() > 0 {
				v := atag.Eq(0).AttrOr("href", "")
				if v == "" {
					return
				}
				println(niftyBaseUrl + v)
				property, err := FetchPropertyFromUrlNifty(niftyBaseUrl + v)
				if err != nil {
					return
				}

				if property.Price == 0 {
					dl := div.Find("div.nayose_head > div > div > div.itemContentWrapper.clearfix > dl").Eq(0)
					dts := dl.Find("dt")
					dds := dl.Find("dd")
					dts.Each(func(i int, dt *goquery.Selection) {
						if strings.Contains(dt.Text(), "価格") {
							rep := regexp.MustCompile(`[^0-9]`)
							price, _ := strconv.ParseUint(rep.ReplaceAllString(dds.Eq(i).Text(), ""), 10, 32)
							if (price != 0) {
								property.Price = uint32(price)
								return
							}
						}
						if strings.Contains(dt.Text(), "所在地") {
							if strings.Contains(dds.Eq(i).Text(), "市") {
								property.City = dds.Eq(i).Text()[:strings.Index(dds.Eq(i).Text(), "市")] + "市"
							} else if strings.Contains(dds.Eq(i).Text(), "町") {
								property.City = dds.Eq(i).Text()[:strings.Index(dds.Eq(i).Text(), "町")] + "町"
							} else if strings.Contains(dds.Eq(i).Text(), "村") {
								property.City = dds.Eq(i).Text()[:strings.Index(dds.Eq(i).Text(), "村")] + "村"
							}
							return
						}
						if strings.Contains(dt.Text(), "交通") {
							property.Access = strings.TrimSpace(dds.Eq(i).Text())
							if strings.Contains(dds.Eq(i).Text(), "/") {
								property.Station = strings.Split(strings.Fields(dds.Eq(i).Text())[0], "/")[1]
							}
							return
						}
						if strings.Contains(dt.Text(), "間取り") {
							property.Layout = strings.TrimSpace(dds.Eq(i).Text())
							return
						}
						if strings.Contains(dt.Text(), "土地面積") {
							landAreaString := strings.Split(dds.Eq(i).Text(), "m")[0]
							rep := regexp.MustCompile(`[^0-9\.]`)
							landArea, _ := strconv.ParseFloat(rep.ReplaceAllString(landAreaString, ""), 32)
							property.LandArea = float32(landArea)
							return
						}
						if strings.Contains(dt.Text(), "建物面積") {
							buildingAreaString := strings.Split(dds.Eq(i).Text(), "m")[0]
							rep := regexp.MustCompile(`[^0-9\.]`)
							buildingArea, _ := strconv.ParseFloat(rep.ReplaceAllString(buildingAreaString, ""), 32)
							property.BuildingArea = float32(buildingArea)
							return
						}
						if strings.Contains(dt.Text(), "築年月") {
							buildYear, _ := strconv.ParseUint(strings.Split(dds.Eq(i).Text(), "年")[0], 10, 16)
							property.BuildYear = uint16(buildYear)
							return
						}
					})
				}

      	fmt.Printf("p: %#v\n", property)
      	ps = append(ps, property)
      	time.Sleep(time.Second * 10)
			}
		})
   
		nextUrl := ""
		pageNations := doc.Find("ul.pageNation > li")
		pageNations.Each(func(i int, li *goquery.Selection) {
			if strings.Contains(li.Text(), "次へ") {
				nextUrl = li.Find("a").Eq(0).AttrOr("href", "")
			}
		})
		if nextUrl == "" {
			break
		}
    println("next page", niftyBaseUrl + nextUrl)
    url = niftyBaseUrl + nextUrl
	}
  return ps
}
