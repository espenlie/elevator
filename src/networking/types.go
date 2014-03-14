package networking

type Order struct {
    Direction   string
    Floor       int
    InOut       bool
}

type Status struct {
    State       string
    LastFloor   int
    Source      string
}

type Networkmessage struct {
    Order   Order
    Status  Status
}
