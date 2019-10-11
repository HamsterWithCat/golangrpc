package middle

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/smallnest/rpcx/share"
	"log"
)

func GenerateSpanWithContext(ctx context.Context, operationName string) (opentracing.Span, context.Context, error) {
	md := ctx.Value(share.ReqMetaDataKey)

	var span opentracing.Span
	var parentSpan opentracing.Span

	tracer := opentracing.GlobalTracer()

	if md != nil {
		carrier := opentracing.TextMapCarrier(md.(map[string]string))
		spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("metadata error %s\n", err)
		} else {
			parentSpan = tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
		}
	}

	if parentSpan != nil {
		span = tracer.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		span = opentracing.StartSpan(operationName)
	}

	metadata := opentracing.TextMapCarrier(make(map[string]string))

	_ = tracer.Inject(span.Context(), opentracing.TextMap, metadata)

	ctx = context.WithValue(context.Background(), share.ReqMetaDataKey, (map[string]string)(metadata))
	return span, ctx, nil
}
