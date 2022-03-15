# go-db-helper

Helper to work with database. Support only postgres.
----

### Create your connection db params by implementing PostgreConfigParams:

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

### Create database source:

```
package main

func main() {
    ds = db.GetDb(config.PostgreParams{}, 0, 0)
    ds.GetConn().SetMaxOpenConns(20)
    ds.GetConn().SetMaxIdleConns(20)
    ds.GetConn().SetConnMaxIdleTime(5 * time.Minute)
}

```