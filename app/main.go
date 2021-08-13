package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"

	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_authorRepo "gitlab.com/content-management-services/content-service/author/repository/mysql"
	_contentHttpDelivery "gitlab.com/content-management-services/content-service/content/delivery/http"
	_contentHttpDeliveryMiddleware "gitlab.com/content-management-services/content-service/content/delivery/http/middleware"
	_contentRepo "gitlab.com/content-management-services/content-service/content/repository/mysql"
	_contentUcase "gitlab.com/content-management-services/content-service/content/usecase"
)

func init() {
	viper.SetConfigFile(`.env`)
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if viper.GetBool(`DEBUG`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`DB_HOST`)
	dbPort := viper.GetString(`DB_PORT`)
	dbUser := viper.GetString(`DB_USER`)
	dbPass := viper.GetString(`DB_PASS`)
	dbName := viper.GetString(`DB_DATABASE`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _contentHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middL.SENTRY)
	e.Use(middleware.Logger())

	// restricted group
	r := e.Group("")
	r.Use(middL.JWT)

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString(`SENTRY_DSN`),
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _contentRepo.NewMysqlContentRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("APP_TIMEOUT")) * time.Second
	au := _contentUcase.NewContentUsecase(ar, authorRepo, timeoutContext)
	_contentHttpDelivery.NewContentHandler(e, r, au)

	log.Fatal(e.Start(viper.GetString("APP_ADDRESS")))
}
