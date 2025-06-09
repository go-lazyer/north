package norm

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
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
	update     map[string]any
	updates    []map[string]any
	insert     map[string]any
	inserts    []map[string]any
	joins      []*Join
	tableName  string
	tableAlias string
	primary    string //主键
	columns    []string
}

func CreateOrm() *Orm {
	return new(Orm)
}

func (s *Orm) Where(query ...BaseQuery) *Orm {
	if s.querys == nil {
		s.querys = make([]BaseQuery, 0)
	}
	s.querys = append(s.querys, query...)
	return s
}

func (s *Orm) Update(m map[string]any) *Orm {
	s.update = m
	return s
}
func (s *Orm) Updates(m []map[string]any) *Orm {
	s.updates = m
	return s
}
func (s *Orm) Insert(m map[string]any) *Orm {
	s.insert = m
	return s
}
func (s *Orm) Inserts(m []map[string]any) *Orm {
	s.inserts = m
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
func (s *Orm) Primary(primary string) *Orm {
	s.primary = primary
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

func (s *Orm) CountSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	params := make([]any, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("select ")

	result := " count(*) count  "
	if s.columns != nil && len(s.columns) == 1 {
		result = strings.Join(s.columns, ",")
	}

	sql.WriteString(result)
	sql.WriteString(" from  " + s.tableName + "")

	if s.tableAlias != "" {
		sql.WriteString(" " + s.tableAlias + " ")
	}

	if s.joins != nil && len(s.joins) > 0 {
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
	}

	if s.querys != nil && len(s.querys) > 0 {
		n := 0
		table := s.tableName
		if s.tableAlias != "" {
			table = s.tableAlias
		}
		for _, query := range s.querys {
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
			n = n + 1
		}
	}

	return sql.String(), params, nil
}

func (s *Orm) SelectSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
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

	if s.joins != nil && len(s.joins) > 0 {
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
	}

	if s.querys != nil && len(s.querys) > 0 {
		n := 0
		table := s.tableName
		if s.tableAlias != "" {
			table = s.tableAlias
		}
		for _, query := range s.querys {
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
			n = n + 1
		}
	}
	if s.groupBy != nil && len(s.groupBy) > 0 {
		sql.WriteString(" group by   ")
		for n, v := range s.groupBy {
			if n != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(v)
		}
	}
	if s.orderBy != nil && len(s.orderBy) > 0 {
		sql.WriteString(" order by   ")
		for n, v := range s.orderBy {
			if n != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(v)
		}
	}
	if s.pageSize > 0 {
		if s.pageNum > 0 {
			s.pageStart = (s.pageNum - 1) * s.pageSize
		}
		params = append(params, s.pageSize, s.pageStart)
		if prepare {
			sql.WriteString(fmt.Sprintf(" limit %s offset %s", PLACE_HOLDER_GO, PLACE_HOLDER_GO))
		} else {
			sql.WriteString(fmt.Sprintf(" limit %d offset %d", s.pageSize, s.pageStart))
		}
	}

	return sql.String(), params, nil
}

func (s *Orm) DeleteSql(prepare bool) (string, []any, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName cannot be empty")
	}
	if s.querys == nil || len(s.querys) != 1 {
		return "", nil, errors.New("the querys size must be 1")
	}
	params := make([]any, 0, 10)
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

	return sql.String(), params, nil
}
func (s *Orm) InsertSql(prepare bool) (string, []any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName  cannot be empty")
	}
	n := 0
	params := make([]any, 0)
	fields := make([]string, 0)
	var sql bytes.Buffer
	sql.WriteString("insert into " + s.tableName + " ")
	sql.WriteString("(")
	if s.inserts != nil && len(s.inserts) > 0 {
		//把所有要修改的字段提取出来

		for field, _ := range s.inserts[0] {
			fields = append(fields, field)
		}

		for _, field := range fields {
			if n != 0 {
				sql.WriteString(",")
			}
			sql.WriteString(" " + field + " ")
			n++
		}
		sql.WriteString(") values")
		n = 0

		for _, maps := range s.inserts {
			if n != 0 {
				sql.WriteString(",")
			}
			sql.WriteString("(")
			m := 0
			for _, field := range fields {
				if m != 0 {
					sql.WriteString(",")
				}
				params = append(params, maps[field])
				if prepare {
					sql.WriteString(fmt.Sprintf(" %s ", PLACE_HOLDER_GO))
				} else {
					sql.WriteString(fmt.Sprintf(" '%v' ", maps[field]))
				}
				m++
			}
			sql.WriteString(")")
			n++
		}
	} else {
		for field, _ := range s.insert {
			fields = append(fields, field)
		}
		for _, field := range fields {
			if n != 0 {
				sql.WriteString(",")
			}
			sql.WriteString(" " + field + " ")
			n++
		}
		sql.WriteString(") values")
		n = 0
		sql.WriteString("(")
		m := 0
		for _, field := range fields {
			if m != 0 {
				sql.WriteString(",")
			}
			params = append(params, s.insert[field])
			if prepare {
				sql.WriteString(fmt.Sprintf(" %s ", PLACE_HOLDER_GO))
			} else {
				sql.WriteString(fmt.Sprintf(" '%v' ", s.insert[field]))
			}
			m++
		}
		sql.WriteString(")")
	}

	return sql.String(), params, nil
}

func (s *Orm) UpdateSql(prepare bool) (string, []any, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName  cannot be empty")
	}

	if s.querys == nil || len(s.querys) != 1 {
		return "", nil, errors.New("the querys size must be 1")
	}

	params := make([]any, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("update " + s.tableName + " set ")
	n := 0
	if s.updates != nil && len(s.updates) > 0 { //批量更新

		if s.primary == "" {
			return "", nil, errors.New("primary cannot be empty")
		}

		//把所有要修改的字段提取出来
		fields := make(map[string]string)
		for _, setMap := range s.updates {
			for name, _ := range setMap {
				fields[name] = ""
			}
		}

		for field, _ := range fields {
			if n != 0 {
				sql.WriteString(",")
			}
			sql.WriteString(fmt.Sprintf("%v = CASE %v", field, s.primary))
			for _, setMap := range s.updates {
				v, ok := setMap[field]
				if !ok {
					continue
				}
				params = append(params, setMap[s.primary], v)
				if prepare {
					sql.WriteString(fmt.Sprintf(" WHEN %s THEN %s", PLACE_HOLDER_GO, PLACE_HOLDER_GO))
				} else {
					sql.WriteString(fmt.Sprintf(" WHEN '%v' THEN '%v'", setMap[s.primary], v))
				}
			}
			sql.WriteString(" END ")
			n++
		}
	} else { //单个更新
		for name, value := range s.update {
			if n != 0 {
				sql.WriteString(",")
			}
			if prepare {
				sql.WriteString(fmt.Sprintf("%v=%s", name, PLACE_HOLDER_GO))
			} else {
				sql.WriteString(fmt.Sprintf("%v='%v'", name, value))
			}
			params = append(params, value)
			n++
		}
	}

	source, param, _ := s.querys[0].Source(s.tableName, prepare)
	sql.WriteString(" where " + source + " ")
	params = append(params, param...)

	return sql.String(), params, nil
}
