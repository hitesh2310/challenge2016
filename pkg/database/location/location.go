package database

import (
	"encoding/csv"
	"fmt"
	"main/logs"
	"main/pkg/constants"
	"os"
)

var (
	AreaMap map[string]map[string][]string
)

func SetUpLocations() {

	logs.InfoLog("Path of the CSV data:: %v ", constants.ApplicationConfig.Application.DataPath)
	countryProvinceMap := make(map[string]map[string][]string)
	provinceCityMap := make(map[string][]string)

	logs.InfoLog("Parsing over the CSV...")

	file, err := os.Open(constants.ApplicationConfig.Application.DataPath)
	if err != nil {
		fmt.Println("Error to open the CSV file")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read all records from CSV
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records from CSV...")
	}

	fmt.Printf("Number of rows ::%v \n", len(records))

	for index, row := range records {
		if index == 0 {
			fmt.Println("Headers:")
			for hcount, element := range row {
				fmt.Printf("Header %d: %s\n", hcount, element)
			}
			continue
		}

		city, province, country := row[0], row[1], row[2]

		// add cty to province map
		provinceCityMap[province] = append(provinceCityMap[province], city)

		// make map  country if it doesn't exist
		if _, exists := countryProvinceMap[country]; !exists {
			countryProvinceMap[country] = make(map[string][]string)
		}
		// add province map to the country
		countryProvinceMap[country][province] = provinceCityMap[province]
	}
	fmt.Println("MAP plotted...")
	logs.InfoLog("%v", countryProvinceMap)
	AreaMap = countryProvinceMap
}
