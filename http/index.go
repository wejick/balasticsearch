package http

import (
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

//CreateIndex handle http request for index creation
func (H *Handler) CreateIndex(context echo.Context) (err error) {
	indexName := context.Param(H.indexNameIdentifier)
	requestBody, err := ioutil.ReadAll(context.Request().Body)

	if err != nil {
		err = fmt.Errorf("[ERROR] couldn't create index %v", err)
		context.Error(err)
		return
	}

	err = H.indexModules.Create(indexName, string(requestBody[:]))
	if err != nil {
		err = fmt.Errorf("[ERROR] couldn't create index %v", err)
	}

	return
}

//GetIndex handle http request for index information
func (H *Handler) GetIndex(context echo.Context) (err error) {
	indexName := context.Param(H.indexNameIdentifier)
	infos, err := H.indexModules.Get(indexName)
	if err != nil {
		context.Response().Status = 404
	}

	err = context.JSON(200, infos)
	return
}
