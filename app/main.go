package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/rachadiannovansyah/go-echo-clean-arch/modules/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/rachadiannovansyah/go-echo-clean-arch/modules/article/delivery/http/middleware"
	_articleRepo "github.com/rachadiannovansyah/go-echo-clean-arch/modules/article/repository/mysql"
	_articleUcase "github.com/rachadiannovansyah/go-echo-clean-arch/modules/article/usecase"
	_authorRepo "github.com/rachadiannovansyah/go-echo-clean-arch/modules/author/repository/mysql"
	_userHttpDelivery "github.com/rachadiannovansyah/go-echo-clean-arch/modules/user/delivery/http"
	_userRepo "github.com/rachadiannovansyah/go-echo-clean-arch/modules/user/repository/mysql"
	_userUcase "github.com/rachadiannovansyah/go-echo-clean-arch/modules/user/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
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

	// use echo
	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// init repo
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	articleRepo := _articleRepo.NewMysqlArticleRepository(dbConn)
	userRepo := _userRepo.NewMysqlUserRepository(dbConn)

	// init usecase
	articleUsecase := _articleUcase.NewArticleUsecase(articleRepo, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, articleUsecase)
	userUsecase := _userUcase.NewUserUsecase(userRepo, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, userUsecase)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
