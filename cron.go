package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
)

func setupPriceCron(db *sqlx.DB) {
	c := cron.New()
	// cronExp := "*/5 * * * *" // each 5 minutes
	cronExp := "0 6,18 * * *" // each day at 06 and 18 hours
	_, err := c.AddFunc(cronExp, func() {
		fetchUSTDPrices(db)
	})
	panicIfErr(err)
	c.Start()
}

type USTDResp struct {
	Data []struct {
		Adv struct {
			Price string `json:"price"`
		} `json:"adv"`
	} `json:"data"`
}

func fetchUSTDPrices(db *sqlx.DB) {
	url := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	jsonStr := `{
			"fiat": "BOB",
			"page": 1,
			"rows": 20,
			"tradeType": "BUY",
			"asset": "USDT",
			"countries": [],
			"proMerchantAds": false,
			"shieldMerchantAds": false,
			"filterType": "all",
			"periods": [],
			"additionalKycVerifyFilter": 0,
			"publisherType": "merchant",
			"payTypes": [],
			"classifies": ["mass", "profession", "fiat_trade"]
		}`

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		log.Println("ERROR: post request", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: reading body", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("ERROR: not OK request", resp)
		return
	}

	var data USTDResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("ERROR: unmarshaling body", err)
		return
	}
	calculateAndSaveUSTDPrice(data, db)
}

func calculateAndSaveUSTDPrice(data USTDResp, db *sqlx.DB) {
	var totalPrice int64
	for i := range data.Data {
		priceResp := data.Data[i]
		priceParts := strings.Split(priceResp.Adv.Price, ".")
		bs, _ := strconv.ParseInt(priceParts[0], 10, 64)
		cents, _ := strconv.ParseInt(priceParts[1], 10, 64)
		cents += bs * 100
		totalPrice += cents
	}
	var avgPrice int64 = totalPrice / int64(len(data.Data))
	now := time.Now().UTC()

	sql := "INSERT INTO ustd_price (id, price, created_at, updated_at) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(sql, uuid.New(), avgPrice, now, now)
	if err != nil {
		log.Println("ERROR: inserting to db", err)
	}
}
