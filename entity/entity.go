package entity

// Result 表示数据库表的结构信息
type Result struct {
	TableSchema  string        `json:"tableSchema"`
	TableName    string        `json:"tableName"`
	TableDetails []TableDetail `json:"tableDetails"`
}

// TableDetail 表示表中列的详细信息
type TableDetail struct {
	ColumnName    string `json:"columnName" field:"字段名"`
	ColumnType    string `json:"columnType" field:"类型"`
	IsNullable    string `json:"isNullable" field:"是否为空"`
	ColumnKey     string `json:"columnKey" field:"索引"`
	ColumnDefault string `json:"columnDefault" field:"默认值"`
	ColumnComment string `json:"columnComment" field:"说明"`
}

// GetFieldMapping 返回字段名与显示名称的映射
func GetFieldMapping() map[string]string {
	return map[string]string{
		"ColumnName":    "字段名",
		"ColumnType":    "类型",
		"IsNullable":    "是否为空",
		"ColumnKey":     "索引",
		"ColumnDefault": "默认值",
		"ColumnComment": "说明",
	}
}
