package services

import (
	"encoding/json"
	"fmt"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSWriteOperationLog(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	var log entities.SystemLog
	err := json.Unmarshal(body, &log)
	if err != nil {
		fmt.Printf("Write Operation Log Error: %s\r\n", err.Error())
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		log.LogType = "operate" //限定操作日志
		err = tools.WriteSystemLog(&log)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			fmt.Printf("Write Operation Log Error: %s\r\n", err.Error())
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
