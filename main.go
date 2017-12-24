package main

import (
	"errors"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	index "github.com/wejick/balasticsearch/index"
)

const (
	indexNameStr = "indexName"
)

var (
	e            *echo.Echo // echo instance
	indexModules *index.Index
)

//initModules initialize modules used in balastic
//some functionality  depend on echo
func initModules() (err error) {
	//initialize index module
	indexModules, err = index.New(
		index.UseDataDir("./data"),
		index.UseIndexNameLookupKey(indexNameStr),
	)
	if err != nil {
		return errors.New("[ERROR] couldn't initialize Index module :" + err.Error())
	}

	return
}

func main() {
	e = echo.New()
	e.HideBanner = true

	err := initModules()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//index API
	e.GET("/:"+indexNameStr, indexModules.GetIndexHandler)

	//cat API
	e.GET("/_cat/indices", indexModules.GetIndexListHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
