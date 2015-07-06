# go_mysql_tools
common db operation patterns (esp. for mysql)

# Sample Usage 
```go
  import gmtools "github.com/tt7/go_mysql_tools"
  
  dbcfg, err := gmtools.ReadConfigFromJsonFile(configFilename)
  if err!= nil{
    /* ... */
  }
  
  err = gmtools.UseDb(dbcfg, func(db *sql.DB) error {
    if err:=doSomething(); err != nil{
      return err
    }
    return nil
  })
  /* process err ... */
  
  //query
  err = gmtools.QueryDb(db,
    "select a,b,c,d from `sometable` where a=? or b=?",
    []interface{}{1, 2},
    func(rowno int, rows *sql.Rows) error{
      var a,b,c,d int
      err := rows.Scan(&a, &b, &c, &d)
      /* ... */
      return nil  
  })
  
  // transaction
  err = gmtools.InTxWithDB(db, []func(tx *sql.Tx) error {
    func(tx *sql.Tx) error{
      /* do some thing first */
    },
    func(tx *sql.Tx) error{
      /* do some other things */
    }})
```
