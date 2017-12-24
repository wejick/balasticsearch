package index

import (
	"log"

	bleveHttp "github.com/blevesearch/bleve/http"
)

//Index main struct
type Index struct {
	dataDirectory string

	//bleveHttpHandler
	indexNameLookupKey string
	bleveGetIndex      *bleveHttp.GetIndexHandler
	bleveGetIndexList  *bleveHttp.ListIndexesHandler
}

//New initialize new instance of index api
func New(options ...func(*Index)) (I *Index, err error) {
	I = &Index{
		bleveGetIndex:     bleveHttp.NewGetIndexHandler(),
		bleveGetIndexList: bleveHttp.NewListIndexesHandler(),
	}

	for _, option := range options {
		option(I)
	}

	if I.indexNameLookupKey == "" {
		log.Println("[WARNING] indexNameLookupKey is not set, using indexName as default key")
	}
	return
}

//UseDataDir specify data directory to be used by balasticsearch
func UseDataDir(dataDirectoryPath string) func(*Index) {
	return func(i *Index) {
		i.dataDirectory = dataDirectoryPath
	}
}

//UseIndexNameLookupKey specify function to lookup indexname on http request
func UseIndexNameLookupKey(indexNameLookupKey string) func(*Index) {
	return func(i *Index) {
		i.indexNameLookupKey = indexNameLookupKey
	}
}
