package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"starter-restapi-golang/app/di"
	"starter-restapi-golang/app/server"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/lestrrat-go/server-starter/listener"
)

const port = 8080

func netListen(network, addr string) (net.Listener, error) {
	ls, err := listener.ListenAll()
	if err == nil {
		return ls[0], nil
	}
	return net.Listen(network, addr)
}

func createRouter(userHandler server.UserHandler) chi.Router {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Route("/users", func(mux chi.Router) {
		mux.Get("/", userHandler.GetUsersHandler)
	})
	return mux
}

func main() {
	run()
}

func run() {
	ctx := context.Background()

	app, cleanup, err := di.NewApp(ctx)
	if err != nil {
		log.Fatalf("server closed with %v", err)
		return
	}
	defer cleanup()

	mux := createRouter(app.UserHandler)
	server := http.Server{
		Handler: mux,
	}

	l, err := netListen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failded to listen: %v", err)
		return
	}

	go func() {
		log.Printf("starting server on %s", l.Addr())
		if err := server.Serve(l); err != nil {
			log.Fatalf("server closed with %v", err)
			return
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %v received, then shutting down...", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to graceful shutdown: %v", err)
	}
	log.Printf("server shutdown")

}
