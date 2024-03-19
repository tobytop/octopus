package regcenter

type EtcdCenter struct {
}

func NewCenter(address string) *EtcdCenter {
	return &EtcdCenter{}
}
