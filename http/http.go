package http

import (
	"log"

	"github.com/wejick/balasticsearch/index"
)

//Handler http handler main struct
type Handler struct {
	indexModules *index.Index

	indexNameIdentifier string
}

//New initialize new http instance
func New(options ...func(*Handler)) (newHTTPhandler *Handler) {
	newHTTPhandler = &Handler{}

	for _, option := range options {
		option(newHTTPhandler)
	}

	if newHTTPhandler.indexModules == nil {
		log.Println("[ERROR] Index instance is not yet set")
	}

	return
}

//UseIndex specify index instance to use
func UseIndex(indexModules *index.Index) func(*Handler) {
	return func(H *Handler) {
		H.indexModules = indexModules
	}
}

//UseIndexNameIdentifier specify index name identifier used in the http router
func UseIndexNameIdentifier(indexNameIdentifier string) func(*Handler) {
	return func(H *Handler) {
		H.indexNameIdentifier = indexNameIdentifier
	}
}
