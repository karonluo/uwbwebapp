package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func CreateUWBTag(tag *entities.UWBTag) error {
	return Database.Create(tag).Error
}

func GetUWBTag(code string) (entities.UWBTag, error) {
	var tag entities.UWBTag
	res := Database.Where("code =?", code).Find(&tag)
	return tag, res.Error
}

func GetUWBTagCount(queryCondition entities.QueryCondition, companyIds []string, isBound interface{}) (int64, error) {
	var count int64
	var tag entities.UWBTag
	var result *gorm.DB

	if queryCondition.LikeValue != "" {
		whereFields := ` (code LIKE ? OR description LIKE ?) `
		if len(companyIds) > 0 {
			if isBound != nil {
				result = Database.Model(&tag).Where(whereFields+` AND sports_company_id IN ? AND is_bound = ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					companyIds,
					isBound.(bool)).Count(&count)
			} else {
				result = Database.Model(&tag).Where(whereFields+` AND sports_company_id IN ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					companyIds).Count(&count)
			}
		} else {
			if isBound != nil {
				result = Database.Model(&tag).Where(whereFields+` AND is_bound = ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					isBound.(bool)).
					Count(&count)
			} else {
				result = Database.Model(&tag).Where(whereFields,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%").
					Count(&count)
			}
		}
	} else {
		whereFields := ` 1=1 `
		if len(companyIds) > 0 {
			if isBound != nil {
				result = Database.Model(&tag).Where(whereFields+` AND sports_company_id IN ? AND is_bound = ?`, companyIds, isBound.(bool)).Count(&count)
			} else {
				result = Database.Model(&tag).Where(whereFields+` AND sports_company_id IN ?`, companyIds).Count(&count)
			}
		} else {
			if isBound != nil {
				result = Database.Model(&tag).Where(whereFields+` AND is_bound = ?`, isBound.(bool)).Count(&count)
			} else {
				result = Database.Model(&tag).Count(&count)
			}
		}
	}
	return count, result.Error
}

// 查询检索UWB标签信息
// queryCondition 查询条件
// companyIds 要检索的公司集合
// 是否已经绑定了用户（false 否/true 真/nil 全部）
func QueryUWBTags(queryCondition entities.QueryCondition, companyIds []string, isBound interface{}) ([]entities.UWBTag, error) {
	var tag entities.UWBTag
	var tags []entities.UWBTag
	var result *gorm.DB
	selectFields := `code, sports_company_id, sports_company_name, is_bound, description,create_datetime,modify_datetime,Creator,Modifier`
	if queryCondition.LikeValue != "" {
		whereFields := `(code LIKE ? OR description LIKE ?) `
		if len(companyIds) > 0 {
			if isBound != nil {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND sports_company_id IN ? AND is_bound = ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					companyIds,
					isBound.(bool)).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			} else {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND sports_company_id IN ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					companyIds).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			}
		} else {
			if isBound != nil {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND is_bound = ?`,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%",
					isBound.(bool)).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			} else {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields,
					"%"+queryCondition.LikeValue+"%",
					"%"+queryCondition.LikeValue+"%").
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			}
		}

	} else {
		whereFields := " 1=1 "
		if len(companyIds) > 0 {
			if isBound != nil {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND sports_company_id IN ? AND is_bound = ?`,
					companyIds, isBound.(bool)).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			} else {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND sports_company_id IN ?`,
					companyIds).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			}
		} else {

			if isBound != nil {
				result = Database.Model(&tag).Select(selectFields).Where(whereFields+` AND is_bound = ?`, isBound.(bool)).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			} else {
				result = Database.Model(&tag).Select(selectFields).
					Order("modify_datetime DESC").
					Limit(int(queryCondition.PageSize)).
					Offset(int(queryCondition.PageSize * (queryCondition.PageIndex - 1))).
					Find(&tags)
			}
		}
	}

	return tags, result.Error
}

func UpdateUWBTagSportsCompanyName(coId string, coName string) error {
	return Database.Table("uwb_tags").Where("sports_company_id = ?", coId).UpdateColumn("sports_company_name", coName).Error
}

func SetUWBTagSwimmerBoundStatus(uwbCode string, isBound bool) error {
	return Database.Table("uwb_tags").Where("code=?", uwbCode).UpdateColumn("is_bound", isBound).Error
}

func CancelAllUWBTagsBySwimmerId(swimmerId string) error {
	sql := `UPDATE uwb_tags SET is_bound=false WHERE code IN (SELECT uwb_tag_code FROM company_swimmers WHERE swimmer_id=?)`
	return Database.Exec(sql, swimmerId).Error
}

func DeleteUWBTags(codes []string) error {
	var tag entities.UWBTag
	result := Database.Delete(&tag, codes)
	return result.Error
}
