package main

import (
	//"fmt"
  //"time"

	"kouya17/scsave"
)

func main() {
  /*
  property, _ := scsave.FetchPropertyFromUrlNifty("https://myhome.nifty.com/chuko/ikkodate/aichi/nagoyashitempakuku/suumof_70313534/")
  fmt.Printf("%v\n", property)
  time.Sleep(time.Second * 1)
  property, _ = scsave.FetchPropertyFromUrlNifty("https://myhome.nifty.com/chuko/ikkodate/aichi/nagoyashimidoriku/suumof_70763724/")
  fmt.Printf("%v\n", property)
  */
  scsave.FetchPropertyFromSearchPageNifty("https://myhome.nifty.com/chuko/ikkodate/tokai/aichi/?cities=toyohashishi,okazakishi,ichinomiyashi,kasugaishi,nishioshi&subtype=buh&isFromSearch=1&b2=15000000&page=1")
}
