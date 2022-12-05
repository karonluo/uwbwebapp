package dao

func EnumSiteOwners(site_id string) []map[string]interface{} {
	var result []map[string]interface{}
	Database.Raw(`
	SELECT
		u."id" as sys_user_id, 
		u.cellphone,
		u.display_name
	FROM
		sites AS s,
		sys_users AS u,
		site_owners AS o
	WHERE
		s."id" = o.site_id AND
		u."id" = o.sys_user_id AND
		o.site_id = ? 
	ORDER BY u.modify_datetime DESC`,
		site_id).Scan(&result)
	return result
}
