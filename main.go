package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/landaiqing/Go-TriTabula/util"
)

func main() {
	// 解析命令行参数
	var configPath string
	var outputFile string

	// 设置命令行参数
	flag.StringVar(&configPath, "config", "config/dbconfig.json", "数据库配置文件路径")
	flag.StringVar(&outputFile, "output", "output.docx", "输出文件名")
	flag.Parse()

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	// 加载数据库配置
	configFullPath := filepath.Join(wd, configPath)
	config, err := util.LoadDBConfig(configFullPath)
	if err != nil {
		log.Fatalf("加载数据库配置失败: %v", err)
	}

	fmt.Printf("正在连接数据库: %s\n", config.Database)

	// 获取数据库连接
	db, err := util.GetDatabaseConnection(config)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	fmt.Println("数据库连接成功")

	// 获取表结构信息
	fmt.Printf("正在获取数据库 %s 的表结构信息...\n", config.Database)
	results, err := util.GetTableDetails(db, config.Database)
	if err != nil {
		log.Fatalf("获取表结构信息失败: %v", err)
	}

	fmt.Printf("成功获取 %d 个表的结构信息\n", len(results))

	// 创建导出工具
	exportWord := &util.ExportWord{}

	// 创建文档
	fmt.Println("正在生成Word文档...")
	doc := exportWord.CreateDocument(results)

	// 导出文档
	outputFullPath := filepath.Join(wd, outputFile)
	err = exportWord.ExportToFile(doc, outputFullPath)
	if err != nil {
		log.Fatalf("导出文档失败: %v", err)
	}

	fmt.Printf("文档生成成功: %s\n", outputFullPath)
}
