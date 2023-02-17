package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"totoval/bootstrap"
	"totoval/resources/views"
	"totoval/routes"

	c "github.com/totoval/framework/config"
	"github.com/totoval/framework/graceful"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/toto"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/http/middleware"
	"github.com/totoval/framework/request"
	"github.com/totoval/framework/sentry"
)

func init() {
	bootstrap.Initialize()
}

// @caution cannot use config methods to get config in init function
func main() {
	//j := &jobs.ExampleJob{}
	//j.SetParam(&pbs.ExampleJob{Query: "test", PageNumber: 111, ResultPerPage: 222})
	////j.SetDelay(5 * zone.Second)
	//err := job.Dispatch(j)
	//fmt.Println(err)

	//go hub.On("add-user-affiliation")  // go run artisan.go queue:listen add-user-affiliation

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		call := <-quit
		log.Info("system call", toto.V{"call": call})
		cancel()
	}()

	httpServe(ctx)
}

func httpServe(ctx context.Context) {
	r := request.New()

	sentry.Use(r.GinEngine(), false)

	if c.GetBool("app.debug") {
		//r.Use(middleware.RequestLogger())
	}

	r.RedirectTrailingSlash = false

	if c.GetString("app.env") == "production" {
		r.Use(middleware.Logger())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.Locale())

	r.UseGin(gin.BasicAuth(c.Get("webdav.accounts").(gin.Accounts)))

	routes.Register(r)

	views.Initialize(r)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s := &http.Server{
		Addr:           ":" + c.GetString("app.port"),
		Handler:        r,
		ReadTimeout:    zone.Duration(c.GetInt64("app.read_timeout_seconds")) * zone.Second,
		WriteTimeout:   zone.Duration(c.GetInt64("app.write_timeout_seconds")) * zone.Second,
		TLSConfig:      cfg,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-ctx.Done()

	log.Info("Shutdown Server ...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	_ctx, cancel := context.WithTimeout(ctx, 5*zone.Second)
	defer cancel()

	if err := s.Shutdown(_ctx); err != nil {
		log.Fatal("Server Shutdown: ", toto.V{"error": err})
	}

	// totoval framework shutdown
	graceful.ShutDown(false)

	log.Info("Server exiting")
}
