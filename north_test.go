package north

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-lazyer/north/constant"
	"github.com/go-lazyer/north/nmap"
	"github.com/go-lazyer/north/nsql"
	_ "github.com/go-sql-driver/mysql"
)

func GetDataSource() (DataSource, error) {
	username := "root"
	password := "Daoway_Mysql_iO12"
	host := "test.daoway.cn"
	port := "3306"
	dbname := "north"
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 创建数据库连接
	return Open(constant.DRIVER_NAME_MYSQL, connStr, Config{MaxOpenConns: 10, MaxIdleConns: 10})
}

func TestCount(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	orm := nsql.NewCountOrm().Table("user")
	sqlStr, prarm, _ := orm.ToSql(true)
	fmt.Println(Count(sqlStr, prarm, ds))
}

func TestPrepareCount(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	query := nsql.NewBetweenQuery("create_time", "2023-01-01 00:00:00", "2023-12-31 23:59:59")
	orm := nsql.NewCountOrm().Table("user").Where(query)
	sqlStr, prarms, _ := orm.ToPrepareSql()
	params := make([][]any, 0)
	params = append(params, []any{"2024-01-01 00:00:00", "2024-12-31 23:59:59"})
	params = append(params, []any{"2025-01-01 00:00:00", "2025-12-31 23:59:59"})
	res, _ := PrepareCount(sqlStr, prarms, ds)
	fmt.Println(res)
}

type User struct {
	Id         sql.NullString `orm:"id" `             //
	Name       sql.NullString `orm:"name" default:""` // 昵称
	Age        sql.NullInt64  `orm:"age" default:""`  // 年龄
	Sex        sql.NullString `orm:"sex" default:""`  // 性别
	CreateTime sql.NullTime   `orm:"create_time" `    // 用户名
}

func TestQuery(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	query := nsql.NewBetweenQuery("create_time", "2023-01-01 00:00:00", "2023-12-31 23:59:59")
	orm := nsql.NewSelectOrm().Result("id", "nickname", "username", "create_time").Table("user").Where(query).PageSize(2)
	sqlStr, prarm, _ := orm.ToSql(true)
	fmt.Println(Query[User](sqlStr, prarm, ds))
}

func TestPrepareQuery(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	query1 := nsql.NewBetweenQuery("create_time", "2023-01-01 00:00:00", "2023-12-31 23:59:59")
	query2 := nsql.NewBetweenQuery("create_time", "2024-01-01 00:00:00", "2024-12-31 23:59:59")
	orm := nsql.NewSelectOrm().Table("user").Where(query1, query2).PageSize(1)
	sqlStr, prarms, _ := orm.ToPrepareSql()
	res, _ := PrepareQuery[User](sqlStr, prarms, ds)
	fmt.Println(res)
}

func TestInsert(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	m := map[string]any{
		"name":        "lilie",
		"sex":         "boy",
		"age":         "10",
		"create_time": "2023-01-01 00:00:00",
	}

	orm := nsql.NewInsertOrm().Table("user").Insert(m)
	sqlStr, prarm, _ := orm.ToSql(true)
	fmt.Println(Insert(sqlStr, prarm, ds))
}

func TestInserts(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	m := map[string]any{
		"name":        "lili",
		"sex":         "girl",
		"age":         "20",
		"create_time": "2023-01-01 00:00:00",
	}

	m1 := map[string]any{
		"name":        "lilie",
		"sex":         "boy",
		"age":         "10",
		"create_time": "2023-01-01 00:00:00",
	}

	orm := nsql.NewInsertOrm().Table("user").Insert(m, m1)
	sqlStr, prarm, _ := orm.ToSql(true)
	fmt.Println(Insert(sqlStr, prarm, ds))
}
func TestPrepareInsert(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	m := map[string]any{
		"sex":         "girl",
		"name":        "lili",
		"age":         "20",
		"create_time": "2023-01-01 00:00:00",
	}

	m1 := map[string]any{
		"name":        "lilie",
		"sex":         "boy",
		"age":         "10",
		"create_time": "2023-01-01 00:00:00",
	}

	orm := nsql.NewInsertOrm().Table("user").Insert(m, m1)
	sqlStr, prarm, _ := orm.ToPrepareSql()

	fmt.Println(sqlStr, prarm, err)

	fmt.Println(PrepareInsert(sqlStr, prarm, ds))
}

