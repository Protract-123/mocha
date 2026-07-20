package manifest

import (
	"reflect"
	"testing"
)

func TestParseSpec(tester *testing.T) {
	type args struct {
		spec string
	}

	tests := []struct {
		name    string
		args    args
		want    Info
		wantErr bool
	}{
		{
			name:    "valid spec: name",
			args:    args{spec: "test"},
			want:    Info{Name: "test"},
			wantErr: false,
		},
		{
			name:    "valid spec: name and bucket",
			args:    args{spec: "test/test"},
			want:    Info{Name: "test", Bucket: "test"},
			wantErr: false,
		},
		{
			name:    "valid spec: name and version",
			args:    args{spec: "test@1.0.0"},
			want:    Info{Name: "test", Version: "1.0.0"},
			wantErr: false,
		},
		{
			name:    "valid spec: name, bucket and version",
			args:    args{spec: "test/test@1.0.0"},
			want:    Info{Name: "test", Bucket: "test", Version: "1.0.0"},
			wantErr: false,
		},

		{
			name:    "malformed spec: empty string",
			args:    args{spec: ""},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: ending @",
			args:    args{spec: "test@"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: leading /",
			args:    args{spec: "/test"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: contains /@",
			args:    args{spec: "test/@test"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: ending /",
			args:    args{spec: "test/"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: only @",
			args:    args{spec: "@"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: only /",
			args:    args{spec: "/"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: @ before /",
			args:    args{spec: "test@1.0/beta"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: two /",
			args:    args{spec: "a/b/c"},
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "malformed spec: two @",
			args:    args{spec: "test@1.0@2.0"},
			want:    Info{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			got, err := ParseSpec(test.args.spec)
			if (err != nil) != test.wantErr {
				t.Errorf("ParseSpec() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseSpec() got = %v, want %v", got, test.want)
			}
		})
	}
}
