package ngen

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/go-lazyer/north/nfile"
	"github.com/go-lazyer/north/nstring"
)

type Generator struct {
	dsn        string
	project    string
	driverName string
}

func NewGen() *Generator {
	return &Generator{}
}
func (gen *Generator) Dsn(dsn string) *Generator {
	gen.dsn = dsn
	return gen
}
func (gen *Generator) DriverName(driverName string) *Generator {
	gen.driverName = driverName
	return gen
}
func (gen *Generator) Project(project string) *Generator {
	gen.project = project
	return gen
}

type Module struct {
	TableName           string //表名
	TableNameUpperCamel string //表名的大驼峰
	TableNameLowerCamel string //表名的小驼峰
	ModulePath          string //模块名用于生成文件名
	Fields              []Field
	PrimaryKeyFields    []Field //主键
	Model               bool
	ModelFilePath       string //全路径，不包含文件名
	ModelFileName       string //只有文件名
	ModelPackageName    string //只有包名，不包含文件名
	ModelPackagePath    string //包含完整的包名

	Extend            bool
	ExtendFilePath    string //全路径，不包含文件名
	ExtendFileName    string //只有文件名
	ExtendPackageName string //只有包名，不包含文件名
	ExtendPackagePath string //包含完整的包名

	View            bool
	ViewFilePath    string
	ViewFileName    string
	ViewPackageName string
	ViewPackagePath string

	Param            bool
	ParamFilePath    string
	ParamFileName    string
	ParamPackageName string
	ParamPackagePath string

	Dao            bool
	DaoFilePath    string
	DaoFileName    string
	DaoPackageName string
	DaoPackagePath string

	Service            bool
	ServiceFilePath    string
	ServiceFileName    string
	ServicePackageName string
	ServicePackagePath string

	Controller            bool
	ControllerFilePath    string
	ControllerFileName    string
	ControllerPackageName string
	ControllerPackagePath string
	CreateTime            string
}

type Field struct {
	ColumnName           string         //msyql字段名 user_id
	ColumnNameLowerCamel string         //小驼峰 userId
	ColumnNameUpper      string         //字段名大写 USER_ID
	ColumnType           string         //msql 类型 varchat
	ColumnDefault        sql.NullString //默认值
	IsNullable           int            //允许为空
	IsPrimaryKey         int            //是否主键
	FieldName            string         //实体名称 大驼峰  UserId
	FieldNullType        string         //实体golang Null类型 sql.NullString
	FieldNullTypeValue   string         //实体golang Null类型 取值  String
	FieldType            string         //golang 类型  string
	FieldTypeDefault     string         //golang 类型  的默认值
	FieldOrmTag          string         //tag orm:
	FieldJsonTag         string         //tag json
	FieldFormTag         string         //tag form
	FieldDefaultTag      string         //tag 默认值
	Comment              string         //表中字段注释
}

