package main

import (
	"log"

	tracer "Micro_API_Framework/plugins/tracer/jaeger"
	"Micro_API_Framework/plugins/tracer/opentracing/stdhttp"
	"github.com/micro/go-plugins/micro/cors"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	"github.com/opentracing/opentracing-go"
)

//go run main.go --registry=etcd --api_namespace=api.micro.platform.web  api --handler=web

func init() {
	plugin.Register(cors.NewPlugin())

	plugin.Register(plugin.NewPlugin(
		plugin.WithName("tracer"),
		plugin.WithHandler(
			stdhttp.TracerWrapper,
		),
	))
}

const name = "API gateway"

func main() {
	stdhttp.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(name, "")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	cmd.Init()
}



