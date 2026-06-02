package repository

import (
	"strings"

	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

type DataScopeFilter struct {
	Enabled  bool
	All      bool
	Self     bool
	NoAccess bool
	AdminID  uint64
	OrgIDs   []uint64
}

func (scope DataScopeFilter) IsRestricted() bool {
	return scope.Enabled && !scope.All
}

func applyDataScopeFilter(db *gorm.DB, query *gorm.DB, tenantID uint64, adminColumn string, scope DataScopeFilter) *gorm.DB {
	if !scope.Enabled || scope.All {
		return query
	}
	if scope.NoAccess {
		return query.Where("1 = 0")
	}
	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)
	if scope.Self && scope.AdminID > 0 {
		conditions = append(conditions, adminColumn+" = ?")
		args = append(args, scope.AdminID)
	}
	if len(scope.OrgIDs) > 0 {
		orgAdminSubquery := db.Model(&makeadmin.AdminOrg{}).
			Select("admin_id").
			Where("tenant_id = ? AND org_id IN ? AND status = ? AND delete_time = ?", tenantID, scope.OrgIDs, makeadmin.StatusEnabled, 0)
		conditions = append(conditions, adminColumn+" IN (?)")
		args = append(args, orgAdminSubquery)
	}
	if len(conditions) == 0 {
		return query.Where("1 = 0")
	}
	return query.Where("("+strings.Join(conditions, " OR ")+")", args...)
}
