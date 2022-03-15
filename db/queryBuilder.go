package db

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

type QueryBuilder struct {
	TableName string

	querySelect string
	query       []string
	order       string
	limit       string
	offset      string

	params map[string]interface{}

	resultQuery string
}

func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		TableName: tableName,
		params:    map[string]interface{}{},
	}
}

func (q QueryBuilder) Select(querySelect string) *QueryBuilder {
	clone := q
	clone.querySelect = "SELECT " + querySelect

	return &clone
}

func (q QueryBuilder) OrWhere(query string) *QueryBuilder {
	clone := q

	clone.query = append(clone.query, "OR "+query)

	return &clone
}

func (q QueryBuilder) AndWhere(query string) *QueryBuilder {
	clone := q
	clone.query = append(clone.query, "AND "+query)

	return &clone
}

func (q QueryBuilder) OrderBy(sort, order string) *QueryBuilder {
	clone := q
	clone.order = "ORDER BY " + sort + " " + order

	return &clone
}

func (q QueryBuilder) AndOrderBy(sort, order string) *QueryBuilder {
	clone := q

	prefix := clone.order + ", "
	if clone.order == "" {
		prefix = "ORDER BY "
	}

	clone.order = prefix + sort + " " + order

	return &clone
}

func (q QueryBuilder) Limit(limit uint32) *QueryBuilder {
	var sqlLimit sql.NullInt32
	if limit > 0 {
		sqlLimit = sql.NullInt32{Int32: int32(limit), Valid: true}
	}

	clone := q
	clone.limit = "LIMIT :limit"
	clone.params[":limit"] = sqlLimit

	return &clone
}

func (q QueryBuilder) Offset(offset uint32) *QueryBuilder {
	clone := q
	clone.offset = "OFFSET :offset"
	clone.params[":offset"] = offset

	return &clone
}

func (q QueryBuilder) SetParameter(key string, value interface{}) *QueryBuilder {
	clone := q

	clone.params[key] = value

	return &clone
}

func (q *QueryBuilder) GetQuery(namedQuery bool) string {
	query := q.buildQuery()

	if !namedQuery {
		// Normalize offset of query with counted args ($<number>)
		regx := regexp.MustCompile(`:\w+\b`)
		counter := 0
		query = regx.ReplaceAllStringFunc(query, func(part string) string {
			counter++
			return "$" + strconv.Itoa(counter)
		})
	}

	return query
}

func (q *QueryBuilder) GetParams() []interface{} {
	var params []interface{}

	query := q.buildQuery()

	// Find sorted params in query
	regx := regexp.MustCompile(`:\w+\b`)
	regx.ReplaceAllStringFunc(query, func(part string) string {
		params = append(params, q.params[part])
		return part
	})

	return params
}

func (q *QueryBuilder) GetNamedParams() map[string]interface{} {
	namedParams := map[string]interface{}{}
	for key, value := range q.params {
		key = strings.ReplaceAll(key, ":", "")
		namedParams[key] = value
	}
	return namedParams
}

func (q *QueryBuilder) buildQuery() string {
	if q.resultQuery != "" {
		return q.resultQuery
	}

	if q.querySelect == "" {
		q.querySelect = "SELECT *"
	}

	if q.limit == "" {
		q.limit = "LIMIT :limit"
		q.params[":limit"] = sql.NullInt32{}
	}

	if q.offset == "" {
		q.offset = "OFFSET :offset"
		q.params[":offset"] = 0
	}

	if len(q.query) > 0 {
		replacer := strings.NewReplacer("OR", "WHERE", "AND", "WHERE")
		q.query[0] = replacer.Replace(q.query[0])
	}

	q.resultQuery = q.querySelect + " " + // SELECT
		"FROM " + strings.ToLower(q.TableName) + " " + // FROM
		strings.Join(q.query, " ") + " " + // WHERE
		q.order + " " + // ORDER BY
		q.limit + " " + // LIMIT
		q.offset + // OFFSET
		";" // END

	return q.resultQuery
}
