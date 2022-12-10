package biz

import (
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func CreateSysRole(sysrole entities.SysRole) (string, error) {
	sysrole.Id = uuid.New().String()
	if sysrole.Creator == "" {
		sysrole.Creator = "admin"
	}
	if sysrole.Modifier == "" {
		sysrole.Modifier = "admin"

	}
	sysrole.CreateDatetime = time.Now()
	sysrole.ModifyDatetime = time.Now()
	err := dao.CreateSysRole(sysrole)
	return sysrole.Id, err
}
