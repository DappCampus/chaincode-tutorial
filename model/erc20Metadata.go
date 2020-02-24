package model

// ERC20Metadata is the definition of Token Meta Info
type ERC20Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"owner"`
	TotalSupply uint64 `json:"totalSupply"`
}

func NewERC20MetaData(name, symbol, owner string, totalSupply uint64) *ERC20Metadata {
	return &ERC20Metadata{
		Name:        name,
		Symbol:      symbol,
		Owner:       owner,
		TotalSupply: totalSupply,
	}
}

func (erc20 *ERC20Metadata) GetTotalSupply() *uint64 {
	return &erc20.TotalSupply
}
