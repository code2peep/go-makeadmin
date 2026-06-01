package makeadmin

const (
	GlobalTenantID uint64 = 0

	StatusDisabled uint8 = 0
	StatusEnabled  uint8 = 1

	MenuTypeCatalog = "catalog"
	MenuTypePage    = "page"
	MenuTypeAction  = "action"

	ScopeTypeAll       = "all"
	ScopeTypeSelf      = "self"
	ScopeTypeOrg       = "org"
	ScopeTypeOrgTree   = "org_tree"
	ScopeTypeCustomOrg = "custom_org"
)

// Admin is the P1 account identity draft. It is not wired into runtime yet.
type Admin struct {
	ID            uint64 `gorm:"primaryKey;comment:'primary key'"`
	Username      string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_admin_username_live,priority:1;comment:'username'"`
	PasswordHash  string `gorm:"size:255;not null;default:'';comment:'password hash'"`
	PasswordSalt  string `gorm:"size:64;not null;default:'';comment:'password salt'"`
	IsSuper       uint8  `gorm:"not null;default:0;comment:'super admin flag'"`
	Status        uint8  `gorm:"not null;default:1;index;comment:'status: 1=enabled, 0=disabled'"`
	LastLoginIP   string `gorm:"column:last_login_ip;size:64;not null;default:'';comment:'last login ip'"`
	LastLoginTime int64  `gorm:"not null;default:0;comment:'last login time'"`
	CreateTime    int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime    int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime    int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_username_live,priority:2;comment:'delete time'"`
}

func (Admin) TableName() string {
	return "ma_admin"
}

type AdminProfile struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	AdminID    uint64 `gorm:"not null;default:0;uniqueIndex;comment:'admin id'"`
	Nickname   string `gorm:"size:64;not null;default:'';comment:'nickname'"`
	Avatar     string `gorm:"size:255;not null;default:'';comment:'avatar'"`
	Email      string `gorm:"size:128;not null;default:'';comment:'email'"`
	Mobile     string `gorm:"size:32;not null;default:'';comment:'mobile'"`
	Remark     string `gorm:"size:255;not null;default:'';comment:'remark'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
}

func (AdminProfile) TableName() string {
	return "ma_admin_profile"
}

type Tenant struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_tenant_code_live,priority:1;comment:'tenant code'"`
	Name       string `gorm:"size:128;not null;default:'';comment:'tenant name'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_tenant_code_live,priority:2;comment:'delete time'"`
}

func (Tenant) TableName() string {
	return "ma_tenant"
}

type TenantMember struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_tenant_member_live,priority:1;comment:'tenant id'"`
	AdminID    uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_tenant_member_live,priority:2;index;comment:'admin id'"`
	MemberType string `gorm:"size:32;not null;default:'member';comment:'member type'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_tenant_member_live,priority:3;comment:'delete time'"`
}

func (TenantMember) TableName() string {
	return "ma_tenant_member"
}

type TenantSetting struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_tenant_setting_key,priority:1;comment:'tenant id'"`
	SettingKey   string `gorm:"size:128;not null;default:'';uniqueIndex:uk_ma_tenant_setting_key,priority:2;comment:'setting key'"`
	SettingValue string `gorm:"type:text;not null;comment:'setting value'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime   int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
}

func (TenantSetting) TableName() string {
	return "ma_tenant_setting"
}

type Role struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_code_live,priority:1;index:idx_ma_role_tenant_status,priority:1;comment:'tenant id'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_role_code_live,priority:2;comment:'role code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'role name'"`
	Remark     string `gorm:"size:255;not null;default:'';comment:'remark'"`
	IsSystem   uint8  `gorm:"not null;default:0;comment:'system role flag'"`
	Status     uint8  `gorm:"not null;default:1;index:idx_ma_role_tenant_status,priority:2;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_role_code_live,priority:3;comment:'delete time'"`
}

func (Role) TableName() string {
	return "ma_role"
}

