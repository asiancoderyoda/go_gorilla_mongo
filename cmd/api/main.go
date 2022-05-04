package main

import (
	"context"
	"flag"
	"fmt"
	"go-gorilla-mongo/cmd/api/configs"
	"go-gorilla-mongo/cmd/api/routes"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	port int
	env  string
	wait time.Duration
	db   *mongo.Client
}

type ServerStatus struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Env     string `json:"env"`
}

func main() {
	cfg, err := configureConfigs()

	if err != nil {
		fmt.Printf("Error configuring configs: %s\n", err)
		panic(err)
	}

	prepareAndServe(cfg)
}

func configureConfigs() (*Config, error) {
	configs.LoadEnv()
	conn := configs.ConnectToDB(configs.GetEnvFromKey("MONGO_URI"))
	configs.DB = conn

	port, err := strconv.Atoi(configs.GetEnvFromKey("PORT"))

	if err != nil {
		fmt.Printf("Error getting port: %s\n", err)
		return nil, err
	}

	var cfg Config

	flag.IntVar(&cfg.port, "port", port, "port to listen on")
	flag.StringVar(&cfg.env, "env", configs.GetEnvFromKey("ENV"), "environment")
	flag.DurationVar(&cfg.wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	cfg.db = conn
	flag.Parse()

	return &cfg, nil
}

func prepareAndServe(cfg *Config) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      routes.ConfigureRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Printf("Starting server on port %d in %s mode\n", cfg.port, cfg.env)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			// panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Println("shutting down")
	os.Exit(0)
}
