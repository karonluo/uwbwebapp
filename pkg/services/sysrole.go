package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateSysRole(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	bytesRole, _ := ctx.GetBody()
	role, err := entities.UnmarshalSysRole(bytesRole)

	if err != nil {
		tools.ProcessError(`services.WSCreateSysRole`, `role, err := entities.UnmarshalSysRole(bytesRole)`, err, `pkg/services/sysrole.go`)
		message.StatusCode = 500
		message.Message = err.Error()

	} else {
		id, err := biz.CreateSysRole(role)
		if err != nil {
			message.StatusCode = 500
			message.Message = err.Error()
		} else {
			message.Message = id
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
