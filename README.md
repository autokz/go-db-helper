go-db-helper
----
Helper to work with database. Support only postgres!

### DB:
Create your connection db params by implementing PostgreConfigParams:

```
package config

type PostgreParams struct {
}

func (p PostgreParams) GetDbName() string {
	return "db_name"
}

func (p PostgreParams) GetHost() string {
	return "localhost"
}

func (p PostgreParams) GetPort() string {
	return "5432"
}

func (p PostgreParams) GetUser() string {
	return "user"
}

func (p PostgreParams) GetPassword() string {
	return "password"
}
```

Create database source:

```
package main

func main() {
    ds = db.GetDb(config.PostgreParams{}, 0, 0)
    ds.GetConn().SetMaxOpenConns(20)
    ds.GetConn().SetMaxIdleConns(20)
    ds.GetConn().SetConnMaxIdleTime(5 * time.Minute)
}

```

### QueryBuilder:

```go
    queryBuilder := db.NewQueryBuilder("your_table_name").
            Limit(limit).
            Offset(offset).
            OrderBy("created_at", "DESC")

    if dto.Field1 != "" {
        queryBuilder = queryBuilder.
        AndWhere("LOWER(field_1) LIKE :field_1").
        SetParameter(":field_1", strings.ToLower("%"+dto.Field1+"%"))
    }
    
    if dto.Field2 != "" {
        queryBuilder = queryBuilder.
        AndWhere("field_2 = :field_2").
        SetParameter(":field_2", dto.Field2))
    }
    
    if dto.AnotherManyFields != nil && len(dto.AnotherManyFields) > 0 {
        queryBuilder = queryBuilder.StartGroupCondition()
        for i, anotherField := range dto.AnotherManyFields {
            if i == 0 {
                queryBuilder = queryBuilder.AndWhere("another_field = :another_field_" + strconv.Itoa(i))
            } else {
                queryBuilder = queryBuilder.OrWhere("another_field = :another_field_" + strconv.Itoa(i))
            }
                queryBuilder.SetParameter(":another_field_"+strconv.Itoa(i), anotherField)
        }
        queryBuilder = queryBuilder.EndGroupCondition()
    }

    query := queryBuilder.GetQuery(false)
    // query = SELECT COUNT(id) FROM your_table_name WHERE LOWER(field_1) LIKE $1 AND field_2 = $2 AND (another_field = $3 OR another_field = $4) LIMIT $5 OFFSET $6;
	
    err := repo.DbCon.SelectContext(ctx, &deliveryPointUsers, query, queryBuilder.GetParams()...)
```