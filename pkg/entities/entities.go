package entities

import "time"

type QueryCondition struct {
	PageSize  int64  `json:"page_size"`
	PageIndex int64  `json:"page_index"`
	LikeValue string `json:"like_value"`
}

type BetweenDatetime struct {
	BeginDatetime time.Time
	EndDatetime   time.Time
}

func (btet *BetweenDatetime) String() (string, string) {

	return btet.BeginDatetime.Format("2006-01-02 15:04:05 00:00:00.000000000"),
		btet.EndDatetime.Format("2006-01-02 15:04:05 00:00:00.000000000")
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AlertInformation struct {
	Message            string  `json:"message"`
	SwimmerId          string  `json:"swimmerId"`
	SwimmerDisplayName string  `json:"swimmerDisplayName"`
	SwimmerGender      string  `json:"swimmerGender"`
	DevEui             string  `json:"devEui"`
	Type               string  `json:"type"` // normal/danger
	X                  float64 `json:"x"`
	Y                  float64 `json:"y"` // 该游泳者相关告警的最后一次出现的坐标
}

type EROK struct {
	ErrList     []string `json:"errList"`
	SuccessList []string `json:"successList"`
}
