package nsql

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/go-lazyer/north/constant"
	"github.com/go-lazyer/north/nmap"
)

const (
	INNER_JOIN = "inner join" // inner  join
	LEFT_JOIN  = "left join"  // left  join
	RIGHT_JOIN = "right join" // right join
)

type Orm struct {
	orderBy    []string //排序字段
	groupBy    []string //分组字段
	pageStart  int
	pageSize   int
	pageNum    int
	querys     []BaseQuery
	update     []map[string]any
	insert     []map[string]any
	joins      []*Join
	tableName  string
	tableAlias string
	columns    []string
	action     string //操作类型 count,select,update,delete,insert
}

func NewCountOrm() *Orm {
	return &Orm{action: "count"}
}
func NewSelectOrm() *Orm {
	return &Orm{action: "select"}
}

func NewUpdateOrm() *Orm {
	return &Orm{action: "update"}
}
func NewDeleteOrm() *Orm {
	return &Orm{action: "delete"}
}
func NewInsertOrm() *Orm {
	return &Orm{action: "insert"}
}

func (s *Orm) Where(query ...BaseQuery) *Orm {
	if s.querys == nil {
		s.querys = make([]BaseQuery, 0)
	}
	s.querys = append(s.querys, query...)
	return s
}

func (s *Orm) Update(m ...map[string]any) *Orm {
	s.update = m
	return s
}
func (s *Orm) Insert(m ...map[string]any) *Orm {
	s.insert = m
	return s
}
func (s *Orm) Join(join ...*Join) *Orm {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins, join...)
	return s
}
func (s *Orm) Table(tableName string) *Orm {
	s.tableName = tableName
	return s
}

// 表的别名
func (s *Orm) TableAlias(tableName, tableAlias string) *Orm {
	s.tableName = tableName
	s.tableAlias = tableAlias
	return s
}
func (s *Orm) Result(columns ...string) *Orm {
	s.columns = columns
	return s
}
func (s *Orm) PageNum(pageNum int) *Orm {
	s.pageNum = pageNum
	return s
}
func (s *Orm) PageStart(pageStart int) *Orm {
	s.pageStart = pageStart
	return s
}
func (s *Orm) PageSize(pageSize int) *Orm {
	s.pageSize = pageSize
	return s
}
func (s *Orm) OrderBy(orderBy []string) *Orm {
	s.orderBy = orderBy
	return s
}
func (s *Orm) AddOrderBy(name string, orderByType string) *Orm {
	if s.orderBy == nil {
		s.orderBy = make([]string, 0)
	}
	s.orderBy = append(s.orderBy, name+" "+orderByType)
	return s
}

func (s *Orm) GroupBy(groupBy []string) *Orm {
	s.groupBy = groupBy
	return s
}
func (s *Orm) AddGroupBy(tableName, name string) *Orm {
	if s.groupBy == nil {
		s.groupBy = make([]string, 0)
	}
	s.groupBy = append(s.groupBy, tableName+"."+name)
	return s
}

func (s *Orm) countSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	params := make([]any, 0)
	var sql bytes.Buffer
	sql.WriteString("select ")

	result := " count(*) count  "
	if len(s.columns) == 1 {
		result = strings.Join(s.columns, ",")
	}

	sql.WriteString(result)
	sql.WriteString(" from  " + s.tableName + "")

	if s.tableAlias != "" {
		sql.WriteString(" " + s.tableAlias + " ")
	}

	for _, join := range s.joins {
		sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
		for i, query := range join.querys {
			if i == 0 {
				sql.WriteString(" and ")
			} else {
				sql.WriteString(" or ")
			}
			source, param, _ := query.Source(join.tableName, prepare)
			sql.WriteString(" " + source + " ")
			params = append(params, param...)
		}
	}

	for n, query := range s.querys {
		source, param, err := query.Source(table, prepare)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   ")
		} else {
			sql.WriteString(" or ")
		}
		sql.WriteString(" " + source + " ")
		params = append(params, param...)
	}

	return sql.String(), params, nil
}

func (s *Orm) countPrepareSql() (string, [][]any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}

	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	var sql bytes.Buffer
	sql.WriteString("select  count(*) count  ")

	result := ""
	if len(s.columns) == 1 {
		result = strings.Join(s.columns, ",")
	}

	sql.WriteString(result)
	sql.WriteString(" from  " + s.tableName + "")

	if s.tableAlias != "" {
		sql.WriteString(" " + s.tableAlias + " ")
	}
	params := make([][]any, len(s.querys))
	for _, join := range s.joins {
		sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
		for i, query := range join.querys {
			if i == 0 {
				sql.WriteString(" and ")
			} else {
				sql.WriteString(" or ")
			}
			source, p, _ := query.Source(join.tableName, true)
			sql.WriteString(" " + source + " ")
			for i, param := range params {
				params[i] = append(param, p...)
			}
		}
	}

	for n, query := range s.querys {
		source, p, err := query.Source(table, true)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   ")
		} else {
			sql.WriteString(" or ")
		}
		sql.WriteString(" " + source + " ")
		params[n] = append(params[n], p...)
	}

	return sql.String(), params, nil
}

