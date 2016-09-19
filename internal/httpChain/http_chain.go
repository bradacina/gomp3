package httpChain

import (
	"context"
	"net/http"
)

const (
	chainKey = 93874
)

type Chain struct {
	next    *Chain
	handler http.Handler
}

func (c *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	val, ok := ctx.Value(chainKey).(bool)

	if !ok {
		ctx = context.WithValue(ctx, chainKey, false)
		request = request.WithContext(ctx)
	}

	if val {
		return
	}

	if c.handler != nil {
		c.handler.ServeHTTP(writer, request)
	}

	if c.next != nil {
		c.next.ServeHTTP(writer, request)
	}
}

func BreakChain(request *http.Request) {
	ctx := request.Context()

	ctx = context.WithValue(ctx, chainKey, true)

	request = request.WithContext(ctx)
}

func (c *Chain) NextHandler(handler http.Handler) *Chain {
	current := c
	for current.next != nil {
		current = current.next
	}

	current.next = NewChainWithHandler(handler)
	return c
}

func (c *Chain) NextFunc(handleFunc func(http.ResponseWriter, *http.Request)) *Chain {
	current := c
	for current.next != nil {
		current = current.next
	}

	current.next = NewChainWithFunc(handleFunc)
	return c
}

func NewChainWithHandler(handler http.Handler) *Chain {
	return &Chain{handler: handler}
}

func NewChainWithFunc(handleFunc func(http.ResponseWriter, *http.Request)) *Chain {
	return &Chain{handler: http.HandlerFunc(handleFunc)}
}
