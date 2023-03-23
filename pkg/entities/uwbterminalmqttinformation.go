// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBTerminalMQTTInformation, err := UnmarshalUWBTerminalMQTTInformation(bytes)
//    bytes, err = uWBTerminalMQTTInformation.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalUWBTerminalMQTTInformation(data []byte) (UWBTerminalMQTTInformation, error) {
	var r UWBTerminalMQTTInformation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBTerminalMQTTInformation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UWBTerminalMQTTInformation
type UWBTerminalMQTTInformation struct {
	Alt             float32       `json:"alt"`
	DevEui          string        `json:"devEui"`
	Distance        float32       `json:"distance"`
	Lat             float32       `json:"lat"`
	Lng             float32       `json:"lng"`
	RegionID        int           `json:"regionId"`
	StationInfos    []StationInfo `json:"stationInfos"`
	X               float32       `json:"x"`
	Y               float32       `json:"y"`
	Z               float32       `json:"z"`
	Properties      Properties    `json:"properties"`
	SumDistance     float32       `json:"sumDistance"`     // 这个定义用在计算，原始的 UWB Terminal MQTT Server 中不存在。
	InCacheDateTime time.Time     `json:"inCacheDatetime"` // 这个定义用于计算最后一次获取时间是否超过告警时间，原始 UWB Terminal MQTT Server 中不存在
}

type StationInfo struct {
	DevEui   string    `json:"devEui"`
	Distance float32   `json:"distance"`
	Point    []float32 `json:"point"`
}

type Properties struct {
	SwimmerDisplayName string `json:"swimmerDisplayName"`
	SwimmerGender      string `json:"swimmerGender"`
	SwimmerId          string `json:"swimmerId"`
	Test               string `json:"test"` // 为了验证功能。
}
