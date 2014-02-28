package misc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Elevator struct {
	Address		string
}

type Config struct {
	Elevators 			[]Elevator
	Default_Dial_Port	string
	Default_Listen_Port	string
}

var config Config

func LoadConfig(filename string) Config {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Println(err)
	}
	return config
}
