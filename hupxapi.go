package hupxapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://hupx.hu/en/dam/mc-results/homepage.html?date="
)

func GetHUPX(t string) (hupx map[string]float32) {
	// download the target HTML document
	//currentTime := time.Now()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", baseURL, t), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "sr,en-US;q=0.7,en;q=0.3")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://hupx.hu/en/")
	req.Header.Set("Cookie", "cookie_law_dismissed=1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//	bodyText, err := io.ReadAll(resp.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("%s\n", bodyText)
	//	if resp.StatusCode != 200 {
	//		log.Fatalf("HTTP Error %d: %s", resp.StatusCode, resp.Status)
	//	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println(doc.Text())
	// Pronalazi sve elemente u poslednjoj koloni i izvlači tekstualni sadržaj
	var price []float32
	doc.Find("tr td:nth-last-child(2)").Each(func(i int, s *goquery.Selection) {
		// Izdvajanje samo brojeva iz teksta
		text := s.Text()
		values := strings.Fields(text)
		if len(values) > 0 {
			value := strings.Replace(values[0], ",", "", -1)
			val, _ := strconv.ParseFloat(strings.TrimSpace(value), 32)
			fmt.Printf("Red %d: %f\n", i+1, val)
			price = append(price, float32(val))
		} else {
			fmt.Printf("Red %d: No data found\n", i+1)
		}
	})
	fmt.Printf("Red \n", price[0])
	//hupx := make(map[string]float32)
	hupx["Baseload price"] = price[0]
	hupx["Peakload price"] = price[1]
	hupx["Volume"] = price[2]
	fmt.Println("Hupx \n", hupx)
	return hupx
}
