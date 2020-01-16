package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"Micro_API_Framework/basic"
	"Micro_API_Framework/basic/common"
	"Micro_API_Framework/basic/config"
	"Micro_API_Framework/plugins/breaker"
	tracer "Micro_API_Framework/plugins/tracer/jaeger"
	"Micro_API_Framework/plugins/tracer/opentracing/std2micro"
	"Micro_API_Framework/user_api/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/config/source/grpc"
	"github.com/opentracing/opentracing-go"
	_ "Micro_API_Framework/plugins/session"
)

//micro new --namespace=api.micro.platform --type=srv --alias=user Micro_API_Framework/user-api
//micro --registry=etcd --api_namespace=api.micro.platform.web  api --handler=web --selector=cache --client_pool_size=10 


var (
	appName = "user_api"
	cfg     = &userCfg{}
)

type userCfg struct {
	common.AppCfg
}

func main() {
	// 初始化配置
	initCfg()

	// 使用etcd注册
	micReg := etcd.NewRegistry(registryOptions)

	t, io, err := tracer.NewTracer(cfg.Name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 创建新服务
	service := web.NewService(
		web.Name("api.micro.platform.web.user"),
		web.Version(cfg.Version),
		web.RegisterTTL(time.Second*15),
		web.RegisterInterval(time.Second*10),
		web.Registry(micReg),
		web.Address(cfg.Addr()),
	)

	// 初始化服务
	if err := service.Init(
		web.Action(
			func(c *cli.Context) {
				// 初始化handler
				handler.Init()
			}),
	); err != nil {
		log.Fatal(err)
	}

	//设置采样率
	std2micro.SetSamplingFrequency(50)
	// 注册登录接口
	handlerLogin := http.HandlerFunc(handler.Login)
	service.Handle("/user/login", std2micro.TracerWrapper(breaker.BreakerWrapper(handlerLogin)))
	// 注册退出接口
	service.Handle("/user/logout", std2micro.TracerWrapper(http.HandlerFunc(handler.Logout)))

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := &common.Etcd{}
	err := config.C().App("etcd", etcdCfg)
	if err != nil {
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.Host, etcdCfg.Port)}
}


func initCfg() {
	source := grpc.NewSource(
		grpc.WithAddress("127.0.0.1:9600"),
		grpc.WithPath("micro"),
	)

	basic.Init(config.WithSource(source))

	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}

	log.Logf("[initCfg] 配置，cfg：%v", cfg)

	return
}


