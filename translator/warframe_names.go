package translator

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Structs to match the export file
type WarframeAbility struct {
	AbilityName string `json:"abilityName"`
	Description string `json:"description"`
}

type Warframe struct {
	UniqueName string            `json:"uniqueName"`
	Name       string            `json:"name"`
	Abilities  []WarframeAbility `json:"abilities"`
	Passive    string            `json:"passiveDescription"`
}

type WarframeData struct {
	ExportWarframes []Warframe `json:"ExportWarframes"`
}

// Struct to hold the result from the Translate function
type WarframeInfo struct {
	Name       string            `json:"name"`
	UniqueName string            `json:"uniqueName"`
	Abilities  []WarframeAbility `json:"abilities"`
	Passive    string            `json:"passiveDescription"`
}

func loadWarframeData(filePath string) (map[string]Warframe, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var warframeData WarframeData
	err = json.Unmarshal(data, &warframeData)
	if err != nil {
		return nil, err
	}

	warframeMap := make(map[string]Warframe)
	for _, wf := range warframeData.ExportWarframes {
		warframeMap[wf.UniqueName] = wf
	}

	return warframeMap, nil
}

func getWarframeInfo(uniqueName string, warframeMap map[string]Warframe) WarframeInfo {
	if wf, exists := warframeMap[uniqueName]; exists {
		return WarframeInfo{Name: wf.Name, Abilities: wf.Abilities, Passive: wf.Passive}
	}
	return WarframeInfo{Name: "Unknown Warframe", Abilities: nil}
}

// Translates the database name to correct name and searches for warframe ability information
func TranslateAndAddAbilityInfo(rawName string) WarframeInfo {
	warframeMap, err := loadWarframeData("static/WarframeExported/ExportWarframes_en.json")
	if err != nil {
		log.Fatal(err)
	}

	warframeInfo := getWarframeInfo(rawName, warframeMap)

	return warframeInfo
}
