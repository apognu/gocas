package interceptor

import "net/http"

type Interceptor interface {
	Init()
	Intercept(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

var AvailableInterceptors = map[string]Interceptor{
	"throttling": ThrottlingInterceptor{},
}
