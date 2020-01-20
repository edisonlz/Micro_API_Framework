package models

//doc --- http://beego.me/docs/mvc/model/orm.md

import  (
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql" // import your used driver
    "Micro_API_Framework/user_service/model/user"
    "fmt"
)
    

/*
    go run orm_sync.go orm syncdb
*/

func init() {

    fmt.Println("[init database]......")

    orm.Debug = true
    //regiter driver
    orm.RegisterDriver("mysql", orm.DRMySQL)
    // register model
    orm.RegisterModel(new(user.User))


    mysql_config := "root:xsw2CDE#@(127.0.0.1:3306)/micro_book_mall?charset=utf8&parseTime=true"

    // set default database
    orm.RegisterDataBase("default", "mysql", mysql_config)
    //set db params

    orm.SetMaxIdleConns("default", 240)
    orm.SetMaxOpenConns("default", 240)

    // set go
    fmt.Println("[end init database]......")

}

