package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	otel "go-product-service/drivers/tracing"
	"go-product-service/internal/config"
	"go-product-service/internal/http/middleware"
	"go-product-service/internal/router"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// set otel
	shutdownTraceProvider := otel.NewTraceProvider(context.Background())
	defer shutdownTraceProvider()

	app := gin.Default()
	app.Use(middleware.TraceMiddleware())

	app.GET("/ping", controller)

	router.Route(&app.RouterGroup)

	// create server
	restServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port()),
		Handler: app,
	}

	chanSignal := make(chan os.Signal, 1)
	chanErr := make(chan error, 1)
	chanQuit := make(chan struct{}, 1)

	signal.Notify(chanSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-chanSignal:
				logrus.Warn("receive signal interrupt ⚠️")
				gracefullShutdown(restServer)
				chanQuit <- struct{}{}
				return
			case e := <-chanErr:
				logrus.Warnf("receive error ⚠️ : %s", e.Error())
				gracefullShutdown(restServer)
				chanQuit <- struct{}{}
				return
			}
		}
	}()

	go func() {
		logrus.Infof("http server running listening on port %d ⏳", config.Port())
		if err := restServer.ListenAndServe(); err != nil {
			chanErr <- err
			return
		}
	}()

	<-chanQuit
	close(chanQuit)
	close(chanErr)
	close(chanSignal)

	logrus.Info("server exit ‼️")
}

func gracefullShutdown(srv *http.Server) {
	if srv == nil {
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Error(err)
		_ = srv.Close()
		return
	}

	if err := srv.Close(); err != nil {
		logrus.Error(err)
	}

	logrus.Infof("success gracefull shutdown server ❎")
}

func controller(c *gin.Context) {
	ctx, span := otel.Start(c)
	defer span.End()

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"traceID": ctx.Value("traceID").(string),
	})
}

func controllerProduct(c *gin.Context) {
	ctx, span := otel.Start(c)
	defer span.End()

	traceID := ctx.Value("traceID").(string)

	// get id from parms
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"traceID": traceID,
		})
		return
	}

	// call method in service
	product, err := serviceProduct(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"traceID": traceID,
		})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"message": "success get product",
		"traceID": traceID,
		"data":    product,
	})
}

func serviceProduct(ctx context.Context, id int) (string, error) {
	ctx, span := otel.Start(ctx)
	defer span.End()

	respons := spew.Sprintf("product with ID %d", id)
	return respons, nil
}
