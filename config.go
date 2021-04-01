package main

import (
	"encoding/json"
	"io/ioutil"
)

func loadConfig() {
	var c Config
	confFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		LogMsg("%+v\n", err)
		return
	}
	err = json.Unmarshal([]byte(confFile), &c)
	if err != nil {
		LogMsg("%+v\n", err)
		return
	}
	config = c

	LogMsg("Setup completed using config file.")
}

func saveConfig() {
	// defer log.PanicSafe()
	file, err := json.Marshal(config)
	if err != nil {
		LogMsg("%+v\n", err)
		// log.LogErrorType(err)
	}
	err = ioutil.WriteFile(configFile, file, 0600)
	if err != nil {
		LogMsg("%+v\n", err)
		// log.LogErrorType(err)
	}
}
