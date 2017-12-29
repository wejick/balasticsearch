package index

import (
	"reflect"
	"testing"
	"time"

	"github.com/blevesearch/bleve/mapping"
	"github.com/wejick/balasticsearch/registry"
)

func TestIndex_CreateIndex(t *testing.T) {
	type args struct {
		name        string
		mappingJSON string
	}
	I, err := New(registry.New(), UseDataDir("./data_create_"+time.Now().String()))
	if err != nil {
		t.Errorf("Index.CreateIndex() error = %v", err)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "index1",
			args: args{
				name: "index1",
			},
			wantErr: false,
		},
		{
			name: "index1 duplicate",
			args: args{
				name: "index1",
			},
			wantErr: true,
		},
		{
			name: "index2",
			args: args{
				name:        "index2",
				mappingJSON: "{}",
			},
			wantErr: false,
		},
		{
			name: "index3 invalid json mapping",
			args: args{
				name:        "index3",
				mappingJSON: "{}-",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := I.Create(tt.args.name, tt.args.mappingJSON); (err != nil) != tt.wantErr {
				t.Errorf("Index.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIndex_Get(t *testing.T) {
	const (
		index1Name = "index1_get"
	)
	I, err := New(registry.New(), UseDataDir("./data_get_"+time.Now().String()))
	if err != nil {
		t.Errorf("Index.GetIndex() error = %v", err)
	}

	//prepare index 1 for testing
	err = I.Create(index1Name, "")
	if err != nil {
		t.Errorf("Index.GetIndex() error = %v", err)
	}
	index1 := I.indexRegistry.IndexByName(index1Name)
	infoIndex1 := make(map[string]Info, 1)
	typeIndex1 := make(map[string]mapping.IndexMapping)
	typeIndex1[index1Name] = index1.Mapping()
	infoIndex1[index1Name] = Info{
		Types: typeIndex1,
	}

	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantInfos map[string]Info
		wantErr   bool
	}{
		{
			name: "index1",
			args: args{
				name: index1Name,
			},
			wantInfos: infoIndex1,
		},
		{
			name: "index2 notfound",
			args: args{
				name: "index2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfos, err := I.Get(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfos, tt.wantInfos) {
				t.Errorf("Index.Get() = %v, want %v", gotInfos, tt.wantInfos)
			}
		})
	}
}
