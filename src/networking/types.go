package networking

type Order struct {
    Direction   string
    Floor       int
}

type Status struct {
    State       string
    LastFloor   int
    Source      string
}
