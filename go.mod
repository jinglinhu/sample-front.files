module github.com/aws-samples/eks-workshop/content/x-ray/sample-front

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/fasthttp/router v1.4.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/valyala/fasthttp v1.28.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	golang.org/x/net v0.0.0-20210510120150-4163338589ed
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.38.0 // indirect
)
