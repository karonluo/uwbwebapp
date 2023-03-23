// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBDevicePlatformLoginInformation, err := UnmarshalUWBDevicePlatformLoginInformation(bytes)
//    bytes, err = uWBDevicePlatformLoginInformation.Marshal()

package entities

import "encoding/json"

func UnmarshalUWBDevicePlatformLoginInformation(data []byte) (UWBDevicePlatformLoginInformation, error) {
	var r UWBDevicePlatformLoginInformation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBDevicePlatformLoginInformation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// ApifoxModel
//
// UWBDevicePlatformLoginInformation
type UWBDevicePlatformLoginInformation struct {
	Code int64        `json:"code"`
	Data UWBLoginData `json:"data"`
	Msg  string       `json:"msg"`
}

type UWBLoginData struct {
	AccessToken        string     `json:"accessToken"`
	AccessTokenExpire  int64      `json:"accessTokenExpire"`
	RefreshToken       string     `json:"refreshToken"`
	RefreshTokenExpire int64      `json:"refreshTokenExpire"`
	UserInfoVo         UserInfoVo `json:"userInfoVo"`
}

type UserInfoVo struct {
	Admin         bool           `json:"admin"`
	CreatedAt     string         `json:"createdAt"`
	ID            int64          `json:"id"`
	Organizations []Organization `json:"organizations"`
	Phone         string         `json:"phone"`
	RealName      string         `json:"realName"`
	SuperAdmin    bool           `json:"superAdmin"`
}

type Organization struct {
	ID   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}
