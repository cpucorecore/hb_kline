package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Item struct {
	Id     int64
	Open   float32
	Close  float32
	Low    float32
	High   float32
	Amount float32
	Vol    float32
	Count  uint64
}

type BatchItems struct {
	Ch     string
	Status string
	Ts     uint64
	Data   []Item
}

func TR(H, L, PDC float32) decimal.Decimal {
	return decimal.Max(decimal.NewFromFloat32(H-L), decimal.NewFromFloat32(H-PDC), decimal.NewFromFloat32(PDC-L))
}

func main() {
	request, err := http.NewRequest("GET", "https://api.hadax.com/market/history/kline?period=1day&size=300&symbol=dotusdt", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var d BatchItems
	err = json.Unmarshal(body, &d)
	if err != nil {
		log.Fatal(err)
	}

	base20 := d.Data[len(d.Data)-21:]
	after := d.Data[:len(d.Data)-20]
	sum := 0.0
	for i := 20; i > 0; i-- {
		yestoday := base20[i]
		today := base20[i-1]

		tr := TR(today.High-today.Low, today.High-yestoday.Close, yestoday.Close-today.Low)
		f, _ := tr.Float64()
		sum += f
	}
	PDN := sum / 20.0

	for i := len(after) - 1; i > 0; i-- {
		yestoday := after[i]
		today := after[i-1]

		tr, _ := TR(today.High-today.Low, today.High-yestoday.Close, yestoday.Close-today.Low).Float64()
		PDN = (19.0*PDN + tr) / 20.0
		t := time.Unix(today.Id, 0)
		fmt.Printf("%d.%d.%d,%f,%f,%f,%f,%f\n", t.Year(), t.Month(), t.Day(), today.Low, today.High, today.Open, today.Close, PDN)
	}

}
