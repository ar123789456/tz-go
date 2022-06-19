package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	handler "tz/moduls/products/delivery/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server *http.Server
	db     *bolt.DB
}

func NewApp() App {
	db := initDB()
	return App{
		db: db,
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	handler.RegisterProduct(router, a.db)

	a.Server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.Server.Shutdown(ctx)
}

func initDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte("products"))
		if err != nil {
			return err
		}
		_, err = t.CreateBucketIfNotExists([]byte("name"))
		return err
	})
	if err != nil {
		log.Panicln(err)
	}
	return db
}
