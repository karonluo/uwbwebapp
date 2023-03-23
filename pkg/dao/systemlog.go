package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func WriteSystemLogToDB(log *entities.SystemLog) error {
	return Database.Create(log).Error
}
func GetSystemLogCount(queryCondition entities.QueryCondition) (int64, error) {
	var count int64
	var log entities.SystemLog
	var result *gorm.DB

	if queryCondition.LikeValue != "" {
		whereFields := `message LIKE ? OR function_name LIKE ? OR module_name LIKE ? OR user_display_name LIKE ?`
		result = Database.Model(&log).Where(whereFields,
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%").Count(&count)

	} else {

		result = Database.Model(&log).Count(&count)

	}
	return count, result.Error
}

func QuerySystemLog(queryCondition entities.QueryCondition) ([]entities.SystemLog, error) {
	var log entities.SystemLog
	var logs []entities.SystemLog
	var result *gorm.DB
	selectFields := `id, message, log_type, function_name, module_name, user_name, user_display_name, datetime`
	if queryCondition.LikeValue != "" {
		whereFields := `message LIKE ? OR function_name LIKE ? OR module_name LIKE ? OR user_display_name LIKE ?`
		result = Database.Model(&log).Select(selectFields).Where(whereFields,
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%",
			"%"+queryCondition.LikeValue+"%").
			Order("datetime DESC").
			Limit(int(queryCondition.PageSize)).
			Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
			Find(&logs)

	} else {
		result = Database.Model(&log).Select(selectFields).
			Order("datetime DESC").
			Limit(int(queryCondition.PageSize)).
			Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
			Find(&logs)
	}

	return logs, result.Error
}
func GetSystemLog(id string) (entities.SystemLog, error) {
	var log entities.SystemLog
	err := Database.Where("id=?", id).First(&log).Error
	return log, err
}
