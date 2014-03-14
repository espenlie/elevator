package networking

import "elevator"

type Order struct {
    Direction   elevator.Elev_button
    Floor       int
    InOut       int
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
