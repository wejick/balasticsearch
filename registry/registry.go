//  Copyright (c) 2014 Couchbase, Inc.
//  Copyright (c) 2017 Gian Giovani.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve"
)

//An Registry is to keep and maintaint indices available on runtime
type Registry interface {
	RegisterIndexName(name string, idx bleve.Index)
	UnregisterIndexByName(name string) bleve.Index
	IndexByName(name string) bleve.Index
	IndexNames() []string
	UpdateAlias(alias string, add, remove []string) error
}

//IndexRegistry impelements index registry interface
type IndexRegistry struct {
	indexNameMapping     map[string]bleve.Index
	indexNameMappingLock sync.RWMutex
}

//New creates indexRegistry instance
func New() (NewRegistry *IndexRegistry) {
	NewRegistry = &IndexRegistry{}
	NewRegistry.indexNameMapping = make(map[string]bleve.Index)
	return
}

//RegisterIndexName adds new index to registry
func (I *IndexRegistry) RegisterIndexName(name string, idx bleve.Index) {
	I.indexNameMappingLock.Lock()
	defer I.indexNameMappingLock.Unlock()

	I.indexNameMapping[name] = idx
}

//UnregisterIndexByName removes index from registry
func (I *IndexRegistry) UnregisterIndexByName(name string) bleve.Index {
	I.indexNameMappingLock.Lock()
	defer I.indexNameMappingLock.Unlock()

	if I.indexNameMapping == nil {
		return nil
	}
	rv := I.indexNameMapping[name]
	if rv != nil {
		delete(I.indexNameMapping, name)
	}
	return rv
}

//IndexByName gets index by name
func (I *IndexRegistry) IndexByName(name string) bleve.Index {
	I.indexNameMappingLock.RLock()
	defer I.indexNameMappingLock.RUnlock()

	return I.indexNameMapping[name]
}

//IndexNames gets the list of index names
func (I *IndexRegistry) IndexNames() []string {
	I.indexNameMappingLock.RLock()
	defer I.indexNameMappingLock.RUnlock()

	rv := make([]string, len(I.indexNameMapping))
	count := 0
	for k := range I.indexNameMapping {
		rv[count] = k
		count++
	}
	return rv
}

//UpdateAlias adds or removes indexAlias
func (I *IndexRegistry) UpdateAlias(alias string, add, remove []string) error {
	I.indexNameMappingLock.Lock()
	defer I.indexNameMappingLock.Unlock()

	index, exists := I.indexNameMapping[alias]
	if !exists {
		// new alias
		if len(remove) > 0 {
			return fmt.Errorf("cannot remove indexes from a new alias")
		}
		indexes := make([]bleve.Index, len(add))
		for i, addIndexName := range add {
			addIndex, indexExists := I.indexNameMapping[addIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", addIndexName)
			}
			indexes[i] = addIndex
		}
		indexAlias := bleve.NewIndexAlias(indexes...)
		I.indexNameMapping[alias] = indexAlias
	} else {
		// something with this name already exists
		indexAlias, isAlias := index.(bleve.IndexAlias)
		if !isAlias {
			return fmt.Errorf("'%s' is not an alias", alias)
		}
		// build list of add indexes
		addIndexes := make([]bleve.Index, len(add))
		for i, addIndexName := range add {
			addIndex, indexExists := I.indexNameMapping[addIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", addIndexName)
			}
			addIndexes[i] = addIndex
		}
		// build list of remove indexes
		removeIndexes := make([]bleve.Index, len(remove))
		for i, removeIndexName := range remove {
			removeIndex, indexExists := I.indexNameMapping[removeIndexName]
			if !indexExists {
				return fmt.Errorf("index named '%s' does not exist", removeIndexName)
			}
			removeIndexes[i] = removeIndex
		}
		indexAlias.Swap(addIndexes, removeIndexes)
	}
	return nil
}
