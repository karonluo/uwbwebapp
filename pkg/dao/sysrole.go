package dao

import (
	"uwbwebapp/pkg/entities"
)

func CreateSysRole(role entities.SysRole) error {
	result := Database.Create(&role)
	return result.Error
}
