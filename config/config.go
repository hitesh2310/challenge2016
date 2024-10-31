package config

import (
	"fmt"
	"main/logs"
	"main/pkg/constants"
	database "main/pkg/database/location"

	// database "main/pkg/database/location"
	"main/pkg/models"

	"github.com/spf13/viper"
)

func SetUpApplication() {
	fmt.Println("SetUp config details....")
	setupConfig()
	fmt.Println("SettingUp Logs....")
	setUpApplicationLogs()
	fmt.Println("SettingUp Database...")
	setUpApplicationDatabase()

	// updateCounter()

}

func setupConfig() {
	viper.SetConfigName("config")    // Name of the config file (without extension)
	viper.SetConfigType("json")      // Config file type
	viper.AddConfigPath("./config/") // Path to the directory containing the config file

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	var config models.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error to unmarshal config")
	}

	constants.ApplicationConfig = &config

}

func setUpApplicationLogs() {
	logs.SetUpApplicationLogs()
}

func setUpApplicationDatabase() {
	logs.InfoLog("Loading the geographies.....")
	database.SetUpLocations()
}
