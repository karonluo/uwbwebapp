package tools

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const GOOGLE_DATETIME_FORMAT = "2006-01-02 15:04:05.000000000"
const GOOGLE_DATETIME_FORMAT_NO_NANO = "2006-01-02 15:04:05"
const GOOGLE_DATETIME_FORMAT_NO_TIME = "2006-01-02"
const GOOGLE_DATETIME_FORMAT_NO_DATE_NO_NANO = "15:04:05"

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

/*
统一错误日志处理函数

<moduleFuncMsg> 模块功能, example: services.WSGetLoginInformation

<codeMsg> 代码段, example: res, err := biz.GetLoginInformation(ctx.GetHeader("Authorization"), ctx.FormValue("field"))

<err> 原始错误对象结构

[codeFilePath] 代码段所在代码文件路径, example: /pkg/services/sysuser.go
*/
func ProcessError(moduleFuncMsg string, codeMsg string, err error, codeFilePath ...string) {
	moduleFuncMsgSlice := strings.Split(moduleFuncMsg, `.`)
	var moduleName, funcName string
	funcName = ""
	if len(moduleFuncMsgSlice) > 1 {
		moduleName = moduleFuncMsgSlice[0]
		funcName = moduleFuncMsgSlice[1]
	} else {
		moduleName = moduleFuncMsg
	}
	fmt.Printf("func: %s\t", moduleFuncMsg)
	fmt.Printf("code: %s\t", codeMsg)
	if err != nil {
		fmt.Printf("error: %s\r\n", err.Error())
	}
	error_info := make(map[string]interface{})
	error_info["func"] = moduleFuncMsg
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
	// 目前写入 Redis 内存数据库 由 另外的进程进行持久化存储。
	var log entities.SystemLog
	log.CodeRowContent = codeMsg
	log.FunctionName = funcName
	log.ModuleName = moduleName
	log.Message = err.Error()
	log.Source = "server"
	log.UserDisplayName = "admin"
	log.UserName = "admin"
	log.LogType = "error"
	WriteSystemLog(&log)
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

// 通过 golang gorm 对 type field 的标注获取数据库表字段长度。
func GetDatabaseTableFieldSize(entity interface{}, fieldName string) (int, error) {
	fieldType, _ := reflect.TypeOf(entity).FieldByName(fieldName)
	var fieldSize int
	var err error
	value, found := fieldType.Tag.Lookup(`gorm`)
	if found {
		tmp := strings.Split(value, "size:")
		if len(tmp) > 0 {
			fieldSize, err = strconv.Atoi(tmp[1])
		}
	}
	return fieldSize, err
}

func WriteSystemLog(log *entities.SystemLog) error {
	log.Id = uuid.New().String()
	log.Datetime = time.Now()
	if log.UserName == "" {
		log.UserName = "admin"
	}
	if log.UserDisplayName == "" {
		log.UserDisplayName = "admin"
	}
	err := cache.WriteSystemLogToRedis(log)
	return err
}

// 分解日期时间为全天时间段
func SplitDateTimeToAllDay(date time.Time) entities.BetweenDatetime {
	var result entities.BetweenDatetime
	strBegin := fmt.Sprintf("%s 00:00:00.000000000", date.Format("2006-01-02"))
	strEnd := fmt.Sprintf("%s 23:59:59.999999999", date.Format("2006-01-02"))
	result.BeginDatetime, _ = time.Parse(GOOGLE_DATETIME_FORMAT, strBegin)
	result.EndDatetime, _ = time.Parse(GOOGLE_DATETIME_FORMAT, strEnd)
	return result
}

/*
判断日期是否为空
@param date 用于判断的日期时间
@return bool 是否为空
*/
func CheckNoDate(date time.Time) bool {

	return date.Format(GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01"
}

/*
判断两个时间段是否存在交集
@param dynaStartTime 开始时间日期
@param dynaEndTime 结束时间日期
@param fixedStartTime 用于验证的开始日期时间
@param fixedEndTime 用于验证的结束日期时间
@return bool 返回是否有交集
*/
func CheckTimesHasOverlap(dynaStartTime time.Time, dynaEndTime time.Time, fixedStartTime time.Time, fixedEndTime time.Time) (bool, error) {
	result := false
	var err error
	if dynaStartTime.Format(GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01" ||
		dynaEndTime.Format(GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01" ||
		fixedStartTime.Format(GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01" ||
		fixedEndTime.Format(GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01" {
		err = fmt.Errorf("给定的日期时间参数无效")
		result = false
	} else {
		if dynaEndTime.Before(fixedStartTime) || dynaStartTime.After(fixedEndTime) {
			// 结束日期时间在指定开始日期之前或者开始日期时间在指定结束日期时间之后肯定没有交集
			result = false
		} else {
			result = true
		}
	}
	return result, err

}

// 移除指定元素
func DeleteStringSlice(elems []string, elem string) []string {
	j := 0
	for _, v := range elems {
		if v != elem {
			elems[j] = v
			j++
		}
	}
	return elems[:j]
}

func CheckPointInPolygon(point entities.Point, area []entities.Point) bool {
	// 目标点的x, y坐标
	x := point.X
	y := point.Y

	// 多边形的点数
	count := len(area)

	// 点是否在多边形中
	var inInside bool

	// 浮点类型计算与0的容差
	precision := 2e-10

	// 依次计算每条边，根据每边两端点和目标点的状态栏判断
	for i, j := 0, count-1; i < count; j, i = i, i+1 {
		// 记录每条边上的两个点坐标
		x1 := area[i].X
		y1 := area[i].Y
		x2 := area[j].X
		y2 := area[j].Y

		// 判断点与多边形顶点是否重合
		if (x1 == x && y1 == y) || (x2 == x && y2 == y) {
			return true
		}

		// 判断点是否在水平直线上
		if (y == y1) && (y == y2) {
			return true
		}

		// 判断线段两端点是否在射线两侧
		if (y > y1 && y < y2) || (y < y1 && y > y2) {
			// 斜率
			k := (x2 - x1) / (y2 - y1)

			// 相交点的 x 坐标
			_x := x1 + k*(y-y1)

			// 点在多边形的边上
			if _x == x {
				return true
			}

			// 浮点类型计算容差
			if math.Abs(_x-x) < precision {
				return true
			}

			// 射线穿过多边形的边
			if _x > x {
				inInside = !inInside
			}
		}
	}

	return inInside
}
