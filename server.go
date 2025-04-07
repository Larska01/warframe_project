package main

import (
	"GOlang_projekti/translator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
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
	file, err := os.Open("data/cetus_cycle.json")
	if err != nil {
		return ParsedTime{}, err
	}
	defer file.Close()

	responseData, err := ioutil.ReadAll(file)
	if err != nil {
		return ParsedTime{}, err
	}

	var cetusTime ParsedTime
	if err := json.Unmarshal(responseData, &cetusTime); err != nil {
		return ParsedTime{}, err
	}

	return cetusTime, nil
}

func fetchAlerts() ([]map[string]interface{}, error) {
	file, err := os.Open("data/alerts.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	responseData, err := ioutil.ReadAll(file)
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
	file, err := os.Open("data/archon_hunt.json")
	if err != nil {
		return ParsedArchon{}, err
	}
	defer file.Close()

	responseData, err := ioutil.ReadAll(file)
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
type WarframeAbility struct {
	AbilityName string `json:"abilityName"`
	Description string `json:"description"`
}

type WarframeData struct {
	UniqueName  string            `json:"uniqueName"`
	CleanedName string            `json:"cleanedName"`
	EquipTime   float64           `json:"equiptime"`
	Xp          int               `json:"xp"`
	Abilities   []WarframeAbility `json:"abilities"`
	Passive     string            `json:"passiveDescription"`
}

type ApiResponse struct {
	Weapons []WarframeData `json:"weapons"`
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(apiURL)
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

	var filteredWarframes []WarframeData
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

	// Apply name translator and append ability info
	for i := range top5 {
		warframeInfo := translator.TranslateAndAddAbilityInfo(top5[i].UniqueName)

		abilities := make([]WarframeAbility, len(warframeInfo.Abilities))
		for j, ability := range warframeInfo.Abilities {
			abilities[j] = WarframeAbility{
				AbilityName: ability.AbilityName,
				Description: ability.Description,
			}
		}

		top5[i] = WarframeData{
			UniqueName:  top5[i].UniqueName,
			CleanedName: warframeInfo.Name,
			EquipTime:   top5[i].EquipTime,
			Xp:          top5[i].Xp,
			Abilities:   abilities,
			Passive:     warframeInfo.Passive,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(top5)
}

// TODO: Add news page. URL: https://api.warframestat.us/pc/
type Event struct {
	Date       string `json:"date"`
	Message    string `json:"message"`
	Link       string `json:"link"`
	ImageLink  string `json:"imageLink"`
	Timestring string `json:"asString"`
}

type News struct {
	News []Event `json:"news"`
}

func fetchNews() (News, error) {
	file, err := os.Open("data/news.json")
	if err != nil {
		return News{}, err
	}
	defer file.Close()

	responseData, err := ioutil.ReadAll(file)
	if err != nil {
		return News{}, err
	}

	var news News
	if err := json.Unmarshal(responseData, &news); err != nil {
		return News{}, err
	}

	return news, nil
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := fetchNews()
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/alerts", alertsHandler)
	http.HandleFunc("/cycles", cycleHandler)
	http.HandleFunc("/archon", archonHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/news", newsHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
