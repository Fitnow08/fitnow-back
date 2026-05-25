package main

import (
	"context"
	"fmt"
	_ "github.com/Sanchir01/fitnow/docs"
	"github.com/Sanchir01/fitnow/internal/app"
	"github.com/Sanchir01/fitnow/internal/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title 🚀 FITNOW
// @version         0.0.1
// @description This is a gateway server
// @termsOfService  http://swagger.io/terms/

// @host localhost:7111
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token (you can use it with or without "Bearer " prefix). Example: "eyJhbGc..." or "Bearer eyJhbGc..."

// @contact.url https://github.com/Sanchir01
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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := apps.HTTPServer.Gracefull(shutdownCtx); err != nil {
		fmt.Println(err)
	}
	apps.Log.Info("Shutting down application...")
}
