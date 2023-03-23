// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBDevicePlatformTerminalGetResult, err := UnmarshalUWBDevicePlatformTerminalGetResult(bytes)
//    bytes, err = uWBDevicePlatformTerminalGetResult.Marshal()

package entities

import "encoding/json"

func UnmarshalUWBDevicePlatformTerminalGetResult(data []byte) (UWBDevicePlatformTerminalGetResult, error) {
	var r UWBDevicePlatformTerminalGetResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBDevicePlatformTerminalGetResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UWBDevicePlatformTerminalGetResult

type UWBDevicePlatformTerminalGetResult struct {
	Code int64                  `json:"code"`
	Data UWBTerminalInformation `json:"data"`
	Msg  string                 `json:"msg"`
}

type UWBTerminalInformation struct {
	ApplicationID    int64      `json:"applicationId"`
	ApplicationName  string     `json:"applicationName"`
	Battery          int64      `json:"battery"`
	CoordinateSystem int64      `json:"coordinateSystem"`
	Enabled          bool       `json:"enabled"`
	Model            Model      `json:"model"`
	Name             string     `json:"name"`
	Properties       Properties `json:"properties"`
	SerialNumber     string     `json:"serialNumber"`
	State            int64      `json:"state"`
}
