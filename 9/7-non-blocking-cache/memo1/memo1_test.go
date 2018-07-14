package memo

import (
	"reflect"
	"testing"
)

func TestMemo_Get(t *testing.T) {
	type fields struct {
		f     Func
		cache map[string]result
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memo := &Memo{
				f:     tt.fields.f,
				cache: tt.fields.cache,
			}
			got, err := memo.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Memo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Memo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