func TestUpdate(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	m := map[string]any{
		"name": "hanmeimei",
		"sex":  "girl",
		"age":  "20",
	}

	orm := nsql.NewUpdateOrm().Table("user").Update(m).Where(nsql.NewEqualQuery("id", "1"))
	sqlStr, prarm, err := orm.ToSql(true)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(Update(sqlStr, prarm, ds))
}

func TestUpdates(t *testing.T) {

	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	f1 := map[string]any{
		"sex":        "boy",
		"age":        "20",
		"_condition": nsql.NewEqualQuery("id", "1"),
	}
	f2 := map[string]any{
		"sex":        "girl",
		"age":        "40",
		"_condition": nsql.NewEqualQuery("name", "lilie"),
	}

	// query := nsql.NewEqualQuery("create_time", "2025-01-01 00:00:00")
	query := nsql.NewNullQuery("create_time")
	orm := nsql.NewUpdateOrm().Table("user").Where(query).Update(f1, f2)
	sqlStr, prarm, err := orm.ToSql(true)
	fmt.Println(sqlStr, prarm)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(Update(sqlStr, prarm, ds))
}

func TestPrepareUpdates(t *testing.T) {

	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}

	f1 := map[string]any{
		"age": "20",
		"sex": "boy",
	}
	f2 := map[string]any{
		"age":  "40",
		"sex1": "girl",
	}

	query1 := nsql.NewEqualQuery("create_time", "2023-01-01 00:00:00")
	query2 := nsql.NewEqualQuery("create_time", "2024-01-01 00:00:00")
	// query3 := nsql.NewEqualQuery("create_time", "2025-01-01 00:00:00")
	// q := nsql.NewNotNullQuery("create_time")
	orm := nsql.NewUpdateOrm().Table("user").Where(query1, query2).Update(f1, f2)
	sqlStr, prarms, err := orm.ToPrepareSql()
	fmt.Println(sqlStr, prarms, err)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(PrepareUpdate(sqlStr, prarms, ds))
}

func Test_Main(t *testing.T) {
	f2 := map[string]any{
		"age":    "40",
		"sex1":   "girl",
		"name":   "girl",
		"create": "girl",
	}
	fields := nmap.Keys(f2)

	for i := 0; i < 10; i++ {
		for _, field := range fields {
			fmt.Print(field)
		}
		fmt.Println("")
	}

}
func TestDelete(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	query := nsql.NewBetweenQuery("create_time", "2023-01-01 00:00:00", "2025-12-31 23:59:59")
	orm := nsql.NewDeleteOrm().Table("user").Where(query).PageSize(2)
	sqlStr, param, _ := orm.ToSql(true)
	fmt.Print(sqlStr, param)
	fmt.Println(Delete(sqlStr, param, ds))
}

func TestPrepareDelete(t *testing.T) {
	ds, err := GetDataSource()
	if err != nil {
		fmt.Print(err)
	}
	query1 := nsql.NewBetweenQuery("create_time", "2023-01-01 00:00:00", "2023-12-31 23:59:59")
	query2 := nsql.NewBetweenQuery("create_time", "2024-01-01 00:00:00", "2024-12-31 23:59:59")
	orm := nsql.NewDeleteOrm().Table("user").Where(query1, query2).PageSize(1)
	sqlStr, prarms, _ := orm.ToPrepareSql()

	fmt.Print(sqlStr, prarms)

	res, err := PrepareDelete(sqlStr, prarms, ds)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
