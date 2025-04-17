package util

import (
	"fmt"
	"log"

	"github.com/landaiqing/Go-TriTabula/entity"

	"github.com/landaiqing/go-dockit/document"
)

// ExportWord 用于导出Word文档的工具结构体
type ExportWord struct{}

// CreateDocument 创建一个新的Word文档
func (ew *ExportWord) CreateDocument(results []entity.Result) *document.Document {
	// 创建新文档
	doc := document.NewDocument()

	// 设置文档属性
	doc.SetTitle("数据库三线表")
	doc.SetCreator("Go-TriTabula")
	doc.SetDescription("使用go-dockit库创建数据库三线表")

	// 为每个表结果创建表格和标题
	for _, result := range results {
		// 创建表名标题段落
		tableTitlePara := doc.AddParagraph()
		tableTitlePara.SetAlignment("center")
		tableTitlePara.SetSpacingAfter(0)
		tableTitlePara.SetSpacingBefore(0)
		tableTitlePara.SetLineSpacing(1.5, "auto") // 设置1.5倍行距
		tableTitleRun := tableTitlePara.AddRun()
		tableTitleRun.AddText(result.TableName)
		tableTitleRun.SetBold(true)
		tableTitleRun.SetFontSize(21) // 五号字体约为10.5磅(21)
		tableTitleRun.SetFontFamily("宋体")

		// 获取字段映射
		fieldMapping := entity.GetFieldMapping()
		fieldNames := []string{"ColumnName", "ColumnType", "IsNullable", "ColumnKey", "ColumnDefault", "ColumnComment"}

		// 创建表格
		table := doc.AddTable(len(result.TableDetails)+1, len(fieldNames))
		table.SetWidth("100%", "pct") // 与文字齐宽
		table.SetAlignment("center")

		// 填充表头
		for i, fieldName := range fieldNames {
			cellPara := table.Rows[0].Cells[i].AddParagraph()
			cellPara.SetAlignment("center")
			cellPara.SetLineSpacing(1.5, "auto") // 1.5倍行距
			cellRun := cellPara.AddRun()
			cellRun.AddText(fieldMapping[fieldName])
			cellRun.SetBold(false)
			cellRun.SetFontSize(21) // 五号字体
			cellRun.SetFontFamily("宋体")
		}

		// 填充数据行
		for i, detail := range result.TableDetails {
			// 字段名
			para := table.Rows[i+1].Cells[0].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto") // 1.5倍行距
			cellRun := para.AddRun()
			cellRun.AddText(detail.ColumnName)
			cellRun.SetFontSize(21) // 五号字体
			cellRun.SetFontFamily("Times New Roman")

			// 类型
			para = table.Rows[i+1].Cells[1].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto")
			cellRun = para.AddRun()
			cellRun.AddText(detail.ColumnType)
			cellRun.SetFontSize(21)
			cellRun.SetFontFamily("Times New Roman") // 英文使用Times New Roman

			// 是否为空
			para = table.Rows[i+1].Cells[2].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto")
			cellRun = para.AddRun()
			// 将NO和YES转换为更易读的否和是
			isNullableText := "否"
			if detail.IsNullable == "YES" {
				isNullableText = "是"
			}
			cellRun.AddText(isNullableText)
			cellRun.SetFontSize(21)
			cellRun.SetFontFamily("宋体")

			// 索引
			para = table.Rows[i+1].Cells[3].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto")
			cellRun = para.AddRun()
			cellRun.AddText(detail.ColumnKey)
			cellRun.SetFontSize(21)
			cellRun.SetFontFamily("宋体")

			// 默认值
			para = table.Rows[i+1].Cells[4].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto")
			cellRun = para.AddRun()
			// 将"无"替换为"NULL"，使其更符合数据库术语
			defaultValue := "NULL"
			if detail.ColumnDefault != "无" {
				defaultValue = detail.ColumnDefault
			}
			cellRun.AddText(defaultValue)
			cellRun.SetFontSize(21)
			cellRun.SetFontFamily("宋体")

			// 说明
			para = table.Rows[i+1].Cells[5].AddParagraph()
			para.SetAlignment("center")
			para.SetLineSpacing(1.5, "auto")
			cellRun = para.AddRun()
			cellRun.AddText(detail.ColumnComment)
			cellRun.SetFontSize(21)
			cellRun.SetFontFamily("宋体")
		}

		// 设置三线表样式
		// 1. 首先清除所有默认边框
		table.SetBorders("all", "", 0, "")
		// 2. 顶线（表格顶部边框），1.5磅
		table.SetBorders("top", "single", 10, "000000")
		// 3. 表头分隔线（第一行底部边框），1磅
		for i := 0; i < len(fieldNames); i++ {
			table.Rows[0].Cells[i].SetBorders("bottom", "single", 4, "000000")
		}
		// 4. 底线（表格底部边框），1.5磅
		table.SetBorders("bottom", "single", 10, "000000")
		// 5. 显式设置内部边框为"none"
		table.SetBorders("insideH", "none", 0, "000000")
		table.SetBorders("insideV", "none", 0, "000000")

		// 添加空行
		doc.AddParagraph()
	}

	// 添加页脚（页码）
	footer := doc.AddFooterWithReference("default")
	footer.AddPageNumber()

	return doc
}

// ExportToFile 将文档导出到文件
func (ew *ExportWord) ExportToFile(doc *document.Document, filePath string) error {
	err := doc.Save(filePath)
	if err != nil {
		return fmt.Errorf("保存文档失败: %v", err)
	}

	log.Printf("文档已成功保存到: %s", filePath)
	return nil
}
