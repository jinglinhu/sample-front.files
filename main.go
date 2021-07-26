package main

import (
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/context" 
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"log"

	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/fasthttp/router"
    "github.com/valyala/fasthttp"


    "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

)

func init() {
	xray.Configure(xray.Config{
		DaemonAddr:     "xray-service.default:2000",
		LogLevel:       "info",
	})
}


func main() {
    fh := xray.NewFastHTTPInstrumentor(nil)
    r := router.New()
    
    r.GET("/", middleware("x-ray-sample-index-k8s", index, fh))
    r.GET("/frontend/{trace?}", middleware("x-ray-sample-front-k8s", frontend, fh))
    
    r.GET("/opentelemetry",opentelemetry)

    fasthttp.ListenAndServe(":8080", r.Handler)
}

func middleware(name string, h fasthttp.RequestHandler, fh xray.FastHTTPHandler) fasthttp.RequestHandler {
    f := func(ctx *fasthttp.RequestCtx) {
        h(ctx)
    }
    return fh.Handler(xray.NewFixedSegmentNamer(name), f)
}

func index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/html")
	ctx.WriteString(html)
}

func frontend(ctx *fasthttp.RequestCtx) {
	trace := ctx.UserValue("trace")
	switch trace {
      case "middle": 
      	url := `http://x-ray-sample-middle-k8s.default.svc.cluster.local`
	    resp,err := traceUrl(ctx,url)
		if err != nil {
	        fmt.Println("请求失败:", err.Error())
	        return
	    }
	    fmt.Fprintf(ctx,"Trace:Frontend Server->"+string(resp))

      case "backend": 

      	url := `http://x-ray-sample-back-k8s.default.svc.cluster.local`
	    resp,err := traceUrl(ctx,url)
		if err != nil {
	        fmt.Println("请求失败:", err.Error())
	        return
	    }
	    fmt.Fprintf(ctx,"Trace:Frontend Server->"+string(resp))

      case "all": 

      	url := `http://x-ray-sample-middle-k8s.default.svc.cluster.local/all`
	    resp,err := traceUrl(ctx,url)
		if err != nil {
	        fmt.Println("请求失败:", err.Error())
	        return
	    }
	    fmt.Fprintf(ctx,"Trace:Frontend Server->"+string(resp))

      default: 
    	fmt.Fprintf(ctx,"Trace:Frontend Server")
   }
}

var tr = &http.Transport{
	MaxIdleConns: 20,
	IdleConnTimeout: 30 * time.Second,
}

func traceUrl(ctx context.Context,url string) ([]byte, error) {
    resp, err := ctxhttp.Get(ctx, xray.Client(&http.Client{Transport: tr}), url)
    if err != nil {
      return nil, err
    }
    return ioutil.ReadAll(resp.Body)
}



//---------------------分割符------------------
func opentelemetry(ctxx *fasthttp.RequestCtx) {
	tp, err := tracerProvider("http://simplest-collector.default.svc.cluster.local:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-fontend")

	ctx, span := tr.Start(ctx, "fontend")
	defer span.End()

	otmmiddle(ctx)

	fmt.Fprintf(ctxx,"Opentelemetry:Trace-demo")

}

func otmmiddle(ctx context.Context) {
	tr := otel.Tracer("component-middle")
	_, span := tr.Start(ctx, "middle")
	span.SetAttributes(attribute.Key("testset").String("mid-value"))
	defer span.End()

	otmbackend(ctx)
}

func otmbackend(ctx context.Context) {
	tr := otel.Tracer("component-backend")
	_, span := tr.Start(ctx, "backend")
	span.SetAttributes(attribute.Key("testset").String("back-value"))
	defer span.End()
}

const (
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}
//------------分割符---------------




var html = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
    <head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<title>Xray Test</title>
	<style>
body {
    color: #222222;
    background-color: #e0ebf5;
    font-family: Arial, sans-serif;
    font-size:14px;
    -moz-transition-property: text-shadow;
    -moz-transition-duration: 4s;
    -webkit-transition-property: text-shadow;
    -webkit-transition-duration: 4s;
    text-shadow: none;
}
    body.blurry {
	-moz-transition-property: text-shadow;
	-moz-transition-duration: 4s;
	-webkit-transition-property: text-shadow;
	-webkit-transition-duration: 4s;
	text-shadow: #fff 0px 0px 25px;
    }
    a {
	color: #0188cc;
    }
    .textColumn, .linksColumn {
	padding: 6em;
    }
    .textColumn {
	position: absolute;
	top: 0px;
	right: 50%;
	bottom: 0px;
	left: 0px;

	text-align: right;
	padding-top: 11em;
	background-color: #e0ebf5; 
    }
    .textColumn p {
	width: 75%;
	float:right;
    }
    .linksColumn {
	position: absolute;
	top:0px;
	right: 0px;
	bottom: 0px;
	left: 15%;
	background-color: #ffffff;
    }

    h1 {
	font-size: 500%;
	font-weight: normal;
	margin-bottom: 0em;
	color: #375eab;
    }
    h2 {
	font-size: 200%;
	font-weight: normal;
	margin-bottom: 0em;
	color: #375eab;
    }
    ul {
	padding-left: 1em;
	margin: 0px;
    }
    li {
	margin: 1em 0em;
	font-size:20px
    }

	</style>
    </head>
    <body id="sample">
	<div class="linksColumn"> 
	    <h2>Let's Test Amazon Xray & Opentelemetry</h2>
	    <ul>
		<li><a href="/frontend">Frontend</a></li>
		<li><a href="/frontend/middle">Frontend->Middle</a></li>
		<li><a href="/frontend/backend">Frontend->Backend</a></li>
		<li><a href="/frontend/all">Frontend->Middle->Banckend</a></li>
		<li><a href="/opentelemetry">Frontend->Middle->Banckend(Trace By Opentelemetry)</a></li>
	    </ul>
	    <br><br>
		<div id="api-response">
		</div>
	</div>
    </body>
</html>`

