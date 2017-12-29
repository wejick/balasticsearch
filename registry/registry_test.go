package registry

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/blevesearch/bleve"
)

func Test_registry(t *testing.T) {
	const (
		datadir            = "./data"
		index1Name         = "index1"
		index2Name         = "index2"
		index1NameNothing  = "index1Nothing"
		index2NameNothing  = "index2Nothing"
		index1Alias        = "index1Alias"
		index1AliasNothing = "index1AliasNothing"
	)

	R := New()

	//Add index
	indexDir := datadir + time.Now().String()
	newIndex1, err := bleve.New(indexPath(indexDir, index1Name), bleve.NewIndexMapping())
	if err != nil {
		err = fmt.Errorf("[ERROR] creating index: %v", err)
		t.Error(err)
		return
	}
	newIndex2, err := bleve.New(indexPath(indexDir, index2Name), bleve.NewIndexMapping())
	if err != nil {
		err = fmt.Errorf("[ERROR] creating index: %v", err)
		t.Error(err)
		return
	}
	R.RegisterIndexName(index1Name, newIndex1)
	R.RegisterIndexName(index2Name, newIndex2)

	//Test get index
	tIndex := R.IndexByName(index1Name)
	assert.Equal(t, newIndex1, tIndex, "IndexByName found : should be equal to index1")
	assert.Equal(t, R.IndexByName("notfound"), nil, "IndexByName not found : should be equal to nil")

	//Test add aliases
	err = R.UpdateAlias(index1Alias, []string{index1Name}, []string{})
	assert.Equal(t, nil, err, "UpdateAlias add : should be no error")
	err = R.UpdateAlias(index1Alias, []string{index1Name + "nothing"}, []string{})
	assert.NotEqual(t, nil, err, "UpdateAlias add not exist : should be error")
	err = R.UpdateAlias(index1Alias+"nothing", []string{index1Name + "nothing"}, []string{})
	assert.NotEqual(t, nil, err, "UpdateAlias add index not exist : should be error")
	err = R.UpdateAlias(index1Alias, []string{index2Name, index1Name}, []string{})
	assert.Equal(t, nil, err, "UpdateAlias add index2 : should be no error")

	//Test get index list
	assert.Equal(t, []string{index1Name, index1Alias, index2Name}, R.GetIndexNames(), "GetIndexNames : should be contains index1 name")

	//Test remove aliases
	err = R.UpdateAlias(index1Alias, []string{}, []string{index1Name + "nothing"})
	assert.NotEqual(t, nil, err, "UpdateAlias remove not exist : should be error")
	err = R.UpdateAlias(index1Alias, []string{}, []string{index1Name})
	assert.Equal(t, nil, err, "UpdateAlias remove : should be no error")
	err = R.UpdateAlias(index1Name, []string{}, []string{index1Name})
	assert.NotEqual(t, nil, err, "UpdateAlias remove not alias : should be error")
	err = R.UpdateAlias(index1NameNothing, []string{}, []string{index1Name})
	assert.NotEqual(t, nil, err, "UpdateAlias remove index of not exist alias : should be error")

	//Test unregister index
	tIndex = R.UnregisterIndexByName(index1Name)
	assert.Equal(t, newIndex1, tIndex, "UnregisterIndexByName : should be equal to index1")
	assert.Equal(t, R.UnregisterIndexByName(index1Name), nil, "UnregisterIndexByName : should be equal to nil")
	assert.Equal(t, R.IndexByName(index1Name), nil, "UnregisterIndexByName  IndexByName: should be equal to nil")
}

func indexPath(dataDirectory, name string) string {
	return dataDirectory + string(os.PathSeparator) + name
}
