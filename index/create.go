package index

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blevesearch/bleve"
)

//CreateIndex instantiates and index
func (I *Index) CreateIndex(name string, mappingJSON string) (err error) {
	//Do nothing if index with the same name exist
	if I.indexRegistry.IndexByName(name) != nil {
		err = errors.New("[ERROR] Index is already exist : " + name)
		return
	}

	indexMapping := bleve.NewIndexMapping()
	//Try to get index mapping
	if mappingJSON != "" {
		//TODO we can do preprocessing here to make the mapping resemble elasticsearch mapping if it desirable at the future
		json.Unmarshal([]byte(mappingJSON), &indexMapping)
		if err != nil {
			err = fmt.Errorf("error parsing index mapping: %v", err)
			return
		}
	}

	newIndex, err := bleve.New(I.indexPath(name), indexMapping)
	if err != nil {
		err = fmt.Errorf("[ERROR] creating index: %v", err)
		return
	}
	newIndex.SetName(name)
	I.indexRegistry.RegisterIndexName(name, newIndex)

	return
}
