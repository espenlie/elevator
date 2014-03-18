package networking

import "elevator"

type Order struct {
    Direction   elevator.Elev_button
    Floor       int
    InOut       int
    Source      string
}


type Status struct {
    State       string
    LastFloor   int
    Inhouse     bool
    Source      string
}

type Networkmessage struct {
    Order   Order
    Status  Status
}

var EmptyOrder = []Order{{ Direction  : elevator.BUTTON_CALL_UP,
                        Floor       : 0,
                        InOut       : 0,
                        Source      : "",
                        }}

