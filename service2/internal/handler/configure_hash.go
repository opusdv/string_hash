// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"service2/conf"
	"service2/internal/handler/operations"
	"service2/internal/handlers"
	"service2/internal/reqid"
	"service2/system/initlog"

	st "github.com/pkg/errors"
)

var hs *handlers.HashStruct
var l *logrus.Logger

//go:generate swagger generate server --target ../../../service2 --name Hash --spec ../../api/api.yml --server-package internal/handler --principal interface{}

func configureFlags(api *operations.HashAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.HashAPI) http.Handler {
	var once sync.Once
	once.Do(func() {
		conf.InitDefaultConfig()
	})

	cfg := conf.GetConfig()
	l = initlog.InitLog(cfg.LOG_LEVEL)
	hs = handlers.NewHashStruct(cfg, l)
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = l.Infof
	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetCheckHandler == nil {
		api.GetCheckHandler = operations.GetCheckHandlerFunc(func(params operations.GetCheckParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetCheck has not yet been implemented")
		})
	}
	if api.PostSendHandler == nil {
		api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostSend has not yet been implemented")
		})
	}

	api.PostSendHandler = operations.PostSendHandlerFunc(hs.SendHashHandler)
	api.GetCheckHandler = operations.GetCheckHandlerFunc(hs.CheckHashHandler)

	api.PreServerShutdown = func() {
		hs.Wg.Wait()
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = reqid.GenerateReqId(l)
		}
		ctx := context.WithValue(r.Context(), "requestID", reqID)
		l.WithFields(logrus.Fields{
			"Message":     "Debug Info",
			"Method":      r.Method,
			"Body":        r.Body,
			"Request URI": r.RequestURI,
			"Request ID":  reqID,
		}).Debug("Debug midleware")
		err := hs.UpdateConfig()
		if err != nil {
			l.WithFields(logrus.Fields{
				"Error":       err,
				"Stack Trace": st.WithStack(err),
				"Request ID":  reqID,
			}).Error("Error update config")
		}
		l.Debug("Confog Update")

		handler.ServeHTTP(rw, r.WithContext(ctx))
	})
}
