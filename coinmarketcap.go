package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	baseUrl  = "https://api.coinmarketcap.com/v1"
	graphUrl = "https://graphs.coinmarketcap.com/currencies"
	url      string
)

// GetMarketData - Get information about the global market data of the cryptocurrencies.
func GetMarketData() (GlobalMarketData, error) {
	url = fmt.Sprintf(baseUrl + "/global/")

	resp, err := makeReq(url)

	var data GlobalMarketData
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}

	return data, nil
}

// GetCoinData - Get information about a crypto currency.
func GetCoinData(coin string) (Coin, error) {
	url = fmt.Sprintf("%s/ticker/%s", baseUrl, coin)
	resp, err := makeReq(url)
	if err != nil {
		log.Println(err)
		return Coin{}, err
	}
	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
		return Coin{}, err
	}

	return data[0], nil
}

// GetAllCoinData - Get information about all coins listed in Coin Market Cap.
func GetAllCoinData(limit int) (map[string]Coin, error) {
	var l string
	if limit > 0 {
		l = fmt.Sprintf("?limit=%v", limit)
	}
	url = fmt.Sprintf("%s/ticker/%s", baseUrl, l)

	resp, err := makeReq(url)

	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
	}
	//creating map from the array
	allCoins := make(map[string]Coin)
	for i := 0; i < len(data); i++ {
		allCoins[data[i].ID] = data[i]
	}

	return allCoins, nil
}

// GetCoinGraph - Get graph data points for a crypto currency
func GetCoinGraphData(coin string, start int64, end int64) (CoinGraph, error) {
	url = fmt.Sprintf("%s/%s/%d/%d", graphUrl, coin, start*1000, end*1000)
	resp, err := makeReq(url)
	if err != nil {
		log.Println(err)
		return CoinGraph{}, err
	}
	var data CoinGraph
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.Println(err)
		return CoinGraph{}, err
	}

	return data, nil
}

// doReq - HTTP Client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// makeReq - HTTP Request Helper
func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := doReq(req)
	if err != nil {
		log.Println(err)
	}

	return resp, err
}