// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBDevicePlatformTerminalSearchResult, err := UnmarshalUWBDevicePlatformTerminalSearchResult(bytes)
//    bytes, err = uWBDevicePlatformTerminalSearchResult.Marshal()

package entities

import "encoding/json"

func UnmarshalUWBDevicePlatformTerminalSearchResult(data []byte) (UWBDevicePlatformTerminalSearchResult, error) {
	var r UWBDevicePlatformTerminalSearchResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBDevicePlatformTerminalSearchResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UWBDevicePlatformTerminalSearchResult
type UWBDevicePlatformTerminalSearchResult struct {
	Code int64 `json:"code"`
	Data struct {
		CurrentPage   int64   `json:"currentPage"`
		Data          []Datum `json:"data"`
		Empty         bool    `json:"empty"`
		HasNext       bool    `json:"hasNext"`
		HasPre        bool    `json:"hasPre"`
		PageSize      int64   `json:"pageSize"`
		TotalElements int64   `json:"totalElements"`
		TotalPages    int64   `json:"totalPages"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Datum struct {
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

type Model struct {
	Description string `json:"description"`
	ID          int64  `json:"id"`
	ModelType   int64  `json:"modelType"`
	Name        string `json:"name"`
}
