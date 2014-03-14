package networking

type Order struct {
    Direction   string
    Floor       int
    InOut       bool
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
