package user

import  (
    "github.com/astaxie/beego/orm"
    "fmt"
    proto "Micro_API_Framework/user_service/proto/user"
)

type User struct {
    Id     int64    `orm:"auto"`
    UserName   string `orm:"size(10)"`
    UserId string `orm:"size(10)"`
    Pwd  string `orm:"size(32)"`
}



func QueryUserByName(userName string) (ret *proto.User, err error) {


	var user User

    o := orm.NewOrm()
    qs := o.QueryTable("user")

    errs := qs.Filter("UserName", userName).One(&user)
    fmt.Println(errs)
    if errs != nil {
        fmt.Println(errs)
        
    }
 
    ret = &proto.User{}
    ret.Id = user.Id
    ret.Name = user.UserName
    ret.Pwd = user.Pwd
	return ret , errs
}


