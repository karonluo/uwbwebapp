package entities

type QueryCondition struct {
	PageSize  int64  `json:"page_size"`
	PageIndex int64  `json:"page_index"`
	LikeValue string `json:"like_value"`
}
