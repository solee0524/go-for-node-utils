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

// func (c Column) TableName() string {
// 	return "COLUMNS"
// }

func main() {
	var databaseName string
	var tabelName string
	var argsLen int

	argsLen = len(os.Args)
	if argsLen < 3 {
		fmt.Printf("Wrong args format! Must Be ./js_model_generator [database name] [table name]\n")
		return
	}
	databaseName = os.Args[1]
	tabelName = os.Args[2]

	// 数据库连接配置 for Mysql
	var username string = "username"
	var password string = "password"
	var host string = "localhost"
	var port string = "3306"
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

	// TODO
	// 组装Column Detail内容
	var columnDetails = ""
	for i, v := range columns {

		// 生成DataType
		var columnType = v.ColumnType
		switch v.DataType {
		case "int":
			columnType = strings.Replace(columnType, "int", "DataTypes.INTEGER", -1)
		case "varchar":
			columnType = strings.Replace(columnType, "varchar", "DataTypes.STRING", -1)
		case "text":
			columnType = strings.Replace(columnType, "text", "DataTypes.TEXT", -1)
		}

		var tmpString = "    " + v.ColumnName + ": {type: " + columnType

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
			tmpString += "}"
		} else {
			tmpString += "},\n"
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
	tmplModelFile, err := template.ParseFiles("model.tmpl")
	checkError(err)
	m := ModelInfos{TableName: tabelName, ColumnDetails: columnDetails}
	tmplModelFile.Execute(file, m)

}
