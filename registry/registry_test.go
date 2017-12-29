package registry

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blevesearch/bleve"
)

func Test_registry(t *testing.T) {
	const (
		datadir     = "./data"
		index1Name  = "index1"
		index1Alias = "index1Alias"
	)

	R := New()

	//Add index
	newIndex, err := bleve.New(indexPath(datadir, index1Name), bleve.NewIndexMapping())
	if err != nil {
		err = fmt.Errorf("[ERROR] creating index: %v", err)
		t.Error(err)
		return
	}
	R.RegisterIndexName(index1Name, newIndex)

	//Test get index
	tIndex := R.IndexByName(index1Name)
	assert.Equal(t, newIndex, tIndex, "IndexByName found : should be equal to index1")
	assert.Equal(t, R.IndexByName("notfound"), nil, "IndexByName not found : should be equal to nil")

	//Test add aliases
	err = R.UpdateAlias(index1Alias, []string{index1Name}, []string{})
	assert.Equal(t, nil, err, "UpdateAlias add : should be no error")

	//Test get index list
	assert.Equal(t, []string{index1Name, index1Alias}, R.IndexNames(), "IndexNames : should be contains index1 name")

	//Test remove aliases
	err = R.UpdateAlias(index1Alias, []string{}, []string{index1Name})
	assert.Equal(t, nil, err, "UpdateAlias remove : should be no error")
	assert.Equal(t, []string{index1Name}, R.IndexNames(), "UpdateAlias IndexNames : should be contains index1 name")
	err = R.UpdateAlias(index1Alias, []string{}, []string{index1Name})
	assert.NotEqual(t, nil, err, "UpdateAlias remove : should be error")

	//Test unregister index
	tIndex = R.UnregisterIndexByName(index1Name)
	assert.Equal(t, newIndex, tIndex, "UnregisterIndexByName : should be equal to index1")
	assert.Equal(t, R.UnregisterIndexByName(index1Name), nil, "UnregisterIndexByName : should be equal to nil")
	assert.Equal(t, R.IndexByName(index1Name), nil, "UnregisterIndexByName  IndexByName: should be equal to nil")
}

func indexPath(dataDirectory, name string) string {
	return dataDirectory + string(os.PathSeparator) + name
}