func (s *Orm) ToSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	switch s.action {
	case "count":
		return s.countSql(prepare)
	case "select":
		return s.selectSql(prepare)
	case "update":
		return s.updateSql(prepare)
	case "delete":
		return s.deleteSql(prepare)
	case "insert":
		return s.insertSql(prepare)
	default:
		return "", nil, errors.New("action not supported")
	}
}
func (s *Orm) ToPrepareSql() (string, [][]any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	switch s.action {
	case "count":
		return s.countPrepareSql()
	case "select":
		return s.selectPrepareSql()
	case "update":
		return s.updatePrepareSql()
	case "delete":
		return s.deletePrepareSql()
	case "insert":
		return s.insertPrepareSql()
	default:
		return "", nil, errors.New("action not supported")
	}
}

func (s *Orm) selectSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	params := make([]any, 0)
	var sql bytes.Buffer
	sql.WriteString("select ")
	if s.columns == nil {
		sql.WriteString(" * ")
	} else {
		sql.WriteString(strings.Join(s.columns, ","))
	}
	sql.WriteString(" from  " + s.tableName + "")

	if s.tableAlias != "" {
		sql.WriteString(" " + s.tableAlias + " ")
	}

	for _, join := range s.joins {
		table := join.tableName
		if join.tableAlias != "" {
			table = join.tableAlias
		}
		if join.tableName != "" {
			sql.WriteString(fmt.Sprintf(" %v %v %v on %v", join.joinType, join.tableName, join.tableAlias, join.condition))
		} else {
			sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
		}
		for i, query := range join.querys {
			if i == 0 {
				sql.WriteString(" and ")
			} else {
				sql.WriteString(" or ")
			}
			source, param, _ := query.Source(table, prepare)
			sql.WriteString(" " + source + " ")
			params = append(params, param...)
		}
	}

	for n, query := range s.querys {
		if query == nil {
			continue
		}
		source, param, err := query.Source(table, prepare)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   ")
		} else {
			sql.WriteString(" or ")
		}
		sql.WriteString(" " + source + " ")
		params = append(params, param...)
	}

	for n, v := range s.groupBy {
		if n == 0 {
			sql.WriteString(" group by   ")
		} else {
			sql.WriteString(", ")
		}
		sql.WriteString(v)
	}
	for n, v := range s.orderBy {
		if n == 0 {
			sql.WriteString(" order by   ")
		} else {
			sql.WriteString(", ")

		}
		sql.WriteString(v)
	}
	if s.pageSize > 0 {
		if s.pageNum > 0 {
			s.pageStart = (s.pageNum - 1) * s.pageSize
		}
		params = append(params, s.pageSize, s.pageStart)
		if prepare {
			sql.WriteString(fmt.Sprintf(" limit %s offset %s", constant.PLACE_HOLDER_GO, constant.PLACE_HOLDER_GO))
		} else {
			sql.WriteString(fmt.Sprintf(" limit %d offset %d", s.pageSize, s.pageStart))
		}
	}

	return sql.String(), params, nil
}

func (s *Orm) selectPrepareSql() (string, [][]any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}

	//设置表别名
	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	var sql bytes.Buffer
	sql.WriteString("select ")
	if s.columns == nil {
		sql.WriteString(" * ")
	} else {
		sql.WriteString(strings.Join(s.columns, ","))
	}
	sql.WriteString(" from  " + s.tableName + "")

	if s.tableAlias != "" {
		sql.WriteString(" " + s.tableAlias + " ")
	}
	params := make([][]any, len(s.querys))
	for _, join := range s.joins { //forr 会自动判断joins不为空
		table := join.tableName
		if join.tableAlias != "" {
			table = join.tableAlias
		}
		if join.tableName != "" {
			sql.WriteString(fmt.Sprintf(" %v %v %v on %v", join.joinType, join.tableName, join.tableAlias, join.condition))
		} else {
			sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
		}
		for i, query := range join.querys {
			if i == 0 {
				sql.WriteString(" and ")
			} else {
				sql.WriteString(" or ")
			}
			source, p, _ := query.Source(table, true)
			sql.WriteString(" " + source + " ")

			for i, param := range params {
				params[i] = append(param, p...)
			}
		}
	}

	// 添加查询
	for n, query := range s.querys { //forr 会自动判断 querys 不为空
		if query == nil {
			continue
		}
		source, p, err := query.Source(table, true)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   " + source + " ")
		}
		params[n] = append(params[n], p...)
	}

	// 添加分组
	for n, v := range s.groupBy {
		if n == 0 {
			sql.WriteString(" group by   ")
		} else {
			sql.WriteString(", ")
		}
		sql.WriteString(v)
	}

	// 添加排序
	for n, v := range s.orderBy {
		if n != 0 {
			sql.WriteString(" order by   ")
		} else {
			sql.WriteString(", ")

		}
		sql.WriteString(v)
	}
	if s.pageSize > 0 {
		if s.pageNum > 0 {
			s.pageStart = (s.pageNum - 1) * s.pageSize
		}
		sql.WriteString(fmt.Sprintf(" limit %s offset %s", constant.PLACE_HOLDER_GO, constant.PLACE_HOLDER_GO))
		for i, param := range params {
			params[i] = append(param, s.pageSize, s.pageStart)
		}
	}

	return sql.String(), params, nil
}

