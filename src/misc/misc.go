package misc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"networking"
	"fmt"
)

type Elevator struct {
	Address		string
}

type Config struct {
	Elevators			[]Elevator
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

func PackNetworkMessage(message networking.Networkmessage) []byte {
	pack, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Could not pack message: ",err.Error())
	}
	return pack
}

func UnpackNetworkMessage(pack []byte) networking.Networkmessage{
	var message networking.Networkmessage
	err := json.Unmarshal(pack, &message)
	if err != nil {
		fmt.Println("Could not unpack message: ", err.Error())
	}
	return message
}
