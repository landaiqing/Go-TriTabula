# Go-TriTabula

这是一个使用Go语言实现的数据库三线表生成工具，用于从MySQL数据库中提取表结构信息并生成Word格式的三线表文档。

## 功能特点

- 连接MySQL数据库并获取表结构信息
- 将表结构信息转换为三线表格式
- 生成Word文档格式的三线表
- 支持自定义数据库连接参数

## 使用方法

1. 配置`config/dbconfig.json`文件中的数据库连接信息
2. 运行主程序：`go run main.go`
3. 生成的Word文档将保存在程序运行目录下

## 项目结构

```
Go-TriTabula/
├── config/         # 配置文件目录
├── entity/         # 实体定义
├── util/           # 工具类
└── main.go         # 主程序入口
```