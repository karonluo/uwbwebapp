package tools

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"uwbwebapp/pkg/entities"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 计算 SHA1 值
func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

// 将 struct 对象转换成 map
func ReflectMethod(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

func HaveElementInRange(objs []interface{}, obj interface{}) bool {
	var result bool = false
	for _, tmpObj := range objs {
		if reflect.TypeOf(tmpObj) == reflect.TypeOf(obj) && tmpObj == obj {
			result = true
			break
		}
	}
	return result
}

func ProcessError(funcMsg string, codeMsg string, err error, codeFilePath ...string) {
	fmt.Printf("func: %s\t", funcMsg)
	fmt.Printf("code: %s\t", codeMsg)
	if err != nil {
		fmt.Printf("error: %s\r\n", err.Error())
	}
	error_info := make(map[string]interface{})
	error_info["func"] = funcMsg
	error_info["code"] = codeMsg
	error_info["file_path"] = codeFilePath
	if err != nil {
		error_info["error"] = err.Error()
	} else {
		error_info["error"] = ""
	}
	sbytes, _ := json.Marshal(&error_info)
	fmt.Println(string(sbytes))
	// TODO: 写入错误日志库 (建议使用 MongoDB)
}

// 生成查询条件对象
func GenerateQueryConditionFromWebParameters(pageSize string, pageIndex string, likeValue string) (entities.QueryCondition, error) {
	var queryCondition entities.QueryCondition
	var sizeErr, indexErr, err error
	queryCondition.PageIndex, indexErr = strconv.ParseInt(pageIndex, 10, 64)
	if indexErr != nil {
		// ProcessError("tools.GenericQueryConditionFromWebParameters", `companyQueryCondition.PageIndex, err = strconv.ParseInt(pageIndex, 10, 64)`, indexErr)
		err = indexErr
	} else {
		queryCondition.PageSize, sizeErr = strconv.ParseInt(pageSize, 10, 64)
		if sizeErr != nil {
			// ProcessError("tools.GenericQueryConditionFromWebParameters", `companyQueryCondition.PageSize, err = strconv.ParseInt(ctx.FormValue("page_size"), 10, 64)`, sizeErr)
			err = sizeErr
		} else {
			queryCondition.LikeValue = likeValue
		}
	}
	return queryCondition, err
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)

	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
