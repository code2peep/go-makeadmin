package makeadmin

const (
	SettingValueString = "string"
	SettingValueJSON   = "json"

	FileTypeImage = "image"
	FileTypeVideo = "video"

	CodegenTemplateCRUD = "crud"
	CodegenTemplateTree = "tree"
)

type Setting struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_setting_key,priority:1;comment:'tenant id'"`
	SettingGroup string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_setting_key,priority:2;comment:'setting group'"`
	SettingKey   string `gorm:"size:128;not null;default:'';uniqueIndex:uk_ma_setting_key,priority:3;comment:'setting key'"`
	SettingValue string `gorm:"type:text;not null;comment:'setting value'"`
	ValueType    string `gorm:"size:32;not null;default:'string';comment:'value type'"`
	IsPublic     uint8  `gorm:"not null;default:0;comment:'public readable flag'"`
	Remark       string `gorm:"size:255;not null;default:'';comment:'remark'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime   int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
}

func (Setting) TableName() string {
	return "ma_setting"
}

type DictType struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_dict_type_code_live,priority:1;comment:'dict type code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'dict type name'"`
	Remark     string `gorm:"size:255;not null;default:'';comment:'remark'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_dict_type_code_live,priority:2;comment:'delete time'"`
}

func (DictType) TableName() string {
	return "ma_dict_type"
}

type DictItem struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TypeID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_dict_item_value_live,priority:1;index;comment:'dict type id'"`
	ItemLabel  string `gorm:"size:64;not null;default:'';comment:'item label'"`
	ItemValue  string `gorm:"size:128;not null;default:'';uniqueIndex:uk_ma_dict_item_value_live,priority:2;comment:'item value'"`
	Remark     string `gorm:"size:255;not null;default:'';comment:'remark'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_dict_item_value_live,priority:3;comment:'delete time'"`
}

func (DictItem) TableName() string {
	return "ma_dict_item"
}

type FileCategory struct {
	ID         uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID   uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_file_category_code_live,priority:1;index:idx_ma_file_category_parent_sort,priority:1;comment:'tenant id'"`
	ParentID   uint64 `gorm:"not null;default:0;index:idx_ma_file_category_parent_sort,priority:2;comment:'parent category id'"`
	Code       string `gorm:"size:64;not null;default:'';uniqueIndex:uk_ma_file_category_code_live,priority:2;comment:'category code'"`
	Name       string `gorm:"size:64;not null;default:'';comment:'category name'"`
	FileType   string `gorm:"size:32;not null;default:'image';comment:'file type'"`
	Status     uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	Sort       uint16 `gorm:"not null;default:0;index:idx_ma_file_category_parent_sort,priority:3;comment:'sort order'"`
	CreateTime int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_file_category_code_live,priority:3;comment:'delete time'"`
}

func (FileCategory) TableName() string {
	return "ma_file_category"
}

type File struct {
	ID            uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID      uint64 `gorm:"not null;default:0;index:idx_ma_file_tenant_category,priority:1;comment:'tenant id'"`
	CategoryID    uint64 `gorm:"not null;default:0;index:idx_ma_file_tenant_category,priority:2;comment:'category id'"`
	OwnerAdminID  uint64 `gorm:"not null;default:0;index;comment:'owner admin id'"`
	FileType      string `gorm:"size:32;not null;default:'image';comment:'file type'"`
	StorageDriver string `gorm:"size:32;not null;default:'local';comment:'storage driver'"`
	OriginalName  string `gorm:"size:255;not null;default:'';comment:'original file name'"`
	FileName      string `gorm:"size:255;not null;default:'';comment:'stored file name'"`
	URI           string `gorm:"column:uri;size:512;not null;default:'';comment:'storage uri'"`
	URL           string `gorm:"column:url;size:512;not null;default:'';comment:'access url'"`
	MimeType      string `gorm:"size:128;not null;default:'';comment:'mime type'"`
	Ext           string `gorm:"size:32;not null;default:'';comment:'extension'"`
	Size          int64  `gorm:"not null;default:0;comment:'file size bytes'"`
	Checksum      string `gorm:"size:128;not null;default:'';index;comment:'checksum'"`
	Status        uint8  `gorm:"not null;default:1;comment:'status: 1=enabled, 0=disabled'"`
	CreateTime    int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime    int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime    int64  `gorm:"not null;default:0;comment:'delete time'"`
}