var dbType = map[string]goType{
	"int":                {"int64", "sql.NullInt64", "Int64", "0"},
	"int4":               {"int64", "sql.NullInt64", "Int64", "0"},
	"int8":               {"int64", "sql.NullInt64", "Int64", "0"},
	"integer":            {"int64", "sql.NullInt64", "Int64", "0"},
	"tinyint":            {"int64", "sql.NullInt64", "Int64", "0"},
	"smallint":           {"int64", "sql.NullInt64", "Int64", "0"},
	"mediumint":          {"int64", "sql.NullInt64", "Int64", "0"},
	"bigint":             {"int64", "sql.NullInt64", "Int64", "0"},
	"int unsigned":       {"int64", "sql.NullInt64", "Int64", "0"},
	"integer unsigned":   {"int64", "sql.NullInt64", "Int64", "0"},
	"tinyint unsigned":   {"int64", "sql.NullInt64", "Int64", "0"},
	"smallint unsigned":  {"int64", "sql.NullInt64", "Int64", "0"},
	"mediumint unsigned": {"int64", "sql.NullInt64", "Int64", "0"},
	"bigint unsigned":    {"int64", "sql.NullInt64", "Int64", "0"},
	"bit":                {"int64", "sql.NullInt64", "Int64", "0"},
	"bool":               {"bool", "sql.NullBool", "Bool", "false"},
	"enum":               {"string", "sql.NullString", "String", "\"\""},
	"set":                {"string", "sql.NullString", "String", "\"\""},
	"varchar":            {"string", "sql.NullString", "String", "\"\""},
	"char":               {"string", "sql.NullString", "String", "\"\""},
	"tinytext":           {"string", "sql.NullString", "String", "\"\""},
	"mediumtext":         {"string", "sql.NullString", "String", "\"\""},
	"text":               {"string", "sql.NullString", "String", "\"\""},
	"longtext":           {"string", "sql.NullString", "String", "\"\""},
	"blob":               {"string", "sql.NullString", "String", "\"\""},
	"tinyblob":           {"string", "sql.NullString", "String", "\"\""},
	"mediumblob":         {"string", "sql.NullString", "String", "\"\""},
	"longblob":           {"string", "sql.NullString", "String", "\"\""},
	"date":               {"time.Time", "sql.NullTime", "Time", "nil"},
	"datetime":           {"time.Time", "sql.NullTime", "Time", "nil"},
	"timestamp":          {"time.Time", "sql.NullTime", "Time", "nil"},
	"timestamptz":        {"time.Time", "sql.NullTime", "Time", "nil"},
	"time":               {"time.Time", "sql.NullTime", "Time", "nil"},
	"timetz":             {"time.Time", "sql.NullTime", "Time", "nil"},
	"float":              {"float64", "sql.NullFloat64", "Float64", "0"},
	"float4":             {"float64", "sql.NullFloat64", "Float64", "0"},
	"float8":             {"float64", "sql.NullFloat64", "Float64", "0"},
	"double":             {"float64", "sql.NullFloat64", "Float64", "0"},
	"decimal":            {"float64", "sql.NullFloat64", "Float64", "0"},
	"binary":             {"string", "sql.NullString", "String", "\"\""},
	"varbinary":          {"string", "sql.NullString", "String", "\"\""},
}

type goType struct {
	baseType      string
	nullType      string
	nullTypeValue string
	defaultValue  string
}

func getFields(tableName, driverName string, db *sql.DB) ([]Field, []Field, error) {
	var sqlStr = `select
					column_name name,
					data_type type,
					if('YES'=is_nullable,true,false) is_nullable,
					if('PRI'=column_key,true,false) is_primary_key,
					column_comment comment,column_default 'default'
				from
					information_schema.COLUMNS t
				where
					table_schema = DATABASE() `
	sqlStr += fmt.Sprintf(" and t.table_name = '%s' order by is_primary_key desc", tableName)

	if driverName == "postgres" {
		sqlStr = `SELECT 
						t.column_name name ,
						t.udt_name type,
						CASE WHEN t.is_nullable='YES' THEN 1  ELSE 0  END is_nullable,
						CASE WHEN tc.constraint_type='PRIMARY KEY' THEN 1  ELSE 0  END is_primary_key,
						CASE WHEN t.column_comment is null THEN ''	ELSE t.column_comment END  comment,
						t.column_default default
					FROM 
						information_schema.columns t 
						left join information_schema.key_column_usage kcu on kcu.table_name=t.table_name and kcu.column_name=t.column_name
						left join information_schema.table_constraints tc on tc.table_name=kcu.table_name and tc.constraint_name=kcu.constraint_name and tc.constraint_type='PRIMARY KEY'
					WHERE 
						t.table_catalog=current_database() and t.table_schema='public'
		`
		sqlStr += fmt.Sprintf(" and t.table_name = '%s' order by is_primary_key desc", tableName)
	}
	rows, err := db.Query(sqlStr)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fields := make([]Field, 0)
	primaryKeyFields := make([]Field, 0)
	for rows.Next() {
		field := Field{}
		err = rows.Scan(&field.ColumnName, &field.ColumnType, &field.IsNullable, &field.IsPrimaryKey, &field.Comment, &field.ColumnDefault)
		if err != nil {
			panic(err)
		}
		field.FieldName = nstring.ToUpperCamelCase(field.ColumnName)
		field.ColumnNameLowerCamel = nstring.ToLowerCamelCase(field.ColumnName)
		field.ColumnNameUpper = strings.ToUpper(field.ColumnName)
		field.FieldType = dbType[field.ColumnType].baseType
		field.FieldTypeDefault = dbType[field.ColumnType].defaultValue
		field.FieldNullType = dbType[field.ColumnType].nullType
		field.FieldNullTypeValue = dbType[field.ColumnType].nullTypeValue
		field.FieldOrmTag = fmt.Sprintf("orm:\"%v\"", field.ColumnName)
		field.FieldJsonTag = fmt.Sprintf("json:\"%v\"", field.ColumnName)
		field.FieldFormTag = fmt.Sprintf("form:\"%v\"", field.ColumnName)
		val, _ := field.ColumnDefault.Value()
		if val != nil {
			field.FieldDefaultTag = fmt.Sprintf("default:\"%v\"", val)
		}
		if field.IsPrimaryKey == 1 {
			primaryKeyFields = append(primaryKeyFields, field)
		}
		fields = append(fields, field)
	}
	return fields, primaryKeyFields, nil
}

