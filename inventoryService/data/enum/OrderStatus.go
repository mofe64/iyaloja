package enum

type OrderStatus int

const (
	Pending OrderType = iota
	Fulfilled
	Cancelled
)

func (s OrderStatus) String() string {
	return [...]string{"Pending", "Fulfilled", "Cancelled"}[s]
}

func (s OrderStatus) EnumIndex() int {
	return int(s)
}