func (File) TableName() string {
	return "ma_file"
}

type CodegenTable struct {
	ID           uint64 `gorm:"primaryKey;comment:'primary key'"`
	TenantID     uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_codegen_table_live,priority:1;comment:'tenant id'"`
	SourceTable  string `gorm:"column:table_name;size:128;not null;default:'';uniqueIndex:uk_ma_codegen_table_live,priority:2;comment:'table name'"`
	TableComment string `gorm:"size:255;not null;default:'';comment:'table comment'"`
	ModuleName   string `gorm:"size:64;not null;default:'';comment:'module name'"`
	PackageName  string `gorm:"size:128;not null;default:'';comment:'package name'"`
	BusinessName string `gorm:"size:64;not null;default:'';comment:'business name'"`
	EntityName   string `gorm:"size:64;not null;default:'';comment:'entity name'"`
	FunctionName string `gorm:"size:64;not null;default:'';comment:'function name'"`
	AuthorName   string `gorm:"size:64;not null;default:'';comment:'author name'"`
	TemplateType string `gorm:"size:32;not null;default:'crud';comment:'template type'"`
	GenerateType string `gorm:"size:32;not null;default:'zip';comment:'generate type'"`
	GeneratePath string `gorm:"size:255;not null;default:'';comment:'generate path'"`
	Options      string `gorm:"type:text;not null;comment:'generator options json'"`
	Remark       string `gorm:"size:255;not null;default:'';comment:'remark'"`
	CreateTime   int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime   int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
	DeleteTime   int64  `gorm:"not null;default:0;uniqueIndex:uk_ma_codegen_table_live,priority:3;comment:'delete time'"`
}

func (CodegenTable) TableName() string {
	return "ma_codegen_table"
}

type CodegenColumn struct {
	ID            uint64 `gorm:"primaryKey;comment:'primary key'"`
	TableID       uint64 `gorm:"not null;default:0;uniqueIndex:uk_ma_codegen_column,priority:1;comment:'codegen table id'"`
	ColumnName    string `gorm:"size:128;not null;default:'';uniqueIndex:uk_ma_codegen_column,priority:2;comment:'column name'"`
	ColumnComment string `gorm:"size:255;not null;default:'';comment:'column comment'"`
	ColumnType    string `gorm:"size:64;not null;default:'';comment:'column type'"`
	ColumnLength  int    `gorm:"not null;default:0;comment:'column length'"`
	GoType        string `gorm:"size:64;not null;default:'';comment:'go type'"`
	GoField       string `gorm:"size:64;not null;default:'';comment:'go field'"`
	JSONField     string `gorm:"column:json_field;size:64;not null;default:'';comment:'json field'"`
	IsPK          uint8  `gorm:"column:is_pk;not null;default:0;comment:'primary key flag'"`
	IsIncrement   uint8  `gorm:"not null;default:0;comment:'increment flag'"`
	IsRequired    uint8  `gorm:"not null;default:0;comment:'required flag'"`
	IsInsert      uint8  `gorm:"not null;default:0;comment:'insert flag'"`
	IsEdit        uint8  `gorm:"not null;default:0;comment:'edit flag'"`
	IsList        uint8  `gorm:"not null;default:0;comment:'list flag'"`
	IsQuery       uint8  `gorm:"not null;default:0;comment:'query flag'"`
	QueryType     string `gorm:"size:32;not null;default:'=';comment:'query type'"`
	HTMLType      string `gorm:"column:html_type;size:32;not null;default:'';comment:'html type'"`
	DictType      string `gorm:"size:64;not null;default:'';comment:'dict type'"`
	Sort          uint16 `gorm:"not null;default:0;comment:'sort order'"`
	CreateTime    int64  `gorm:"autoCreateTime;not null;comment:'create time'"`
	UpdateTime    int64  `gorm:"autoUpdateTime;not null;comment:'update time'"`
}

func (CodegenColumn) TableName() string {
	return "ma_codegen_column"
}
