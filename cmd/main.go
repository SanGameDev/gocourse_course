package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SanGameDev/gocourse_course/internal/course"
	"github.com/SanGameDev/gocourse_course/pkg/bootstrap"
	"github.com/SanGameDev/gocourse_course/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()
	courseRepo := course.NewRepo(l, db)
	courseSrv := course.NewService(l, courseRepo)

	h := handler.NewCourseHTTPServer(ctx, course.MakeEndpoints(courseSrv, course.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", port)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
