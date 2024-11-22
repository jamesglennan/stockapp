package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

const (
	stockEnvVar    = "SYMBOL"
	numDaysEnvVar  = "NDAYS"
	apiTokenEnvVar = "APIKEY"
	baseURL        = "https://www.alphavantage.co/query"
)

func main() {
	stock := os.Getenv(stockEnvVar)
	nDaysStr := os.Getenv(numDaysEnvVar)
	APIToken := os.Getenv(apiTokenEnvVar)

	if len(stock) < 1 {
		log.Fatalf("Env Var %s is set unset", stockEnvVar)
	}

	if len(APIToken) < 1 {
		log.Fatalf("Env Var %s is set unset", apiTokenEnvVar)
	}

	if len(nDaysStr) < 1 {
		log.Fatalf("Env Var %s is set unset", numDaysEnvVar)
	}

	nDays, err := strconv.Atoi(nDaysStr)
	if err != nil {
		log.Fatalf("Unable to convert nDays to int: %v, err: %v", nDaysStr, err)
	}

	data := getData(stock, nDays, APIToken)

	//use the basic go http library, if we want to get serious and enable
	// middlewares, logging, auth, tracing etc probably use a more sophisticated library
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jEncode := json.NewEncoder(w)
		err := jEncode.Encode(data)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding json: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Serving request from %v", r.RemoteAddr)

	})
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}

func buildURL(stock, APIToken string) string {
	params := url.Values{}
	params.Add("symbol", stock)
	params.Add("apikey", APIToken)
	params.Add("function", "TIME_SERIES_DAILY")
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func sendRequest(requestURL string, data *map[string]interface{}) {
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

}

func getClose(date string, timeSeries *map[string]interface{}) float64 {
	dayData, ok := (*timeSeries)[date].(map[string]interface{})

	if !ok {
		log.Fatalf("can't get the given date: %v val: %v", date, (*timeSeries)[date])
	}
	//Get the close price using its key
	closePriceStr, ok := dayData["4. close"].(string)
	if !ok {
		log.Fatalf("Could not extract close price: %v", dayData)
	}

	closePrice, err := strconv.ParseFloat(closePriceStr, 64)
	if err != nil {
		log.Fatalf("Error parsing float: val: %v, err: %v", closePriceStr, err)
	}

	return closePrice
}

// Get data
func getData(stock string, nDays int, APIToken string) map[string]float64 {
	results := map[string]float64{"average": 0.0}
	var jsonData map[string]interface{}

	//Moved these to their own functions for clarity
	requestURL := buildURL(stock, APIToken)
	sendRequest(requestURL, &jsonData)

	//Get Time Series
	timeSeries, ok := jsonData["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		log.Fatalf("Failed to parse time series data")
	}

	//Iterate through map and extract the dates
	//Sort them so its easy to get the last nDays
	//Originally approached this by just taking time.Now().After(0,0,-nDays)
	//But quickly realised weekends and public holidays
	//meant that it was easier to deal with it this way
	var dates []string
	for date := range timeSeries {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	//Iterate through the last nDays to and extract the close price
	for _, date := range dates[len(dates)-nDays:] {
		//Pass by reference incase timeSeries is a huge structure and save memory
		closePrice := getClose(date, &timeSeries)
		results[date] = closePrice
		results["average"] += closePrice
	}

	//Divide the summed values by nDays at the end to set average
	results["average"] /= float64(nDays)
	return results
}

// Serve data
