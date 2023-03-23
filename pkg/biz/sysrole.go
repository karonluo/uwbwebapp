package biz

import (
	"math"
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

// 将系统角色和系统功能页面进行绑定
func SysRoleJoinInSysFuncPages(sysRoleId string, sysFuncPageIds []string) error {
	return dao.SysRoleJoinInSysFuncPages(sysRoleId, sysFuncPageIds)
}
func ClearSysRoleSysFuncPages(sysRoleId string) error {
	return dao.ClearSysRoleSysFuncPages(sysRoleId)
}

func EnumAllFuncPagesByRoleId(sysRoleId string) ([]entities.SysFuncPage, error) {
	return dao.EnumAllFuncPagesByRoleId(sysRoleId)
}

func GetSysRoleById(sysRoleId string) (entities.SysRole, error) {
	return dao.GetSysRoleById(sysRoleId)
}

func QuerySysRoles(queryCondition entities.QueryCondition) ([]entities.SysRole, int64, int64, error) {
	var roles []entities.SysRole
	dataRecordCount, err := dao.GetSysRoleCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		roles, err = dao.QuerySysRoles(queryCondition)
	}
	return roles, int64(math.Ceil(pageCount)), dataRecordCount, err

}

func UpdateSysRole(role *entities.SysRole) error {
	var tmpRole entities.SysRole
	var err error
	tmpRole, err = dao.GetSysRoleById(role.Id)
	if err == nil {
		// 防止以下字段被修改
		role.CreateDatetime = tmpRole.CreateDatetime
		role.Creator = tmpRole.Creator
		role.ModifyDatetime = time.Now()

		if role.Modifier == "" {
			role.Modifier = "admin"
		}
		err = dao.UpdateSysRole(role)
	}
	return err
}
