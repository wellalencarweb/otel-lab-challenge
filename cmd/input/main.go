package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/wellalencarweb/otel-lab-challenge/config"
	"github.com/wellalencarweb/otel-lab-challenge/internal/pkg/dependencies"
	opentelemetry "github.com/wellalencarweb/otel-lab-challenge/internal/pkg/otel"
)

func main() {
	configs, configsErr := config.LoadConfig(".")
	if configsErr != nil {
		log.Fatal(configsErr)
	}

	deps := dependencies.ResolveInputServiceDependencies(configs)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	otelProviderShutdownFn, err := opentelemetry.InitProvider(
		ctx,
		deps.ServiceName,
		configs.OtelCollectorURL,
	)
	if err != nil {
		log.Fatalf("failed to initialize the otel provider: %s", err)
	}

	deps.WebServer.Start()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("shutting webserver down...\n")

		if err := deps.WebServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("error shutting webserver down: %s\n", err)
		}

		if err := otelProviderShutdownFn(shutdownCtx); err != nil {
			log.Fatalf("error shutting otel provider down: %s\n", err)
		}
	}()

	wg.Wait()
}