type AdminRole struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_role,priority:1;comment:'tenant id'"`
	AdminID    uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_role,priority:2;comment:'admin id'"`
	RoleID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_role,priority:3;index;comment:'role id'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
}

func (AdminRole) TableName() string {
	return "ma_admin_role"
}

type Permission struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	Code       string `gorm:"size:128;not null;default:'';uniqueIndex;comment:'permission code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'permission name'"`
	Module     string `gorm:"size:64;not null;default:'';index:idx_ma_permission_module,priority:1;comment:'module'"`
	Resource   string `gorm:"size:64;not null;default:'';index:idx_ma_permission_module,priority:2;comment:'resource'"`
	Action     string `gorm:"size:64;not null;default:'';index:idx_ma_permission_module,priority:3;comment:'action'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
}

func (Permission) TableName() string {
	return "ma_permission"
}

type RolePermission struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_permission,priority:1;comment:'tenant id'"`
	RoleID       uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_permission,priority:2;comment:'role id'"`
	PermissionID uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_permission,priority:3;index;comment:'permission id'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
}

func (RolePermission) TableName() string {
	return "ma_role_permission"
}

type Menu struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	ParentID   uint64 `gorm:"not null;default:0;index:idx_ma_menu_parent_sort,priority:1;comment:'parent menu id'"`
	MenuType   string `gorm:"size:16;not null;default:'page';comment:'menu type'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'menu name'"`
	Icon       string `gorm:"size:64;not null;default:'';comment:'menu icon'"`
	RoutePath  string `gorm:"size:128;not null;default:'';comment:'route path'"`
	RouteName  string `gorm:"size:128;not null;default:'';comment:'route name'"`
	Component  string `gorm:"size:255;not null;default:'';comment:'component path'"`
	Redirect   string `gorm:"size:255;not null;default:'';comment:'redirect path'"`
	ActivePath string `gorm:"size:128;not null;default:'';comment:'active path'"`
	Meta       string `gorm:"type:text;not null;comment:'route meta json'"`
	IsVisible  uint8  `gorm:"not null;default:1;comment:'visible flag'"`
	IsCache    uint8  `gorm:"not null;default:0;comment:'cache flag'"`
	Status     uint8  `gorm:"not null;default:1;index;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;index:idx_ma_menu_parent_sort,priority:2;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;comment:'delete time'"`
}

func (Menu) TableName() string {
	return "ma_menu"
}

type MenuPermission struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	MenuID       uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_menu_permission,priority:1;comment:'menu id'"`
	PermissionID uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_menu_permission,priority:2;index;comment:'permission id'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
}

func (MenuPermission) TableName() string {
	return "ma_menu_permission"
}

type OrgUnit struct {
	ID            uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID      uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_org_unit_code_live,priority:1;index:idx_ma_org_unit_parent_sort,priority:1;comment:'tenant id'"`
	ParentID      uint64 `gorm:"not null;default:0;index:idx_ma_org_unit_parent_sort,priority:2;comment:'parent org id'"`
	Code          string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_org_unit_code_live,priority:2;comment:'org code'"`
	Name          string `gorm:"size:128;not null;default:'';comment:'org name'"`
	LeaderAdminID uint64 `gorm:"not null;default:0;comment:'leader admin id'"`
	Status        uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort          uint16 `gorm:"not null;default:0;index:idx_ma_org_unit_parent_sort,priority:3;comment:'sort order'"`
	CreateTime    int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime    int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime    int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_org_unit_code_live,priority:3;comment:'delete time'"`
}

func (OrgUnit) TableName() string {
	return "ma_org_unit"
}

type Position struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_position_code_live,priority:1;comment:'tenant id'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_position_code_live,priority:2;comment:'position code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'position name'"`
	Remark     string `gorm:"size:255;not null;default:'';comment:'remark'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_position_code_live,priority:3;comment:'delete time'"`
}

func (Position) TableName() string {
	return "ma_position"
}

