package ordercontroller

import (
	"elevator"
	"misc"
	"networking"
)
//This function does a BFS-search through all orders to find the most effective solution
func Nextorder(myip string, Elevatorlist []misc.Elevator) networking.Order {
	var statelist = make(map[string]networking.Status)
	statuslist := networking.GetStatusList()
	for host, status := range statuslist {
		statelist[host] = status
	}
	orderlist := networking.GetOrderList()
insideloop:
	for _, order := range orderlist {
		if order.Direction != elevator.BUTTON_COMMAND {
			continue insideloop
		}
		for _, elevator := range Elevatorlist {
			if status, ok := statelist[elevator.Address]; ok {
				if ((status.State == "UP" || status.State == "IDLE") && status.LastFloor <= order.Floor) || ((status.State == "DOWN" || status.State == "IDLE") && status.LastFloor >= order.Floor) {
					if status.Source == order.Source {
						if status.Source == myip {
							return order
						} else {
							delete(statelist, elevator.Address)
							continue insideloop
						}
					}
				}
			}
		}
		for _, elevator := range Elevatorlist {
			if status, ok := statelist[elevator.Address]; ok {
				if (status.State == "UP" && status.LastFloor >= order.Floor) || (status.State == "DOWN" && status.LastFloor <= order.Floor){
					if status.Source == order.Source {
						if status.Source == myip {
							return order
						} else {
							delete(statelist, elevator.Address)
							continue insideloop
						}
					}
				}
			}
		}
	}
orderloop:
	for _, order := range orderlist {
		if order.Direction == elevator.BUTTON_COMMAND {
			continue orderloop
		}
		for i := 0; i < elevator.N_FLOORS; i++ {
			for _, elevator := range Elevatorlist {
				if status, ok := statelist[elevator.Address]; ok {
					if i != 0 && (status.State == "UP" && status.LastFloor+i == order.Floor) || (status.State == "DOWN" && status.LastFloor-i == order.Floor) {
						if statelist[elevator.Address].Source == myip {
							return order
						} else {
							delete(statelist, elevator.Address)
							continue orderloop
						}
					}
				}
			}
			for _, elevator := range Elevatorlist {
				if status, ok := statelist[elevator.Address]; ok {
					if status.State == "IDLE" && (status.LastFloor == order.Floor+i || status.LastFloor == order.Floor-i) {
						if statelist[elevator.Address].Source == myip {
							return order
						} else {
							delete(statelist, elevator.Address)
							continue orderloop
						}
					}
				}
			}
		}
	}
	return networking.EmptyOrder[0]
}

//This function return orders the elevator should stop for
func Stop(myip string, mystate string) []networking.Order {
	var takeorder []networking.Order
	orderlist := networking.GetOrderList()
	for _, order := range orderlist {
		if (order.Direction == elevator.BUTTON_COMMAND && order.Source == myip) || (order.Direction == elevator.BUTTON_CALL_UP && mystate == "UP") || (order.Direction == elevator.BUTTON_CALL_DOWN && mystate == "DOWN") {
			if order.Floor == elevator.CurrentFloor() && elevator.ElevAtFloor() {
				takeorder = append(takeorder, order)
			}
		}
	}
	return takeorder
}
//This function returns the next state for the elevator
func Nextstate(myip string, elevators []misc.Elevator, mystate string) (string, []networking.Order) {
	if elevator.ElevGetObstructionSignal() {
		elevator.ElevSetStopLamp(1)
		return "ERROR", nil
	} else if mystate == "ERROR" {
		elevator.ElevSetStopLamp(0)
		return "INIT", nil
	}

	stop := Stop(myip, mystate)
	if len(stop) != 0 {
		return "DOOR_OPEN", stop
	}

	next := Nextorder(myip, elevators)
	if elevator.ElevAtFloor() && next.Floor == elevator.CurrentFloor() {
		return "DOOR_OPEN", append(stop, next)
	}
	if next.Floor > elevator.CurrentFloor() {
		return "UP", nil
	} else if next.Floor < elevator.CurrentFloor() && next.Floor != 0 {
		return "DOWN", nil
	} else if elevator.ElevAtFloor() {
		return "IDLE", nil
	} else {
		return mystate, nil
	}
}
