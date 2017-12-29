package index

import (
	"os"

	"github.com/wejick/balasticsearch/registry"
)

//IIndex interface for index
type IIndex interface {
	Get() (map[string]Info, error)
	Create() error
}

//Index main struct
type Index struct {
	dataDirectory string
	indexRegistry *registry.IndexRegistry
}

//New initialize new instance of index api
func New(indexRegistry *registry.IndexRegistry, options ...func(*Index)) (I *Index, err error) {
	I = &Index{
		indexRegistry: indexRegistry,
	}

	for _, option := range options {
		option(I)
	}

	return
}

//UseDataDir specify data directory to be used by balasticsearch
func UseDataDir(dataDirectoryPath string) func(*Index) {
	return func(i *Index) {
		i.dataDirectory = dataDirectoryPath
	}
}

func (I *Index) indexPath(name string) string {
	return I.dataDirectory + string(os.PathSeparator) + name
}
