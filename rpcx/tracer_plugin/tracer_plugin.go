package tracer_plugin

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"net"
)

func GetTracer(service string) (tracer opentracing.Tracer, err error) {
	// zipkin上报
	reporter := zipkinhttp.NewReporter("http://49.235.1.29:9411/api/v2/spans")
	defer reporter.Close()

	// localIp := getLocalHostIP()

	endpoint, err := zipkin.NewEndpoint(service, "127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("zipkin.NewEndPoint new endpoint error,err_msg = %v", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("zipkin.NewTracer create tracer error,err_msg = %v", err)
	}
	tracer = zipkinot.Wrap(nativeTracer)

	return
}

func getLocalHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
