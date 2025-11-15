package utils

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"sync"

	"github.com/fuermoya/design/server/global"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const (
	totalMax = 8000000 //数据总数
	sheetMax = 500000  //每页最大行数1048576
	poolMax  = 5       //最大协程
)

func ExportExcel(data [][]interface{}, c *gin.Context, filename string) error {
	if len(data) <= 0 {
		return fmt.Errorf("无数据")
	}
	if len(data[0])*len(data) > totalMax {
		return streamExportExcel(data, c, filename)
	}
	return ordExportExcel(data, c, filename)
}

// 普通导出excel
func ordExportExcel(data [][]interface{}, c *gin.Context, filename string) error {
	// 创建一个新的 Excel 文件
	file := excelize.NewFile()
	//多建一个sheet
	for i := 0; i <= (len(data[0]))/sheetMax; i++ {
		sheetStr := fmt.Sprintf("sheet%d", i+1)
		_, err := file.NewSheet(sheetStr)
		if err != nil {
			return err
		}
	}
	pMax := len(data)
	if pMax > poolMax {
		pMax = poolMax
	}
	pool := NewWorkerPool(pMax)
	fun := func(colIdx int, colData []interface{}) {
		index := 1
		name, _ := excelize.ColumnNumberToName(colIdx + 1)
		dataLen := len(colData) - 1
		sheet := dataLen / sheetMax

		for i := 0; i <= sheet; i++ {
			sheetName := fmt.Sprintf("sheet%d", i+1)
			jLen := sheetMax
			if i == sheet {
				jLen = dataLen % sheetMax
			}

			//设置表头
			cell, err := excelize.CoordinatesToCellName(colIdx+1, 1)
			if err != nil {
				global.LOG.Error("excel设置表头失败", zap.Error(err))
				break
			}
			file.SetCellValue(sheetName, cell, colData[0])

			//写入数据
			for j := 0; j < jLen; j++ {
				cell2, err := excelize.CoordinatesToCellName(colIdx+1, j+2)
				if err != nil {
					global.LOG.Error("excel写入数据失败", zap.Error(err))
					break
				}
				file.SetCellValue(sheetName, cell2, colData[index])
				index++
			}

			// 设置列宽度为 20
			err = file.SetColWidth(sheetName, name, name, 20)
			if err != nil {
				global.LOG.Error("excel设置列宽度失败", zap.Error(err))
			}
		}
	}

	// 将数据按列写入工作表
	for k, v := range data {
		pool.Do(func() {
			fun(k, v)
		})
	}
	pool.Wait()
	// 将文件写入响应
	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename)+".xlsx")
	err := file.Write(c.Writer)
	return err
}

// 流式导出
func streamExportExcel(matrix [][]interface{}, c *gin.Context, filename string) error {
	header, data := rowToCol(matrix)
	colLen := len(data)
	matrixLen := len(matrix)
	matrix = nil
	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	styleID, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "000000"}})

	pMax := colLen / sheetMax
	if pMax > poolMax {
		pMax = poolMax
	}
	pool := NewWorkerPool(pMax)
	var fun = func(sw *excelize.StreamWriter, i int, data [][]interface{}, header []interface{}) {
		//设置列宽
		sw.SetColWidth(1, matrixLen, 20)
		minLen := i * sheetMax
		maxLen := (i + 1) * sheetMax
		if maxLen > colLen {
			maxLen = colLen%sheetMax + minLen
		}

		//写入表头
		firstCell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			global.LOG.Error("excelize.CoordinatesToCellName", zap.Error(err))
		}

		if err := sw.SetRow(firstCell, header, excelize.RowOpts{Height: 16, Hidden: false, StyleID: styleID}); err != nil {
			global.LOG.Error("sw.SetRow", zap.Error(err))
		}

		//写入数据
		index := 2
		for j := minLen; j < maxLen; j++ {
			cell, err := excelize.CoordinatesToCellName(1, index)
			if err != nil {
				fmt.Println(err)
				break
			}
			if err := sw.SetRow(cell, data[j]); err != nil {
				fmt.Println(err)
				break
			}
			index++
		}
		if err = sw.Flush(); err != nil {
			global.LOG.Error("sw.Flush", zap.Error(err))
		}
	}

	//多建一个sheet
	for i := 0; i <= (colLen)/sheetMax; i++ {
		sheetStr := fmt.Sprintf("sheet%d", i+1)
		_, e := f.NewSheet(sheetStr)
		if e != nil {
			return e
		}
		sw, e := f.NewStreamWriter(sheetStr)
		if e != nil {
			return e
		}

		pool.Do(func() {
			fun(sw, i, data, header)
		})
	}

	pool.Wait()
	// 将文件写入响应
	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename)+".xlsx")
	err := f.Write(c.Writer)
	return err
}

// ReadExcel 读取excel
func ReadExcel(mufile multipart.File) ([]map[string]string, error) {
	f, err := excelize.OpenReader(mufile)
	if err != nil {
		global.LOG.Error("OpenReader Excel Error!", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			global.LOG.Error("OpenReader Excel Error!", zap.Error(err))
		}
	}()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		global.LOG.Error("GetRows Sheet1 Error!", zap.Error(err))
		return nil, err
	}

	data := make([]map[string]string, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		m := make(map[string]string, len(rows[i]))
		for i2, colCell := range rows[i] {
			m[rows[0][i2]] = colCell
		}
		data = append(data, m)
	}
	return data, nil
}

// 行转列 并提取表头
func rowToCol(matrix [][]interface{}) (header []interface{}, data [][]interface{}) {
	rows := len(matrix)
	if rows == 0 {
		return
	}

	header = make([]interface{}, len(matrix))
	result := make([][]interface{}, len(matrix[0])-1)
	for i := range result {
		result[i] = make([]interface{}, rows)
	}

	var wg sync.WaitGroup
	wg.Add(rows)

	for i, row := range matrix {
		go func(i int, row []interface{}) {
			defer wg.Done()
			rowLen := len(row)
			for j := 1; j < rowLen; j++ {
				result[j-1][i] = row[j]
			}
			for j := 0; j < rowLen; j++ {
				if j > 0 {
					break
				}
				header[i] = row[j]
			}

		}(i, row)
	}
	wg.Wait()

	return header, result
}
