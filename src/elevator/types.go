package elevator

type ElevatorStatus struct {
    Host    string
    State   string
    Floor   int
}

type Order struct {
    Direction   string
    Floor       int
}

type Elevator struct {
    Status  []ElevatorStatus
    Queue   []Order   
}


    
