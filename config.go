package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func loadConfig() {
	var c Config
	confFile, _ := ioutil.ReadFile(configFile)
	err := json.Unmarshal([]byte(confFile), &c)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	config = c

	fmt.Println("Setup completed using config file.")
}

func saveConfig() {
	//defer log.PanicSafe()
	file, err := json.Marshal(config)
	if err != nil {
		//log.LogErrorType(err)
	}
	err = ioutil.WriteFile(configFile, file, 0600)
	if err != nil {
		//log.LogErrorType(err)
	}
}
