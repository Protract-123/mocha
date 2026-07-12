package manifest

import (
	"reflect"
	"testing"
)

func TestParseRefString(tester *testing.T) {
	type args struct {
		refString string
	}

	tests := []struct {
		name    string
		args    args
		want    Ref
		wantErr bool
	}{
		{
			name:    "valid ref: name",
			args:    args{refString: "test"},
			want:    Ref{Name: "test"},
			wantErr: false,
		},
		{
			name:    "valid ref: name and bucket",
			args:    args{refString: "test/test"},
			want:    Ref{Name: "test", Bucket: "test"},
			wantErr: false,
		},
		{
			name:    "valid ref: name and version",
			args:    args{refString: "test@1.0.0"},
			want:    Ref{Name: "test", Version: "1.0.0"},
			wantErr: false,
		},
		{
			name:    "valid ref: name, bucket and version",
			args:    args{refString: "test/test@1.0.0"},
			want:    Ref{Name: "test", Bucket: "test", Version: "1.0.0"},
			wantErr: false,
		},

		{
			name:    "malformed ref: empty string",
			args:    args{refString: ""},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: ending @",
			args:    args{refString: "test@"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: leading /",
			args:    args{refString: "/test"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: contains /@",
			args:    args{refString: "test/@test"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: ending /",
			args:    args{refString: "test/"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: only @",
			args:    args{refString: "@"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: only /",
			args:    args{refString: "/"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: @ before /",
			args:    args{refString: "test@1.0/beta"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: two /",
			args:    args{refString: "a/b/c"},
			want:    Ref{},
			wantErr: true,
		},
		{
			name:    "malformed ref: two @",
			args:    args{refString: "test@1.0@2.0"},
			want:    Ref{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			got, err := ParseRefString(test.args.refString)
			if (err != nil) != test.wantErr {
				t.Errorf("ParseRefString() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseRefString() got = %v, want %v", got, test.want)
			}
		})
	}
}