type AdminOrg struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_org_live,priority:1;index:idx_ma_admin_org_org,priority:1;comment:'tenant id'"`
	AdminID    uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_org_live,priority:2;comment:'admin id'"`
	OrgID      uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_org_live,priority:3;index:idx_ma_admin_org_org,priority:2;comment:'org id'"`
	PositionID uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_org_live,priority:4;comment:'position id'"`
	IsPrimary  uint8  `gorm:"not null;default:0;comment:'primary org flag'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_admin_org_live,priority:5;comment:'delete time'"`
}

func (AdminOrg) TableName() string {
	return "ma_admin_org"
}

type DataScope struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_data_scope_code_live,priority:1;comment:'tenant id'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_data_scope_code_live,priority:2;comment:'scope code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'scope name'"`
	ScopeType  string `gorm:"size:32;not null;default:'self';comment:'scope type'"`
	ScopeValue string `gorm:"type:text;not null;comment:'scope value json'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_data_scope_code_live,priority:3;comment:'delete time'"`
}

func (DataScope) TableName() string {
	return "ma_data_scope"
}

type RoleDataScope struct {
	ID          uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID    uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_data_scope,priority:1;comment:'tenant id'"`
	RoleID      uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_data_scope,priority:2;comment:'role id'"`
	DataScopeID uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_role_data_scope,priority:3;comment:'data scope id'"`
	CreateTime  int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
}

func (RoleDataScope) TableName() string {
	return "ma_role_data_scope"
}

type LoginLog struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;index:idx_ma_login_log_tenant_time,priority:1;comment:'tenant id'"`
	AdminID    uint64 `gorm:"not null;default:0;index:idx_ma_login_log_admin_time,priority:1;comment:'admin id'"`
	Username   string `gorm:"size:64;not null;default:'';comment:'username'"`
	IP         string `gorm:"column:ip;size:64;not null;default:'';comment:'ip address'"`
	OS         string `gorm:"column:os;size:64;not null;default:'';comment:'operating system'"`
	Browser    string `gorm:"size:64;not null;default:'';comment:'browser'"`
	Status     uint8  `gorm:"not null;default:0;comment:'status: 1=success, 0=failure'"`
	Message    string `gorm:"size:255;not null;default:'';comment:'message'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;index:idx_ma_login_log_tenant_time,priority:2;index:idx_ma_login_log_admin_time,priority:2;comment:'create time'"`
}

func (LoginLog) TableName() string {
	return "ma_login_log"
}

type AuditLog struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID     uint64 `gorm:"not null;default:0;index:idx_ma_audit_log_tenant_time,priority:1;comment:'tenant id'"`
	AdminID      uint64 `gorm:"not null;default:0;index:idx_ma_audit_log_admin_time,priority:1;comment:'admin id'"`
	TraceID      string `gorm:"size:64;not null;default:'';index;comment:'trace id'"`
	Action       string `gorm:"size:128;not null;default:'';comment:'action code'"`
	Method       string `gorm:"size:16;not null;default:'';comment:'http method'"`
	Path         string `gorm:"size:255;not null;default:'';comment:'request path'"`
	RequestBody  string `gorm:"type:text;not null;comment:'request body'"`
	ResponseCode int    `gorm:"not null;default:0;comment:'response code'"`
	Error        string `gorm:"type:text;not null;comment:'error message'"`
	Status       uint8  `gorm:"not null;default:0;comment:'status: 1=success, 0=failure'"`
	StartTime    int64  `gorm:"not null;default:0;comment:'start time'"`
	EndTime      int64  `gorm:"not null;default:0;comment:'end time'"`
	DurationMS   int64  `gorm:"column:duration_ms;not null;default:0;comment:'duration milliseconds'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;index:idx_ma_audit_log_tenant_time,priority:2;index:idx_ma_audit_log_admin_time,priority:2;comment:'create time'"`
}

func (AuditLog) TableName() string {
	return "ma_audit_log"
}
