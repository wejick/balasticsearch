package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	index "github.com/wejick/balasticsearch/index"
	"github.com/wejick/balasticsearch/registry"
)

const (
	indexNameStr = "indexName"
	indexDir     = "./data"
)

var (
	e             *echo.Echo // echo instance
	indexRegistry *registry.IndexRegistry
	indexModules  *index.Index
)

func initIndexRegistry(dataDir string) (err error) {
	if indexRegistry != nil {
		log.Println("[WARNING] index registry get reinitialized")
	}

	indexRegistry = registry.New()

	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(dataDir)
	if err != nil {
		err = errors.New("error reading data dir: " + err.Error())
	}

	for _, dirInfo := range dirEntries {
		indexPath := dataDir + string(os.PathSeparator) + dirInfo.Name()

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("[INFO] not registering %s, skipping", indexPath)
			continue
		}

		i, err := bleve.Open(indexPath)
		if err != nil {
			err = errors.New("[ERROR] opening index " + indexPath + " : " + err.Error())
		} else {
			log.Printf("[INFOR] registered index: %s", dirInfo.Name())
			//TODO : need a way to save alias name and associate it to existing index
			indexRegistry.RegisterIndexName(dirInfo.Name(), i)
			// set correct name in stats
			i.SetName(dirInfo.Name())
		}
	}

	return
}

//initModules initialize modules used in balasticsearch
func initModules() (err error) {
	//initialize index registry
	initIndexRegistry(indexDir)

	//initialize index module
	indexModules, err = index.New(
		indexRegistry,
		index.UseDataDir(indexDir), //hardcode for now
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

	e.Logger.Fatal(e.Start(":1323"))
}
