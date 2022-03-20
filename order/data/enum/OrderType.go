package enum

type OrderType int

const (
	BuyOrder OrderType = iota
	SellOrder
)

func (t OrderType) String() string {
	return [...]string{"Buy Order", "Sell Order"}[t]
}

func (t OrderType) EnumIndex() int {
	return int(t)
}
