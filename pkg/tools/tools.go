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
	"time"

	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/paulsmith/gogeos/geos"
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

// HaveElementInRange 接收一个 interface 类型的切片和一个 interface 类型的元素，
// 判断切片中是否存在与给定元素类型相同且值相等的元素。
// 如果存在则返回 true，否则返回 false。
func HaveElementInRange(objs []interface{}, obj interface{}) bool {
	objType := reflect.TypeOf(obj)

	for _, tmpObj := range objs {
		// 检查元素类型和值是否和给定元素相同
		if reflect.TypeOf(tmpObj) == objType && tmpObj == obj {
			return true
		}
	}

	return false
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
	if dynaStartTime.IsZero() || dynaEndTime.IsZero() || fixedStartTime.IsZero() || fixedEndTime.IsZero() {
		return false, fmt.Errorf("给定的日期时间参数无效")
	}
	if dynaEndTime.Before(fixedStartTime) || dynaStartTime.After(fixedEndTime) {
		// 结束日期时间在指定开始日期之前或者开始日期时间在指定结束日期时间之后肯定没有交集
		return false, nil
	}
	return true, nil
}

// 移除指定字符串元素
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

// 移除指定元素
func DeleteSlice[T any](elems []T, elem T, compare func(a, b T) bool) []T {
	var filtered []T
	for _, v := range elems {
		if !compare(v, elem) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// 移除指定元素 linq 方案
func DeleteSliceByLinQ(elems []interface{}, elem interface{}, compare func(a, b interface{}) bool) []interface{} {
	var filtered []interface{}
	query := linq.From(elems).Where(func(v interface{}) bool {
		return !compare(v, elem)
	})
	query.ToSlice(&filtered)
	return filtered
}

func GetGatewayIP(ip string) string {
	ips := strings.Split(ip, ".")
	ip = fmt.Sprintf("%s.%s.%s.1", ips[0], ips[1], ips[2])
	return ip
}

// UseGeosInPolygon 使用 GEOS 库判断给定点是否在多边形内部
// 参数：
//   - point: 需要判断的点
//   - polygon: 多边形的顶点列表
//
// 返回值：
//   - bool: 点是否在多边形内部
//   - error: 若出现错误则返回错误信息，否则返回 nil
func UseGeosInPolygon(point entities.Point, polygon []entities.Point) (bool, error) {
	var result bool = false
	// var coords []geos.Coord
	// 将多边形顶点转换为 GEOS 库的坐标格式
	coords := make([]geos.Coord, 0, len(polygon))
	for _, p := range polygon {
		coords = append(coords, geos.Coord{X: p.X, Y: p.Y})
	}
	// 创建 GEOS 库的线性环
	shell, err := geos.NewLinearRing(coords...)
	if err == nil {
		// 创建 GEOS 库的多边形
		var geosPolygon *geos.Geometry
		geosPolygon, err = geos.PolygonFromGeom(shell)
		if err == nil {
			// 创建 GEOS 库的点
			var geosPoint *geos.Geometry
			geosPoint, err = geos.NewPoint(geos.Coord{X: point.X, Y: point.Y})
			if err == nil {
				// 判断点是否在多边形内部
				result, err = geosPolygon.Contains(geosPoint)
			}
		}
	}
	return result, err
}
