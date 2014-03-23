package main

import (
	"drivers"
	"elevator"
	. "fmt"
	"misc"
	"networking"
	"ordercontroller"
	"runtime"
	"time"
)

func main() {
	var mystatus networking.Status
	var takeorders []networking.Order

	runtime.GOMAXPROCS(runtime.NumCPU())

	myip := misc.GetLocalIP()
	Println(myip)
	mystatus.Source = myip

	conf := misc.LoadConfig("./config/conf.json")

	generatedmessages_c := make(chan networking.Networkmessage, 100)
	go networking.TCPPeerToPeer(conf, myip, generatedmessages_c)

	state := "INIT"
	drivers.IoInit()
	elevator.ElevInit()

	for {
		time.Sleep(10 * time.Millisecond)
		mystatus.State = state
		elevator.FloorUpdater()
		mystatus.LastFloor = elevator.CurrentFloor()
		networking.NewStatus(mystatus, generatedmessages_c)
		switch state {
		case "INIT":
			{
				elevator.ElevSetSpeed(-300)
			}
		case "IDLE":
			{
				elevator.ElevSetSpeed(0)
			}
		case "UP":
			{
				elevator.ElevSetSpeed(300)
			}
		case "DOWN":
			{
				elevator.ElevSetSpeed(-300)
			}
		case "DOOR_OPEN":
			{
				elevator.ElevSetDoorOpenLamp(1)
				for _, order := range takeorders {
					order.InOut = 0
					Println("Deleting order: ", order)
					time.Sleep(10 * time.Millisecond)
					networking.Neworder(generatedmessages_c, order)
				}
				elevator.ElevSetSpeed(0)
				time.Sleep(3000 * time.Millisecond)
				elevator.ElevSetDoorOpenLamp(0)
			}
		case "ERROR":
			{
				elevator.ElevSetSpeed(0)
			}
		}
		state, takeorders = ordercontroller.Nextstate(myip, conf.Elevators, mystatus.State)
	}
}