func (s *Orm) deleteSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	if s.querys == nil || len(s.querys) != 1 {
		return "", nil, errors.New("the querys size must be 1")
	}
	params := make([]any, 0)
	var sql bytes.Buffer
	sql.WriteString("delete from " + s.tableName + " ")

	sql.WriteString(" where   ")
	for i, query := range s.querys {
		if i != 0 {
			sql.WriteString(" or ")
		}
		source, param, _ := query.Source(s.tableName, prepare)
		sql.WriteString(" " + source + " ")
		params = append(params, param...)
	}
	if s.pageSize > 0 {
		sql.WriteString(fmt.Sprintf(" limit %s", constant.PLACE_HOLDER_GO))
		params = append(params, s.pageSize)
	}
	return sql.String(), params, nil
}
func (s *Orm) deletePrepareSql() (string, [][]any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	if len(s.querys) <= 0 {
		return "", nil, errors.New("the querys cannot be empty")
	}
	params := make([][]any, len(s.querys))
	var sql bytes.Buffer
	sql.WriteString("delete from " + s.tableName + " ")

	// 添加查询
	for n, query := range s.querys { //forr 会自动判断 querys 不为空
		if query == nil {
			continue
		}
		source, p, err := query.Source(s.tableName, true)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   " + source + " ")
		}
		params[n] = append(params[n], p...)
	}

	if s.pageSize > 0 {
		sql.WriteString(fmt.Sprintf(" limit %s", constant.PLACE_HOLDER_GO))
		for i, param := range params {
			params[i] = append(param, s.pageSize)
		}
	}
	return sql.String(), params, nil
}
func (s *Orm) insertSql(prepare bool) (string, []any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	if len(s.insert) == 0 {
		return "", nil, errors.New("insert cannot be empty")
	}
	params := make([]any, 0)
	var sql bytes.Buffer
	sql.WriteString("insert into " + s.tableName + " ")
	sql.WriteString("(")
	//把所有要修改的字段提取出来

	fields := nmap.Keys(s.insert[0])

	for n, field := range fields {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString(" " + field + " ")
	}
	sql.WriteString(") values")

	for n, maps := range s.insert {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString("(")
		for m, field := range fields {
			if m != 0 {
				sql.WriteString(",")
			}
			params = append(params, maps[field])
			if prepare {
				sql.WriteString(fmt.Sprintf(" %s ", constant.PLACE_HOLDER_GO))
			} else {
				sql.WriteString(fmt.Sprintf(" '%v' ", maps[field]))
			}
		}
		sql.WriteString(")")
	}

	return sql.String(), params, nil
}

func (s *Orm) insertPrepareSql() (string, [][]any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	if len(s.insert) == 0 {
		return "", nil, errors.New("insert cannot be empty")
	}

	var sql bytes.Buffer
	sql.WriteString("insert into " + s.tableName + " ")
	sql.WriteString("(")
	//把所有要修改的字段提取出来

	fields := nmap.Keys(s.insert[0])

	var valuesSql bytes.Buffer
	for n, field := range fields {
		if n != 0 {
			sql.WriteString(",")
			valuesSql.WriteString(",")
		}
		sql.WriteString(" " + field + " ")
		valuesSql.WriteString(" " + constant.PLACE_HOLDER_GO + " ")
	}
	sql.WriteString(") values (" + valuesSql.String() + ")")
	params := make([][]any, 0, len(s.insert))
	for _, maps := range s.insert {
		param := make([]any, 0)
		for _, field := range fields {
			param = append(param, maps[field])
		}
		params = append(params, param)
	}

	return sql.String(), params, nil
}

