package chain_scanner

type Replay struct {
	ChainId int64
	Start   uint64
	End     *uint64
}