func (gen *Generator) Gen(modules []Module) error {
	if gen.project == "" {
		return errors.New("project can not nil")
	}
	if gen.dsn == "" {
		return errors.New("dsn can not nil")
	}
	if gen.driverName == "" {
		gen.driverName = "mysql"
	}
	db, err := sql.Open(gen.driverName, gen.dsn)
	if err != nil {
		return err
	}
	if modules == nil {
		return errors.New("modules can not nil")
	}

	for _, module := range modules {
		fields, primaryKeyFields, err := getFields(module.TableName, gen.driverName, db)
		if len(fields) == 0 || err != nil {
			fmt.Printf("error:create table %v error=%v", module.TableName, err)
			continue
		}
		// if primaryKeyFields == nil || len(primaryKeyFields) == 0 {
		// 	fmt.Printf("error:table %v no primary key", module.TableName)
		// 	continue
		// }
		tableName := module.TableName
		module.CreateTime = time.Now().Format("2006-01:02 15:04:05.006")
		module.Fields = fields
		module.PrimaryKeyFields = primaryKeyFields
		module.TableNameUpperCamel = nstring.ToUpperCamelCase(tableName)
		module.TableNameLowerCamel = nstring.ToLowerCamelCase(tableName)
		urls := strings.Split(module.ModulePath, gen.project)

		if module.ModelPackageName == "" {
			module.ModelPackageName = "model"
		}

		module.ModelPackagePath = gen.project + urls[1] + "/" + module.ModelPackageName
		module.ModelFileName = tableName + "_" + module.ModelPackageName + ".go"
		module.ModelFilePath = module.ModulePath + "/" + module.ModelPackageName

		module.ExtendPackageName = "extend"
		module.ExtendPackagePath = gen.project + urls[1] + "/" + module.ModelPackageName
		module.ExtendFileName = tableName + "_" + module.ExtendPackageName + ".go"
		module.ExtendFilePath = module.ModulePath + "/" + module.ModelPackageName

		module.ViewPackageName = "view"
		module.ViewPackagePath = gen.project + urls[1] + "/" + module.ViewPackageName
		module.ViewFileName = tableName + "_" + module.ViewPackageName + ".go"
		module.ViewFilePath = module.ModulePath + "/" + module.ViewPackageName

		module.ParamPackageName = "param"
		module.ParamPackagePath = gen.project + urls[1] + "/" + module.ParamPackageName
		module.ParamFileName = tableName + "_" + module.ParamPackageName + ".go"
		module.ParamFilePath = module.ModulePath + "/" + module.ParamPackageName

		module.DaoPackageName = "dao"
		module.DaoPackagePath = gen.project + urls[1] + "/" + module.DaoPackageName
		module.DaoFileName = tableName + "_" + module.DaoPackageName + ".go"
		module.DaoFilePath = module.ModulePath + "/" + module.DaoPackageName

		module.ServicePackageName = "service"
		module.ServicePackagePath = gen.project + urls[1] + "/" + module.ServicePackageName
		module.ServiceFileName = tableName + "_" + module.ServicePackageName + ".go"
		module.ServiceFilePath = module.ModulePath + "/" + module.ServicePackageName

		module.ControllerPackageName = "controller"
		module.ControllerPackagePath = gen.project + urls[1] + "/" + module.ControllerPackageName
		module.ControllerFileName = tableName + "_" + module.ControllerPackageName + ".go"
		module.ControllerFilePath = module.ModulePath + "/" + module.ControllerPackageName

		if module.Model {
			genFile(&module, module.ModelPackageName)
		}
		if module.Extend {
			genFile(&module, module.ExtendPackageName)
		}
		if module.View {
			genFile(&module, module.ViewPackageName)
		}
		if module.Param {
			genFile(&module, module.ParamPackageName)
		}
		if module.Dao {
			genFile(&module, module.DaoPackageName)
		}
		if module.Service {
			genFile(&module, module.ServicePackageName)
		}
		if module.Controller {
			genFile(&module, module.ControllerPackageName)
		}

	}
	return nil
}

