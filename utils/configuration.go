package utils

import "flag"

var AppConfigurationFileName = flag.String("config", "app-config.yml", "Location of the config file.")

func LoadConfig(fileName *string, settings interface{}) {
	file := loadFile(*fileName)
	setSettings(settings, file)
}
