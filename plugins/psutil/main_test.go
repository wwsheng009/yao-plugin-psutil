package main

import (
	"reflect"
	"testing"

	"github.com/yaoapp/kun/grpc"
)

func TestPsUtilPlugin_Exec(t *testing.T) {
	type args struct {
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		plugin  *PsUtilPlugin
		args    args
		want    *grpc.Response
		wantErr bool
	}{
		{
			name:   "test",
			plugin: &PsUtilPlugin{},
			args: struct {
				name string
				args []interface{}
			}{
				name: "disk",
				args: []interface{}{},
			},
			want: &grpc.Response{Bytes: []byte(`{"code":200,"message":"接收成功"}`), Type: "map"},
		},
		{
			name:   "test",
			plugin: &PsUtilPlugin{},
			args: struct {
				name string
				args []interface{}
			}{
				name: "mem",
				args: []interface{}{},
			},
			want: &grpc.Response{Bytes: []byte(`{"code":200,"message":"接收成功"}`), Type: "map"},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.Exec(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PsUtilPlugin.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PsUtilPlugin.Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
