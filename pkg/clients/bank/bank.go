package bank

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	DailyRates = "http://www.cbr.ru/scripts/XML_daily.asp?date_req="
	GetMethod  = "GET"
	USD        = "Доллар США"
	EURO       = "Евро"
)

var (
	currentTime = time.Now().Format("02/01/2006")
	url         = DailyRates + currentTime
)

func GetDailyRates() map[string]float64 {
	// send request for daily rates of usd and euro
	req, err := http.NewRequest(GetMethod, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	req.Header.Add("Cookie", "__ddg1_=bMz7QAI3fDT4y8GS26rJ")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "BatPhone/7.26.8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// decode response to map
	data := string(body)
	valCurs := new(CBRValCurs)
	r := bytes.NewReader([]byte(data))
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReaderLabel
	err = d.Decode(&valCurs)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return addToMyValutes(valCurs.Val, USD, EURO)
}

func addToMyValutes(vals []Valute, names ...string) map[string]float64 {
	myValutes := make(map[string]float64)
	// TODO - replace linear search to binary search
	for _, name := range names {
		for _, val := range vals {
			if val.Name == name {
				f := val.Value
				f = strings.Replace(f, ",", ".", -1)
				s, err := strconv.ParseFloat(f, 64)
				if err != nil {
					fmt.Println(err.Error())
					return nil
				}
				myValutes[name] = s
				break
			}
		}
	}
	return myValutes
}

type CBRValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Val     []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
