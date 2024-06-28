package macro

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	apiKey = "8c02c7db2bc569741ad1a1b96beb73b8" // 替换为你的实际API密钥
	apiUrl = "https://api.stlouisfed.org/fred/series/observations"
)

type Observation struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

type ApiResponse struct {
	Observations []Observation `json:"observations"`
}

func fetchEconomicData(seriesID, startDate, endDate string) (*ApiResponse, error) {
	url := fmt.Sprintf("%s?series_id=%s&api_key=%s&file_type=json&observation_start=%s&observation_end=%s", apiUrl, seriesID, apiKey, startDate, endDate)
	log.Println("fred url: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func RunAndGetData(seriesId string) (content string) {
	//seriesIDs := map[string]string{
	//	//"10年期国债收益率": "DGS10",
	//	//"SOFR利率":    "SOFR",
	//	"联邦基金利率": "FEDFUNDS",
	//}
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")

	data, err := fetchEconomicData(seriesId, startDate, endDate)
	if err != nil {
		fmt.Println(os.Stderr, "error fetching data for %s: %v\n", seriesId, err)
	}

	for _, obs := range data.Observations {
		date, _ := time.Parse("2006-01-02", obs.Date)
		//fmt.Printf("日期: %s, 值: %s\n", date.Format("2006-01-02"), obs.Value)
		content += fmt.Sprintf("日期: %s, 值: %s\n", date.Format("2006-01-02"), obs.Value)
	}

	return content
}
