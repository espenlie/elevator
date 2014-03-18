package misc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"os/exec"
	"fmt"
)

type Elevator struct {
	Address		string
}

type Config struct {
	Elevators			[]Elevator
	DefaultDialPort		string
	DefaultListenPort	string
	Timeout				int
	NumFloors			int
	StopReverseTime		int	
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

func GetLocalIP() string {
    oneliner := "ifconfig | grep 129.241.187 | cut -d':' -f2 | cut -d' ' -f1"
    cmd := exec.Command("bash", "-c", oneliner)
    out, err := cmd.Output()
    if err != nil {
        fmt.Println(err)
    }
    ip := strings.TrimSpace(string(out))
    return ip
}

func Abs(i int) int {
	if i < 0 {
		return i*-1
	}
	return i
}
