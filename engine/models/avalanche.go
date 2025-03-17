package models

type Avalanche struct {
	Window   [][]int   `json:"window"`
	PayItems []PayItem `json:"pay_items"`
}
