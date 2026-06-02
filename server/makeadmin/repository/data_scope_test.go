package repository

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-makeadmin/model/makeadmin"
)

func TestApplyDataScopeFilterBuildsSelfAndOrgConstraint(t *testing.T) {
	db := newDryRunDataScopeDB(t)
	query := db.Model(&makeadmin.Admin{}).Where("ma_admin.delete_time = ?", 0)

	scoped := applyDataScopeFilter(db, query, 7, "ma_admin.id", DataScopeFilter{
		Enabled: true,
		Self:    true,
		AdminID: 3,
		OrgIDs:  []uint64{10, 11},
	})
	var admins []makeadmin.Admin
	statement := scoped.Find(&admins).Statement
	sql := statement.SQL.String()

	for _, want := range []string{
		"ma_admin.id = ?",
		"ma_admin.id IN (SELECT",
		"FROM `ma_admin_org`",
		"tenant_id = ? AND org_id IN",
		"status = ? AND delete_time = ?",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("scoped SQL missing %q: %s", want, sql)
		}
	}
}

func TestApplyDataScopeFilterBuildsNoAccessConstraint(t *testing.T) {
	db := newDryRunDataScopeDB(t)
	query := db.Model(&makeadmin.LoginLog{}).Where("tenant_id = ?", 7)

	scoped := applyDataScopeFilter(db, query, 7, "admin_id", DataScopeFilter{
		Enabled:  true,
		NoAccess: true,
	})
	var logs []makeadmin.LoginLog
	sql := scoped.Find(&logs).Statement.SQL.String()

	if !strings.Contains(sql, "1 = 0") {
		t.Fatalf("no-access SQL missing guard: %s", sql)
	}
}

func newDryRunDataScopeDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "go_makeadmin:gorm@tcp(127.0.0.1:9910)/go_makeadmin?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open dry-run gorm db: %v", err)
	}
	return db
}