func (s *Orm) updateSql(prepare bool) (string, []any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName  cannot be empty")
	}

	if len(s.querys) == 0 {
		return "", nil, errors.New("the querys cannot be empty")
	}

	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	params := make([]any, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("update " + s.tableName + " set ")
	n := 0
	if len(s.update) > 1 { //批量更新
		//把所有要修改的字段提取出来
		fields := make(map[string]string)
		for _, setMap := range s.update {
			for name := range setMap {
				if name == constant.CONDITION {
					continue
				}
				fields[name] = ""
			}
		}

		for field := range fields {
			if n != 0 {
				sql.WriteString(",")
			}
			sql.WriteString(fmt.Sprintf("%v = CASE ", field))
			for _, setMap := range s.update {
				condition, ok := setMap[constant.CONDITION]
				if !ok {
					return "", nil, errors.New("missing condition in update map for batch update")
				}

				baseQuery, ok := condition.(BaseQuery)
				if !ok {
					return "", nil, errors.New("condition in update map is not of type query")
				}

				v, ok := setMap[field]
				if !ok {
					continue
				}

				source, param, err := baseQuery.Source(s.tableName, prepare)
				if err != nil {
					return "", nil, err
				}
				params = append(params, param...)
				sql.WriteString(fmt.Sprintf(" WHEN %v THEN %v", source, constant.PLACE_HOLDER_GO))
				params = append(params, v)
			}
			sql.WriteString(" else " + field)
			sql.WriteString(" END ")
			n++
		}
	} else { //单个更新
		for name, value := range s.update[0] {
			if n != 0 {
				sql.WriteString(",")
			}
			if prepare {
				sql.WriteString(fmt.Sprintf("%v=%s", name, constant.PLACE_HOLDER_GO))
			} else {
				sql.WriteString(fmt.Sprintf("%v='%v'", name, value))
			}
			params = append(params, value)
			n++
		}
	}

	// 添加查询
	for n, query := range s.querys {
		if query == nil {
			continue
		}
		source, param, err := query.Source(table, prepare)
		if err != nil {
			return "", nil, err
		}
		if source == "" {
			continue
		}
		if n == 0 {
			sql.WriteString(" where   ")
		} else {
			sql.WriteString(" or ")
		}
		sql.WriteString(" " + source + " ")
		params = append(params, param...)
	}

	return sql.String(), params, nil
}

// 用于批量更新，单独更新不建议使用
// 生成updatePrepareSql 是orm 中的 update 和 querys 中len的关系为，一对多，多对一，或者n对n 只有这3种情况
// 多个update 的长度和key必须一致
// 多个querys 的长度和query必须一致

func (s *Orm) updatePrepareSql() (string, [][]any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName  cannot be empty")
	}

	if len(s.update) == 0 {
		return "", nil, errors.New("the update cannot be empty")
	}
	if len(s.querys) == 0 {
		return "", nil, errors.New("the querys cannot be empty")
	}

	if len(s.update) > 1 && len(s.querys) > 1 && len(s.update) != len(s.querys) {
		return "", nil, errors.New("the update and querys must be one to one or many to one")
	}

	table := s.tableName
	if s.tableAlias != "" {
		table = s.tableAlias
	}

	params := make([][]any, max(len(s.update), len(s.querys)))

	var sql bytes.Buffer
	sql.WriteString("update " + s.tableName + " set ")

	fields := nmap.Keys(s.update[0])

	for i, update := range s.update {
		if len(update) != len(s.update[0]) {
			return "", nil, fmt.Errorf("update map %d has different length from the first update map", i)
		}
		n := 0
		for _, field := range fields {
			value, ok := update[field]
			if !ok {
				return "", nil, fmt.Errorf("missing value for field %s in update map", field)
			}

			if i == 0 { //多个update ，只把第一个拼写到sql 中
				if n != 0 {
					sql.WriteString(",")
				}
				sql.WriteString(fmt.Sprintf("%v=%s", field, constant.PLACE_HOLDER_GO))
			}

			if len(s.update) == 1 {
				for i, param := range params {
					params[i] = append(param, value)
				}
			} else {
				params[i] = append(params[i], value)
			}
			n = n + 1
		}
	}
	// 添加查询
	source0, _, _ := s.querys[0].Source(table, true)
	for n, query := range s.querys { //forr 会自动判断 querys 不为空
		if query == nil {
			continue
		}
		source, p, err := query.Source(table, true)
		if err != nil {
			return "", nil, err
		}
		//querys中结构必须完全相同
		if source == "" || source != source0 {
			return "", nil, errors.New("the structure of query must be exactly the same")
		}
		if n == 0 {
			sql.WriteString(" where   " + source + " ")
		}
		if len(s.querys) == 1 {
			for i, param := range params {
				params[i] = append(param, p...)
			}
		} else {
			params[n] = append(params[n], p...)
		}
	}

	return sql.String(), params, nil
}
