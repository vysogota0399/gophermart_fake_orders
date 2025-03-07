package models

type Product struct {
	Match      string `json:"match"`
	Reward     int64  `json:"reward"`
	RewardType string `json:"reward_type"`
}
