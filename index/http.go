package index

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/wejick/balasticsearch/util"
)

//GetIndexHandler http handler to get index
func (i *Index) GetIndexHandler(context echo.Context) (err error) {
	i.bleveGetIndex.IndexNameLookup = func(*http.Request) string {
		return context.Param(i.indexNameLookupKey)
	}
	w := util.NewHTTPResponseInterceptor(context.Response().Writer)
	i.bleveGetIndex.ServeHTTP(w, context.Request())
	context.Response().Status = w.StatusCode

	return
}

//GetIndexListHandler http handler get index list
func (i *Index) GetIndexListHandler(context echo.Context) (err error) {
	w := util.NewHTTPResponseInterceptor(context.Response().Writer)
	i.bleveGetIndexList.ServeHTTP(w, context.Request())
	context.Response().Status = w.StatusCode

	return
}
