package entities

import (
	"encoding/json"
	"time"
)

// 系统日志
type SystemLog struct {
	Id              string    `json:"Id"`
	Message         string    `json:"Message"`             // custom message
	LogType         string    `json:"LogType"`             // error, operate, debug, info, warn
	FunctionName    string    `json:"FunctionName"`        // browser func or server func
	ModuleName      string    `json:"ModelName"`           // pkg, entities, biz, dao, etc...
	UserName        string    `json:"UserName"`            // current login user
	UserDisplayName string    `json:"UserDisplayName"`     // current login user display name
	Source          string    `json:"Source"`              // browser or server
	Datetime        time.Time `json:"Datetime"`            // Log produce datetime.
	CodeRowContent  string    `json:"ErrorCodeRowContent"` // code row text
}

func (selfLog *SystemLog) String() string {
	var result string
	byteSystemLog, err := json.Marshal(selfLog)
	if err != nil {
		result = err.Error()
	} else {
		result = string(byteSystemLog)
	}
	return result
}
