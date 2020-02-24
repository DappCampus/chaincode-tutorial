package model

// Approval is the definition of Approval Event & Data format
type Approval struct {
	Owner     string `json:"owner"`
	Spender   string `json:"spender"`
	Allowance int    `json:"allowance"`
}

func NewApproval(owner, spender string, allowance int) *Approval {
	return &Approval{
		Owner:     owner,
		Spender:   spender,
		Allowance: allowance,
	}
}
