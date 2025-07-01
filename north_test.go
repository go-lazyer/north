package north

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	username := "root"
	password := "******"
	host := "test.daoway.cn"
	port := "3306"
	dbname := "daowei"
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 创建数据库连接
	ds, err := Open(DRIVER_NAME_MYSQL, connStr, Config{MaxOpenConns: 10, MaxIdleConns: 10})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ds)
	// sql1 := "select * from test t where t.id='a'"
	// params := make([]any, 0)
	// ts, err := Query[Test1](sql1, params, ds)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, testInstance := range ts {
	// 	fmt.Print(testInstance)
	// 	// 如果有更多字段，可以继续访问它们
	// }
}
