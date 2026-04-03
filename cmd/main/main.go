package main

import (
	"context"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/app"
	"github.com/Sanchir01/fitnow/internal/handlers"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	apps, err := app.NewApp(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer apps.CancelLogger()
	go func() {
		if err := apps.HTTPServer.Run(handlers.StartHttpHandlers(apps.Handlers)); err != nil {
			fmt.Println(err)
		}
	}()

	apps.Log.Info("Starting application...")
	<-ctx.Done()
	if err := apps.HTTPServer.Gracefull(ctx); err != nil {
		fmt.Println(err)
	}
	apps.Log.Info("Shutting down application...")
}
