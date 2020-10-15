package go_marketcap_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
)

var idMap []*Currency

type Currency struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Stats
}

type Stats struct {
	Price     float64 `json:"price"`
	Volume24  int64   `json:"volume_24h"`
	MarketCap float64 `json:"market_cap"`

	PercentageChange1H  float64 `json:"percent_change_1h"`
	PercentageChange24H float64 `json:"percent_change_24h"`
	PercentageChange7D  float64 `json:"percent_change_7d"`
}

func init() {
	getIdMapping()
}

func GetCryptoPrice(token string) (*Currency, error) {
	currency, err := getCurrencyBySymbol(token)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	apikey := os.Getenv("CRYPTO_API_KEY")
	if apikey == "" {
		return nil, fmt.Errorf("get me a API_KEY")
	}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	q := url.Values{}
	q.Add("id", fmt.Sprintf("%v", currency.Id))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apikey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	res := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	// todo refactor later, not sure how to parse nested json in go => https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?id=2
	data := res["data"].(map[string]interface{})
	info := data[fmt.Sprintf("%v", currency.Id)].(map[string]interface{})
	quote := info["quote"].(map[string]interface{})
	usd := quote["USD"].(map[string]interface{})

	currency.Price = usd["price"].(float64)
	currency.MarketCap = usd["market_cap"].(float64)
	currency.PercentageChange1H = usd["percent_change_1h"].(float64)
	currency.PercentageChange7D = usd["percent_change_7d"].(float64)
	currency.PercentageChange24H = usd["percent_change_24h"].(float64)

	return currency, nil
}

func getCurrencyBySymbol(symbol string) (*Currency, error) {
	// todo flatten data to just array like =>  ["symbol"] = id;
	if idMap == nil {
		getIdMapping()
	}

	for _, v := range idMap {
		if v.Symbol == strings.ToUpper(symbol) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("crypto not found")
}

func getIdMapping() {
	_, f, _, _ := runtime.Caller(0)
	lastGermanComma := strings.LastIndex(f, "/")
	dir := f[0:lastGermanComma]

	jsonFile, err := os.Open(dir + "/coinIdMap.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		err := jsonFile.Close()
		log.Print(err)
	}()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(byteValue, &idMap)
}