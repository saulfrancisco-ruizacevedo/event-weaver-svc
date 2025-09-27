package mediator

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mu       sync.RWMutex
	handlers = make(map[string]interface{})
)

type MediatorHandler[Req any, Res any] interface {
	Handle(ctx context.Context, request Req) (Res, error)
}

func Register[Req any, Res any](name string, handler MediatorHandler[Req, Res]) {
	mu.Lock()
	defer mu.Unlock()
	handlers[name] = handler
}

func Send[Req any, Res any](name string, req Req, ctx *gin.Context) (Res, error) {
	var zero Res

	mu.RLock()
	handlerInstance, ok := handlers[name]
	mu.RUnlock()

	if !ok {
		if ctx != nil {
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("No handler registered for %s", name)})
		}
		return zero, fmt.Errorf("no handler registered for %s", name)
	}

	typedHandler, ok := handlerInstance.(MediatorHandler[Req, Res])
	if !ok {
		if ctx != nil {
			ctx.JSON(500, gin.H{"error": "handler has incompatible type"})
		}
		return zero, fmt.Errorf("handler has incompatible type")
	}

	res, error := typedHandler.Handle(ctx, req)
	if error != nil {
		return zero, error
	}

	return res, nil
}
