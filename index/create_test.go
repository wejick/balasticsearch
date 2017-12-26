package index

import (
	"testing"

	"github.com/wejick/balasticsearch/registry"
)

func TestIndex_CreateIndex(t *testing.T) {
	type args struct {
		name        string
		mappingJSON string
	}
	I, err := New(registry.New(), UseDataDir("./data"))
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := I.Create(tt.args.name, tt.args.mappingJSON); (err != nil) != tt.wantErr {
				t.Errorf("Index.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