func genFile(table *Module, packageName string) {

	var templateStr, filePath, file string
	switch packageName {
	case "model":
		templateStr = getModelTemplate()
		filePath = table.ModelFilePath
		file = filePath + "/" + table.ModelFileName
	case "extend":
		templateStr = getExtendTemplate()
		filePath = table.ExtendFilePath
		file = filePath + "/" + table.ExtendFileName
		if nfile.IsExist(file) { //view 不覆盖
			return
		}
	case "view":
		templateStr = getViewTemplate()
		filePath = table.ViewFilePath
		file = filePath + "/" + table.ViewFileName
		if nfile.IsExist(file) { //view 不覆盖
			return
		}
	case "param":
		templateStr = getParamTemplate()
		filePath = table.ParamFilePath
		file = filePath + "/" + table.ParamFileName
		if nfile.IsExist(file) { //param 不覆盖
			return
		}
	case "dao":
		templateStr = getDaoTemplate()
		filePath = table.DaoFilePath
		file = filePath + "/" + table.DaoFileName
	case "service":
		templateStr = getServiceTemplate()
		filePath = table.ServiceFilePath
		file = filePath + "/" + table.ServiceFileName
		if nfile.IsExist(file) { //service 不覆盖
			return
		}
	case "controller":
		templateStr = getController()
		filePath = table.ControllerFilePath
		file = filePath + "/" + table.ControllerFileName
		if nfile.IsExist(file) { //controller 不覆盖
			return
		}
	}
	// 第一步，加载模版文件
	tmpl, err := template.New("tmpl").Parse(templateStr)
	if err != nil {
		fmt.Println("create template model, err:", err)
		return
	}
	// 第二步，创建文件目录
	err = nfile.CreateDir(filePath)
	if err != nil {
		fmt.Printf("create path:%v err", filePath)
		return
	}
	// 第三步，创建且打开文件
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Can not write file")
		return
	}
	defer f.Close()

	// 第四步，写入数据
	tmpl.Execute(f, table)

	//第五步，格式化代码
	cmd := exec.Command("gofmt", "-w", file)
	cmd.Run()
}

func getModelTemplate() string {
	return `// Create by code north  {{.CreateTime}}
	package model
	
	import (
		"database/sql"
	)
	
	const (
		{{range $field := .Fields}}
			{{- .ColumnNameUpper -}}  ="{{ .ColumnName }}" // {{ .Comment }}
		{{end}}
		TABLE_NAME  = "{{ .TableName }}" // 表名
	)
	
	type {{.TableNameUpperCamel}}Model struct {
		{{range $field := .Fields}}{{ .FieldName }}  {{ .FieldNullType }} ` + "`{{ .FieldOrmTag }} {{ .FieldDefaultTag }}`" + ` // {{ .Comment }}
		{{end}}
	}
	
	func (m *{{.TableNameUpperCamel}}Model) ToMap(includeEmpty bool) map[string]any {
		view := make(map[string]any)
		{{range $field := .Fields}}
			if m.{{ .FieldName }}.Valid {
				view[{{- .ColumnNameUpper -}}] = m.{{ .FieldName }}.{{ .FieldNullTypeValue}}
			} else if includeEmpty {
				view[{{- .ColumnNameUpper -}}] = nil
			}
		{{end}}
		return view
	}`
}
func getExtendTemplate() string {
	return `package model

			type {{.TableNameUpperCamel}}Extend struct {
				{{.TableNameUpperCamel}}Model
			}`
}

