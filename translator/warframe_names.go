package translator

// This file is used to translate the warframes to their proper names using the warframe public export json files
// TODO: For the top 5 warframes feature maybe more information on them can be viewed using the public export

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Structs to match the export file
type Warframe struct {
	UniqueName string `json:"uniqueName"`
	Name       string `json:"name"`
}

type WarframeData struct {
	ExportWarframes []Warframe `json:"ExportWarframes"`
}

func loadWarframeNames(filePath string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var warframeData WarframeData
	err = json.Unmarshal(data, &warframeData)
	if err != nil {
		return nil, err
	}

	warframeNames := make(map[string]string)
	for _, wf := range warframeData.ExportWarframes {
		warframeNames[wf.UniqueName] = wf.Name
	}

	return warframeNames, nil
}

func getWarframeName(uniqueName string, warframeNames map[string]string) string {
	if name, exists := warframeNames[uniqueName]; exists {
		return name
	}
	return "Unknown Warframe"
}

func Translate(rawName string) string {
	var cleanName string
	warframeNames, err := loadWarframeNames("static/WarframeExported/ExportWarframes_en.json")
	if err != nil {
		log.Fatal(err)
	}

	cleanName = getWarframeName(rawName, warframeNames)

	return cleanName
}
