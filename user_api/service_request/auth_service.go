package service_request

import (
    "fmt"
    "net/http"
    "github.com/afex/hystrix-go/hystrix"
    auth "Micro_API_Framework/auth_service/proto/auth"
    us "Micro_API_Framework/user_service/proto/user"
    "errors"
)

func ServiceAuth(serviceClient auth.Service, r *http.Request,rsp *us.Response)(token_response *auth.Response, err error){

    ctx := r.Context()

    output_chan := make(chan *auth.Response, 1)
    errors_chan := make(chan string, 1)


    // 调用后台服务
    hystrix.ConfigureCommand("MakeAccessToken", hystrix.CommandConfig{
        Timeout:               1000, //1秒超时
        MaxConcurrentRequests: 100, //设置max concurrent
        ErrorPercentThreshold: 50, //错误率达到50%降级操作
    })

    // 根据自身业务需求封装到http client调用处
    hystrix.Go("MakeAccessToken", func() error {
        token_response, err := serviceClient.MakeAccessToken(ctx, &auth.Request{
            UserId:   rsp.User.Id,
            UserName: rsp.User.Name,
        })
        fmt.Println(rsp)
        fmt.Println(err)
        output_chan <- token_response

        return nil
    },
    func(err error) error {
        // 失败重试，降级等具体操作
        fmt.Println("get an error, handle it")
        errors_chan <- "error down user backup data"
        return nil
    })

    select {
        case token_response := <-output_chan:
            close(output_chan)
            return token_response , nil
        case err := <-errors_chan:
            close(errors_chan)
            return nil , errors.New(err)
    }
    
    
}