func getViewTemplate() string {
	return ` // Create by code north  {{.CreateTime}}
	package view
	
	import (
		"{{.ModelPackagePath}}"
		"time"
	)
	type {{.TableNameUpperCamel}}View struct {
		{{range $field := .Fields}}{{ .FieldName }}  {{ .FieldType }} ` + "`{{ .FieldJsonTag }}`" + ` // {{ .Comment }}
		{{end}}
	}
	func Convert(m model.{{.TableNameUpperCamel}}Model) {{.TableNameUpperCamel}}View {
		return {{.TableNameUpperCamel}}View{
			{{range $field := .Fields}}{{ .FieldName }} : m.{{ .FieldName }}.{{ .FieldNullTypeValue}},
			{{end}}
		}
	}
	func Converts(models []model.{{.TableNameUpperCamel}}Model) []{{.TableNameUpperCamel}}View {
		views := make([]{{.TableNameUpperCamel}}View, 0, len(models))
		for _, model := range models {
			views = append(views, Convert(model))
		}
		return views
	}
	
	func ConvertExtend(m model.{{.TableNameUpperCamel}}Extend) {{.TableNameUpperCamel}}View {
		view := Convert(m.{{.TableNameUpperCamel}}Model)
		return view
	}
	func ConvertExtends(extends []model.{{.TableNameUpperCamel}}Extend) []{{.TableNameUpperCamel}}View {
		views := make([]{{.TableNameUpperCamel}}View, 0, len(extends))
		for _, extend := range extends {
			views = append(views, ConvertExtend(extend))
		}
		return views
	}`
}
func getParamTemplate() string {
	return `// Create by code north  {{.CreateTime}}
			package param
			
			import (
				"time"
			)
			type {{.TableNameUpperCamel}}Param struct {
				{{range $field := .Fields}}{{ .FieldName }}  {{ .FieldType }} ` + "`{{.FieldFormTag}} {{ .FieldJsonTag }}`" + ` // {{ .Comment }}
				{{end}}
				PageNum 	int ` + "`form:\"page\" json:\"page\"`" + `
				PageStart 	int ` + "`form:\"start\" json:\"start\"`" + `
				PageSize 	int ` + "`form:\"size\" json:\"size\"`" + `
			}`
}
func getDaoTemplate() string {
	return `// Create by go-ormerator  {{.CreateTime}}
		  package dao
		  
		  import (
		  	 "database/sql"
			 "github.com/go-lazyer/north"
			 "github.com/go-lazyer/north/nsql"
			 "{{.ModelPackagePath}}"
			 "github.com/pkg/errors"
		  )

		  func CountByOrm(orm  *nsql.CountOrm) (int64, error) {
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return 0,errors.WithStack(err)
			 }
			 return CountBySql(sqlStr, params)
		  }
		  func CountBySql(sqlStr string, params []any) (int64, error) {
			 ds, err := database.DataSource()
			 if err != nil {
				return 0, errors.WithStack(err)
			 }
			 count, err := north.Count(sqlStr, params, ds)
			 if err != nil {
				return 0,errors.WithStack(err)
			 }
			 return count,nil
		  }

		  {{ if gt (len .PrimaryKeyFields) 0 -}}
		  func QuerySingleByPrimaryKey({{range $i,$field := .PrimaryKeyFields}} {{if ne $i 0}},{{end}}{{ .ColumnNameLowerCamel }} any  {{end}}) (model.{{.TableNameUpperCamel}}Model, error) {
			 {{ if eq (len .PrimaryKeyFields) 1 -}} 
			 query := nsql.NewEqualQuery(model.{{(index .PrimaryKeyFields 0).ColumnNameUpper}}, {{(index .PrimaryKeyFields 0).ColumnNameLowerCamel}})
			 {{ else -}}
			 query := nsql.NewBoolQuery(){{range $field := .PrimaryKeyFields}} .And(nsql.NewEqualQuery(model.{{ .ColumnNameUpper }}, {{ .ColumnNameLowerCamel }})) {{end}}
			 {{end}}
			 orm := nsql.NewSelectOrm().Table(model.TABLE_NAME).Where(query)
			 return QuerySingleByOrm(orm)
		  }
		  {{ end -}}
		  func QuerySingleByOrm(orm  *nsql.SelectOrm) (model.{{.TableNameUpperCamel}}Model, error) {
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return model.{{.TableNameUpperCamel}}Model{},errors.WithStack(err)
			 }
			 return QuerySingleBySql(sqlStr, params)
		  }
		  func QuerySingleBySql(sqlStr string, params []any) (model.{{.TableNameUpperCamel}}Model, error) {
			 models, err := QueryBySql(sqlStr, params)
 
			 if len(models) == 0 || err != nil {
				return model.{{.TableNameUpperCamel}}Model{}, errors.WithStack(err)
			 }
			 return models[0], nil
		  }
		  {{if eq (len .PrimaryKeyFields) 1}} 
		  func QueryMapByPrimaryKeys(primaryKeys []any) (map[{{(index .PrimaryKeyFields 0).FieldType}}]model.{{.TableNameUpperCamel}}Model, error) {
			 orm := nsql.NewSelectOrm().Table(model.TABLE_NAME).Where(nsql.NewInQuery(model.{{(index .PrimaryKeyFields 0).ColumnNameUpper}}, primaryKeys))
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
			 return QueryMapBySql(sqlStr, params)
		  }
		  
		  func QueryMapByOrm(orm  *nsql.SelectOrm) (map[{{(index .PrimaryKeyFields 0).FieldType}}]model.{{.TableNameUpperCamel}}Model, error) {
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return nil, errors.WithStack(err)
			 }
			 return QueryMapBySql(sqlStr, params)
		  }
		  
		  func QueryMapBySql(sqlStr string, params []any) (map[{{(index .PrimaryKeyFields 0).FieldType}}]model.{{.TableNameUpperCamel}}Model, error) {
			 {{.TableNameLowerCamel}}s,err := QueryBySql(sqlStr, params)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
 
			 if len({{.TableNameLowerCamel}}s) == 0 {
				return nil,nil
			 }
			 {{.TableNameLowerCamel}}Map := make(map[{{(index .PrimaryKeyFields 0).FieldType}}]model.{{.TableNameUpperCamel}}Model, len({{.TableNameLowerCamel}}s))
			 for _, {{.TableNameLowerCamel}} := range {{.TableNameLowerCamel}}s {
				new := {{.TableNameLowerCamel}}
				{{.TableNameLowerCamel}}Map[{{.TableNameLowerCamel}}.{{(index .PrimaryKeyFields 0).FieldName}}.{{(index .PrimaryKeyFields 0).FieldNullTypeValue}}] = new
			 }
			 return {{.TableNameLowerCamel}}Map,nil
		  }
		  {{end}}
 
		  func QueryByOrm(orm  *nsql.SelectOrm) ([]model.{{.TableNameUpperCamel}}Model, error) {
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
			 return QueryBySql(sqlStr, params)
		  }
		  func QueryBySql(sqlStr string, params []any) ([]model.{{.TableNameUpperCamel}}Model, error) {
			 ds, err := database.DataSource()
			 if err != nil {
				return nil, errors.WithStack(err)
			 }
			 {{.TableNameLowerCamel}}s, err := north.Query[model.{{.TableNameUpperCamel}}Model](sqlStr, params, ds)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
			 return {{.TableNameLowerCamel}}s,nil
		  }
 
 
		  func QueryExtendByOrm(orm  *nsql.SelectOrm) ([]model.{{.TableNameUpperCamel}}Extend, error) {
			 sqlStr, params, err := orm.ToSql(true)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
			 return QueryExtendBySql(sqlStr, params)
		  }
		  func QueryExtendBySql(sqlStr string, params []any) ([]model.{{.TableNameUpperCamel}}Extend, error) {
			 ds, err := database.DataSource()
			 if err != nil {
				return nil, errors.WithStack(err)
			 }
			 {{.TableNameLowerCamel}}Extends, err := north.Query[model.{{.TableNameUpperCamel}}Extend](sqlStr, params, ds)
			 if err != nil {
				return nil,errors.WithStack(err)
			 }
			 return {{.TableNameLowerCamel}}Extends,nil
		  }

		func Insert(m model.{{.TableNameUpperCamel}}Model) (int64, error) {
			return InsertTx(m,nil)
		}
		
		func InsertByMap(insertMap map[string]any) (int64, error) {
			return InsertByMapTx(insertMap,nil)
		}
	
		func InsertByOrm(orm  *nsql.InsertOrm) (int64, error) {
			return InsertByOrmTx(orm,nil)
		}
			
		func InsertBySql(sqlStr string, params []any) (int64, error) {
			return InsertBySqlTx(sqlStr,params,nil)
		}

		func InsertTx(m model.{{.TableNameUpperCamel}}Model, tx *sql.Tx) (int64, error) {
			return InsertByMapTx(m.ToMap(false), tx)
		}
		
		func InsertByMapTx(insertMap map[string]any, tx *sql.Tx) (int64, error) {
			orm := nsql.NewInsertOrm().Table(model.TABLE_NAME).Insert(insertMap)
			return InsertByOrmTx(orm,tx)
		}
	
		func InsertByOrmTx(orm  *nsql.InsertOrm, tx *sql.Tx) (int64, error) {
			sqlStr, params, err := orm.ToSql(true)
			if err != nil {
			return 0, errors.WithStack(err)
			}
			return InsertBySqlTx(sqlStr, params,tx)
		}
			
		func InsertBySqlTx(sqlStr string, params []any, tx *sql.Tx) (int64, error) {
			ds, err := database.DataSource()
			if err != nil {
			return 0, errors.WithStack(err)
			}
			if tx != nil {
			ds.Tx = tx
			}
			id, err := north.Insert(sqlStr, params,ds)
			if err != nil {
			return 0, errors.WithStack(err)
			}
			return id, nil
		}

		{{ if gt (len .PrimaryKeyFields) 0 -}}
		func Update(m model.{{.TableNameUpperCamel}}Model) (int64, error) {
			return UpdateTx(m,nil)
		}
		{{end}}

		func UpdateByOrm(orm  *nsql.UpdateOrm) (int64, error) {
			return UpdateByOrmTx(orm,nil)
		}
		func UpdateBySql(sqlStr string, params []any) (int64, error) {
			 return UpdateBySqlTx(sqlStr, params,nil)
		}

		{{ if gt (len .PrimaryKeyFields) 0 -}}
		func UpdateTx(m model.{{.TableNameUpperCamel}}Model, tx *sql.Tx) (int64, error) {
			{{ if eq (len .PrimaryKeyFields) 1 -}} 
			query := nsql.NewEqualQuery(model.{{(index .PrimaryKeyFields 0).ColumnNameUpper}}, m.{{(index .PrimaryKeyFields 0).FieldName}}.{{(index .PrimaryKeyFields 0).FieldNullTypeValue}})
			{{ else -}}
			query := nsql.NewBoolQuery(){{range $field := .PrimaryKeyFields}} .And(nsql.NewEqualQuery(model.{{ .ColumnNameUpper }}, m.{{.FieldName}}.{{.FieldNullTypeValue}})) {{end}}
			{{end -}}
			orm := nsql.NewUpdateOrm().Table(model.TABLE_NAME).Update(m.ToMap(false)).Where(query)
			return UpdateByOrmTx(orm,tx)
		}
		{{end}}

		func UpdateByOrmTx(orm  *nsql.UpdateOrm, tx *sql.Tx) (int64, error) {
			sqlStr, params, err := orm.ToSql(true)
			if err != nil {
			   return 0, errors.WithStack(err)
			}
			return UpdateBySqlTx(sqlStr, params,tx)
		 }
		func UpdateBySqlTx(sqlStr string, params []any, tx *sql.Tx) (int64, error) {
			 ds, err := database.DataSource()
			 if err != nil {
				return 0, errors.WithStack(err)
			 }
			 if tx != nil {
				ds.Tx = tx
			 }
			 count, err := north.Update(sqlStr, params,ds)
			 if err != nil {
				return 0, errors.WithStack(err)
			 }
			 return count, nil
		}

		{{ if gt (len .PrimaryKeyFields) 0 -}}
		func DeleteByPrimaryKey({{range $i,$field := .PrimaryKeyFields}} {{if ne $i 0}},{{end}}{{ .ColumnNameLowerCamel }} any  {{end}}) (int64, error) {
			return DeleteByPrimaryKeyTx({{range $i,$field := .PrimaryKeyFields}} {{ .ColumnNameLowerCamel }},{{end}} nil)
		}
		{{ end -}}
		{{ if eq (len .PrimaryKeyFields) 1 -}}
		func DeleteByPrimaryKeys(primaryKeys []any) (int64, error) {
			return DeleteByPrimaryKeysTx(primaryKeys,nil)
		}
		{{ end -}}
		func DeleteByOrm(orm  *nsql.DeleteOrm) (int64, error) {
			return DeleteByOrmTx(orm, nil)
		}
		func DeleteBySql(sqlStr string, params []any) (int64, error) {
			return DeleteBySqlTx( sqlStr, params,nil)
		}
		
		{{ if gt (len .PrimaryKeyFields) 0 -}}
		func DeleteByPrimaryKeyTx({{range $i,$field := .PrimaryKeyFields}} {{if ne $i 0}},{{end}}{{ .ColumnNameLowerCamel }} any  {{end}}, tx *sql.Tx) (int64, error) {
			{{ if eq (len .PrimaryKeyFields) 1 -}} 
			orm := nsql.NewDeleteOrm().Table(model.TABLE_NAME).Where(nsql.NewEqualQuery(model.{{(index .PrimaryKeyFields 0).ColumnNameUpper}}, {{(index .PrimaryKeyFields 0).ColumnNameLowerCamel}}))
			{{ else -}}
			query := nsql.NewBoolQuery(){{range $field := .PrimaryKeyFields}} .And(nsql.NewEqualQuery(model.{{ .ColumnNameUpper }}, {{ .ColumnNameLowerCamel }})) {{end}}
			orm := nsql.NewDeleteOrm().Table(model.TABLE_NAME).Where(query)
			{{ end -}}
			return DeleteByOrmTx(orm,tx)
		}
		{{ end -}}
		{{ if eq (len .PrimaryKeyFields) 1 -}}
		func DeleteByPrimaryKeysTx(primaryKeys []any, tx *sql.Tx) (int64, error) {
			orm := nsql.NewDeleteOrm().Table(model.TABLE_NAME).Where(nsql.NewInQuery(model.{{(index .PrimaryKeyFields 0).ColumnNameUpper}}, primaryKeys))
			return DeleteByOrmTx(orm,tx)
		}
		{{ end -}}
		func DeleteByOrmTx(orm  *nsql.DeleteOrm, tx *sql.Tx) (int64, error) {
			sqlStr, params, err := orm.ToSql(true)
			if err != nil {
			return 0, errors.WithStack(err)
			}
			return DeleteBySqlTx(sqlStr, params,tx)
		}
		func DeleteBySqlTx(sqlStr string, params []any, tx *sql.Tx) (int64, error) {
			ds, err := database.DataSource()
			if err != nil {
			return 0, errors.WithStack(err)
			}
			if tx != nil {
			ds.Tx = tx
			}
			count, err := north.Delete(sqlStr, params,ds)
			if err != nil {
			return 0, errors.WithStack(err)
			}
			return count, nil
		}
		func GetDataSource() (north.DataSource, error) {
			return database.DataSource()
		}`
}
func getServiceTemplate() string {
	return `// Create by code north  {{.CreateTime}}
			package service

			import (
			
				"{{.DaoPackagePath}}"
				"{{.ModelPackagePath}}"
				"{{.ParamPackagePath}}"
			
				norm "github.com/go-lazyer/north/orm"
			)

			{{ if gt (len .PrimaryKeyFields) 0 -}} 
			func QuerySingleByPrimaryKey({{range $i,$field := .PrimaryKeyFields}} {{if ne $i 0}},{{end}}{{ .ColumnNameLowerCamel }} any  {{end}}) (model.{{.TableNameUpperCamel}}Model, error) {
				{{.TableNameLowerCamel}}, err := dao.QuerySingleByPrimaryKey({{range $i,$field := .PrimaryKeyFields}} {{if ne $i 0}},{{end}}{{ .ColumnNameLowerCamel }}   {{end}})
				if err != nil {
					return model.{{.TableNameUpperCamel}}Model{},err
				}
				return {{.TableNameLowerCamel}},nil
			}
			{{end}}

			func QueryByParam({{.TableNameLowerCamel}}Param param.{{.TableNameUpperCamel}}Param) ([]model.{{.TableNameUpperCamel}}Model, error) {
				query := nsql.NewBoolQuery()
				orm := nsql.NewSelectOrm().PageNum({{.TableNameLowerCamel}}Param.PageNum).PageStart({{.TableNameLowerCamel}}Param.PageStart).PageSize({{.TableNameLowerCamel}}Param.PageSize).Table(model.TABLE_NAME).Where(query)
				{{.TableNameLowerCamel}}s, err := dao.QueryByOrm(orm)
				if err != nil {
					return nil,err
				}
				return {{.TableNameLowerCamel}}s,nil
			}`
}
func getController() string {
	return `// Create by code north  {{.CreateTime}}
			package controller
			
			import (
				"net/http"
			
				"github.com/gin-gonic/gin"
			)
			
			func Index(g *gin.Context) {
				data := gin.H{
					"code": 200,
				}
				g.JSON(http.StatusOK, data)
			}`
}
