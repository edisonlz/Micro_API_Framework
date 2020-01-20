package handler

import (
	"encoding/json"
	"net/http"
	"time"

	hystrix_go "github.com/afex/hystrix-go/hystrix"
	auth "Micro_API_Framework/auth_service/proto/auth"
	"Micro_API_Framework/plugins/session"
	us "Micro_API_Framework/user_service/proto/user"
	user_service_api "Micro_API_Framework/user_api/service_request"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
)

var (
	serviceClient us.UserService
	authClient    auth.Service
)

// Error 错误结构体
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Init() {
	hystrix_go.DefaultVolumeThreshold = 20
	hystrix_go.DefaultErrorPercentThreshold = 50
	cl := hystrix.NewClientWrapper()(client.DefaultClient)
	serviceClient = us.NewUserService("api.micro.framework.service.user", cl)
	authClient = auth.NewService("api.micro.framework.service.auth", cl)
}



// Login 登录入口
func Login(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	// 只接受POST请求
	if r.Method != "POST" {
		log.Logf("非法请求")
		http.Error(w, "非法请求", 400)
		return
	}
	r.ParseForm()
	rsp, err := user_service_api.ServiceQueryUserName(serviceClient, r, r.Form.Get("userName"))
	
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// 返回结果
	response := map[string]interface{}{
		"ref": time.Now().UnixNano(),
	}
	
	if rsp.User.Pwd == r.Form.Get("pwd") {
		// 生成token
		token_response, err := user_service_api.ServiceAuth(authClient,r,rsp)
		if err != nil {
			log.Logf("[Login] 创建token失败，err：%s", err)
			http.Error(w, err.Error(), 500)
			return
		}

		response["token"] = token_response.Token
		response["data"] = rsp.User
		response["success"] = true
	
		w.Header().Add("set-cookie", "application/json; charset=utf-8")

		// 过期30分钟
		expire := time.Now().Add(60 * time.Minute)
		cookie := http.Cookie{Name: "remember-me-token", Value: token_response.Token, Path: "/", Expires: expire, MaxAge: 90000}
		http.SetCookie(w, &cookie)

		// 同步到session中
		sess := session.GetSession(w, r)
		sess.Values["userId"] = rsp.User.Id
		sess.Values["userName"] = rsp.User.Name
		_ = sess.Save(r, w)

	} else {
		response["success"] = false
		response["error"] = &Error{
			Detail: "密码错误",
		}
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Logout 退出登录
func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 只接受POST请求
	if r.Method != "POST" {
		log.Logf("非法请求")
		http.Error(w, "非法请求", 400)
		return
	}

	tokenCookie, err := r.Cookie("remember-me-token")
	if err != nil {
		log.Logf("token获取失败")
		http.Error(w, "非法请求", 400)
		return
	}

	// 删除token
	_, err = authClient.DelUserAccessToken(ctx, &auth.Request{
		Token: tokenCookie.Value,
	})
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 清除cookie
	cookie := http.Cookie{Name: "remember-me-token", Value: "", Path: "/", Expires: time.Now().Add(0 * time.Second), MaxAge: 0}
	http.SetCookie(w, &cookie)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// 返回结果
	response := map[string]interface{}{
		"ref":     time.Now().UnixNano(),
		"success": true,
	}

	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
