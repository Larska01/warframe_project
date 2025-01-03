package main

import (
	"GOlang_projekti/translator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

// Data types for parsing
type ParsedTime struct {
	State       string `json:"state"`
	ShortString string `json:"shortString"`
}

type Cycles struct {
	Cycle []ParsedTime
}

type ArchonMission struct {
	NodeKey string `json:"nodeKey"`
	TypeKey string `json:"typeKey"`
}

type ParsedArchon struct {
	Boss     string          `json:"boss"`
	Eta      string          `json:"eta"`
	Missions []ArchonMission `json:"missions"`
}

func fetchDayNight() (ParsedTime, error) {
	response, err := http.Get("https://api.warframestat.us/pc/cetusCycle?language=en")
	if err != nil {
		return ParsedTime{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ParsedTime{}, err
	}

	var cycle ParsedTime
	if err := json.Unmarshal(responseData, &cycle); err != nil {
		return ParsedTime{}, err
	}
	return cycle, nil
}

func fetchAlerts() ([]map[string]interface{}, error) {
	response, err := http.Get("https://api.warframestat.us/pc/alerts?language=en")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var alerts []map[string]interface{}
	if err := json.Unmarshal(responseData, &alerts); err != nil {
		return nil, err
	}
	return alerts, nil
}

func fetchArchon() (ParsedArchon, error) {
	response, err := http.Get("https://api.warframestat.us/pc/archonHunt?language=en")
	if err != nil {
		return ParsedArchon{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ParsedArchon{}, err
	}
	var archonHunt ParsedArchon
	if err := json.Unmarshal(responseData, &archonHunt); err != nil {
		return ParsedArchon{}, err
	}
	return archonHunt, nil
}

func cycleHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchDayNight()
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func alertsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchAlerts()
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func archonHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchArchon()
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Url: https://api.warframestat.us/profile/{username}/stats/
type WarframeStat struct {
	UniqueName  string  `json:"uniqueName"`
	Xp          int     `json:"xp"`
	EquipTime   float64 `json:"equiptime"`
	Assists     int     `json:"assists"`
	Kills       int     `json:"kills"`
	CleanedName string  `json:"cleanedName"`
}

type ApiResponse struct {
	Weapons []WarframeStat `json:"weapons"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Query string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("https://api.warframestat.us/profile/%s/stats/", requestBody.Query)
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "API returned an error", http.StatusInternalServerError)
		return
	}

	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
		return
	}

	var filteredWarframes []WarframeStat
	for _, item := range apiResponse.Weapons {
		if strings.HasPrefix(item.UniqueName, "/Lotus/Powersuits/") && !strings.Contains(item.UniqueName, "/Operator/") {
			filteredWarframes = append(filteredWarframes, item)
		}
	}

	sort.Slice(filteredWarframes, func(i, j int) bool {
		return filteredWarframes[i].EquipTime > filteredWarframes[j].EquipTime
	})

	top5 := filteredWarframes
	if len(filteredWarframes) > 5 {
		top5 = filteredWarframes[:5]
	}

	// Apply name translator
	for i := range top5 {
		top5[i].CleanedName = translator.Translate(top5[i].UniqueName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top5)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/alerts", alertsHandler)
	http.HandleFunc("/cycles", cycleHandler)
	http.HandleFunc("/archon", archonHandler)
	http.HandleFunc("/search", searchHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
