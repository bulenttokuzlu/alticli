# Package alticli
The C API for Altibase is called CLI and is basically the same interface as ODBC. However, it can be used without configuring ODBC which simplifies usage for Altibase-only projects.

This driver is based on code.google.com/p/odbc  and https://bitbucket.org/phiggins/db2cli

This package was built with "CGO"


## How to run ( sample code )
```
package main

import (
    "database/sql"
    "flag"
    "fmt"
    _ "github.com/bulenttokuzlu/alticli"
    "os"
    "time"
)

var (
    connStr = flag.String("conn", "Server=127.0.0.1;User=SYS;Password=MANAGER;PORT=20300", "connection string to use")
    repeat  = flag.Uint("repeat", 2, "number of times to repeat query")
)

func usage() {
    fmt.Fprintf(os.Stderr, `usage: %s [options]

%s connects to altibase and executes a simple SQL statement a configurable
number of times.

Here is a sample connection string:

Server=127.0.0.1;User=SYS;Password=MANAGER;PORT=20300
`, os.Args[0], os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}

func execQuery(st *sql.Stmt) error {
    rows, err := st.Query()
    if err != nil {
        return err 
    }   
    defer rows.Close()
    for rows.Next() {
        var t time.Time
        err = rows.Scan(&t)
        if err != nil {
            return err 
        }   
        fmt.Printf("Time: %v\n", t)
    }   
    return rows.Err()
}
func dbOperations() error {
    db, err := sql.Open("alticli", *connStr)
    if err != nil {
        return err
    }
    defer db.Close()
    st, err := db.Prepare("select sysdate from dual;")
    if err != nil {
        return err
    }
    defer st.Close()

    for i := 0; i < int(*repeat); i++ {
        err = execQuery(st)
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
    flag.Usage = usage
    flag.Parse()
    if *connStr == "" {
        fmt.Fprintln(os.Stderr, "-conn is required")
        flag.Usage()
    }

    if err := dbOperations(); err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
}

```

