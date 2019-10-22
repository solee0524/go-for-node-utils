package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
	"strings"
	"text/template"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

type Column struct {
	TableSchema     string `xorm:"TABLE_SCHEMA varchar(64) notnull"`
	TableName       string `xorm:"TABLE_NAME varchar(64) notnull"`
	ColumnName      string `xorm:"COLUMN_NAME varchar(64) notnull"`
	OrdinalPosition int    `xorm:"ORDINAL_POSITION bigint(21) notnull"`
	IsNullable      string `xorm:"IS_NULLABLE varchar(3) notnull"`
	DataType        string `xorm:"DATA_TYPE varchar(64) notnull"`
	ColumnType      string `xorm:"COLUMN_TYPE longtext notnull"`
	ColumnKey       string `xorm:"COLUMN_KEY varchar(3) notnull"`
	Extra           string `xorm:"EXTRA varchar(27) notnull"`
	ColumnComment   string `xorm:"COLUMN_COMMENT varchar(1024) notnull"`
}

func main() {
	var databaseName string
	var tabelName string
	var argsLen int
	var modelStyle string

	argsLen = len(os.Args)
	if argsLen < 8 {
		fmt.Printf("Wrong args format! Must Be ./js_model_generator [host] [port] [username] [password] [database name] [table name] [model_style:1(undeline_mode) 2(camel_mode)]\n")
		return
	}
	databaseName = os.Args[5]
	tabelName = os.Args[6]
	modelStyle = os.Args[7]

	// 数据库连接配置 for Mysql
	var host string = os.Args[1]
	var port string = os.Args[2]
	var username string = os.Args[3]
	var password string = os.Args[4]
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?timeout=10s", username, password, host, port)
	engine, err := xorm.NewEngine("mysql", connectStr)
	checkError(err)

	var columns []Column
	err = engine.Table("COLUMNS").Where("TABLE_NAME = ?", tabelName).And("TABLE_SCHEMA LIKE ?", databaseName).Cols("TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME,  ORDINAL_POSITION, IS_NULLABLE, DATA_TYPE, COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT").Asc("ORDINAL_POSITION").Find(&columns)
	checkError(err)

	// 创建modle文件
	file, err := os.Create(tabelName + ".js")
	checkError(err)
	defer file.Close()

	// 组装Column Detail内容
	var columnDetails = ""
	for i, v := range columns {

		// 生成DataType
		var columnType = v.ColumnType
		fmt.Sprintf("%s",v.DataType)
		switch v.DataType {
		case "int":
			columnType = strings.Replace(columnType, "int", "DataTypes.INTEGER", -1)
		case "varchar":
			columnType = strings.Replace(columnType, "varchar", "DataTypes.STRING", -1)
		case "text":
		case "mediumtext":
		case "longtext":
			columnType = strings.Replace(columnType, "text", "DataTypes.TEXT", -1)
		}

		columnName := v.ColumnName;
		if modelStyle == "2" {
			// 修改栏目名称为驼峰模式
			temp := strings.Split(columnName, "_")
			var s string
			for i, v := range temp {
				if (i != 0) {
					vv := []rune(v)
					if len(vv) > 0 {
						if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
							vv[0] -= 32
						}
						s += string(vv)
					}
				} else {
						s += v
				}
			}
			columnName = s;
		}
		var tmpString = "    " + columnName + ": {type: " + columnType

		// 生成allowNull
		if v.IsNullable == "YES" {
			tmpString += ", allowNull: true"
		} else {
			tmpString += ", allowNull: false"
		}

		// 判断是否有自增属性
		if strings.Contains(v.Extra, "auto_increment") {
			tmpString += ", autoIncrement: true"
		}
		// 判断是否是主键
		if strings.Contains(v.ColumnKey, "PRI") {
			tmpString += ", primaryKey: true"
		}

		if len(columns)-1 == i {
			tmpString += ", field: '" + v.ColumnName + "'}"
		} else {
			tmpString += ", field: '" + v.ColumnName + "'},\n"
		}

		// 处理注释
		var columnComments = strings.Split(v.ColumnComment, "\n")
		for _, y := range columnComments {
			columnDetails += "    // " + y + "\n"
		}

		columnDetails += tmpString
	}

	// 使用模版文件生成需要的modle内容到创建的文件中
	type ModelInfos struct {
		TableName     string
		ColumnDetails string
	}
	modelFileName := "model.tmpl"
	if modelStyle == "2" {
		modelFileName = "model_camel.tmpl"
	}
	tmplModelFile, err := template.ParseFiles(modelFileName)
	checkError(err)
	m := ModelInfos{TableName: tabelName, ColumnDetails: columnDetails}
	tmplModelFile.Execute(file, m)
}
