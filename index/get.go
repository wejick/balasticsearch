package index

import (
	"errors"

	"github.com/blevesearch/bleve/mapping"
)

type (
	//Info hold index information
	Info struct {
		Aliases map[string]interface{}          `json:"aliases"`
		Types   map[string]mapping.IndexMapping `json:"Mappings"`
	}
)

//Get get index information
func (I *Index) Get(name string) (infos map[string]Info, err error) {
	index := I.indexRegistry.IndexByName(name)
	//Do nothing if index with the same name exist
	if index == nil {
		err = errors.New("[ERROR] No index found : " + name)
		return
	}

	infos = make(map[string]Info, 1)
	infos[name] = Info{
		Types: make(map[string]mapping.IndexMapping),
	}
	infos[name].Types[name] = index.Mapping()

	return
}
