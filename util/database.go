package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log" // 添加 log 包
	"os"
	"path/filepath"

	"github.com/landaiqing/Go-TriTabula/entity"

	"github.com/go-sql-driver/mysql"
)

// DBConfig 数据库配置结构
type DBConfig struct {
	Driver       string `json:"driver"`
	URL          string `json:"url"`
	Database     string `json:"database"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

// LoadDBConfig 从配置文件加载数据库配置
func LoadDBConfig(configPath string) (*DBConfig, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("获取配置文件绝对路径失败: %v", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config DBConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// GetDatabaseConnection 获取数据库连接
func GetDatabaseConnection(config *DBConfig) (*sql.DB, error) {
	// 配置MySQL连接
	mysqlConfig := mysql.Config{
		User:                 config.Username,
		Passwd:               config.Password,
		Net:                  "tcp",
		Addr:                 config.URL,
		DBName:               config.Database,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// 打开数据库连接
	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("测试数据库连接失败: %v", err)
	}

	return db, nil
}

// GetTableDetails 获取数据库表结构详情
func GetTableDetails(db *sql.DB, databaseName string) ([]entity.Result, error) {
	query := `
		SELECT table_schema, table_name, column_name, column_type, column_key, 
		       is_nullable, column_default, column_comment, character_set_name, EXTRA
		FROM information_schema.columns 
		WHERE table_schema = ? 
		ORDER BY table_name, ORDINAL_POSITION
	`

	// 执行查询
	rows, err := db.Query(query, databaseName)
	if err != nil {
		return nil, fmt.Errorf("查询表结构失败: %v", err)
	}
	defer rows.Close()

	// 用于存储结果的映射和列表
	tableMap := make(map[string]int) // Map tableName to index in results slice
	var results []entity.Result

	// 遍历查询结果
	for rows.Next() {
		var tableSchema, tableName, columnName, columnType, columnKey string
		var isNullable, columnDefault, columnComment, characterSetName, extra sql.NullString

		// 扫描行数据
		err := rows.Scan(
			&tableSchema, &tableName, &columnName, &columnType, &columnKey,
			&isNullable, &columnDefault, &columnComment, &characterSetName, &extra,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描行数据失败: %v", err)
		}
		log.Printf("读取到行: Table=%s, Column=%s\n", tableName, columnName) // 添加行读取日志

		// 处理表信息
		idx, exists := tableMap[tableName]
		if !exists {
			// 创建新的表结果并添加到 results slice
			newResult := entity.Result{
				TableSchema:  tableSchema,
				TableName:    tableName,
				TableDetails: []entity.TableDetail{},
			}
			results = append(results, newResult)
			idx = len(results) - 1    // Get the index of the newly added result
			tableMap[tableName] = idx // Store the index in the map
		}

		// 处理列信息
		columnDefaultValue := "无"
		if columnDefault.Valid {
			columnDefaultValue = columnDefault.String
		}

		columnKeyValue := "无"
		if columnKey != "" {
			columnKeyValue = columnKey
		}

		columnCommentValue := ""
		if columnComment.Valid {
			columnCommentValue = columnComment.String
		}

		isNullableValue := "NO"
		if isNullable.Valid {
			isNullableValue = isNullable.String
		}

		// 创建列详情
		tableDetail := entity.TableDetail{
			ColumnName:    columnName,
			ColumnType:    columnType,
			ColumnKey:     columnKeyValue,
			IsNullable:    isNullableValue,
			ColumnDefault: columnDefaultValue,
			ColumnComment: columnCommentValue,
		}

		// 添加到表的列列表 (直接修改 results 切片中的元素)
		results[idx].TableDetails = append(results[idx].TableDetails, tableDetail)
	}

	// 检查是否有错误发生
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果集时发生错误: %v", err)
	}

	log.Printf("GetTableDetails 完成, 共处理 %d 个表的结果。\n", len(results)) // 添加最终结果日志
	return results, nil
}
