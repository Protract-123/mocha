package manifest

import "testing"

// test cases retrieved from https://github.com/ScoopInstaller/Scoop/blob/master/test/Scoop-Versions.Tests.ps1

func TestCompareVersions(tester *testing.T) {
	type args struct {
		version1 string
		version2 string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "major.minor.patch progression 1",
			args: args{version1: "0.1.0", version2: "0.1.1"},
			want: 1,
		},
		{
			name: "major.minor.patch progression 2",
			args: args{version1: "0.1.1", version2: "0.2.0"},
			want: 1,
		},
		{
			name: "major.minor.patch progression 3",
			args: args{version1: "0.2.0", version2: "1.0.0"},
			want: 1,
		},

		{
			name: "pre-release versioning progression 1",
			args: args{version1: "0.4.0", version2: "0.5.0-alpha.1"},
			want: 1,
		},
		{
			name: "pre-release versioning progression 2",
			args: args{version1: "0.5.0-alpha.1", version2: "0.5.0-alpha.2"},
			want: 1,
		},
		{
			name: "pre-release versioning progression 3",
			args: args{version1: "0.5.0-alpha.2", version2: "0.5.0-alpha.10"},
			want: 1,
		},
		{
			name: "pre-release versioning progression 4",
			args: args{version1: "0.5.0-alpha.10", version2: "0.5.0-beta"},
			want: 1,
		},
		{
			name: "pre-release versioning progression 5",
			args: args{version1: "0.5.0-beta", version2: "0.5.0-alpha.10"},
			want: -1,
		},
		{
			name: "pre-release versioning progression 6",
			args: args{version1: "0.5.0-beta", version2: "0.5.0-beta.0"},
			want: 1,
		},

		{
			name: "pre-release tags alphabetic order 1",
			args: args{version1: "0.5.0-rc.1", version2: "0.5.0-z"},
			want: 1,
		},
		{
			name: "pre-release tags alphabetic order 2",
			args: args{version1: "0.5.0-rc.1", version2: "0.5.0-howdy"},
			want: -1,
		},
		{
			name: "pre-release tags alphabetic order 3",
			args: args{version1: "0.5.0-howdy", version2: "0.5.0-rc.1"},
			want: 1,
		},

		{
			name: "windows-styled progression 1",
			args: args{version1: "0.0.0.0", version2: "0.0.0.1"},
			want: 1,
		},
		{
			name: "windows-styled progression 2",
			args: args{version1: "0.0.0.1", version2: "0.0.0.2"},
			want: 1,
		},
		{
			name: "windows-styled progression 3",
			args: args{version1: "0.0.0.2", version2: "0.0.1.0"},
			want: 1,
		},
		{
			name: "windows-styled progression 4",
			args: args{version1: "0.0.1.0", version2: "0.0.1.1"},
			want: 1,
		},
		{
			name: "windows-styled progression 5",
			args: args{version1: "0.0.1.1", version2: "0.0.1.2"},
			want: 1,
		},
		{
			name: "windows-styled progression 6",
			args: args{version1: "0.0.1.2", version2: "0.0.2.0"},
			want: 1,
		},
		{
			name: "windows-styled progression 7",
			args: args{version1: "0.0.2.0", version2: "0.1.0.0"},
			want: 1,
		},
		{
			name: "windows-styled progression 8",
			args: args{version1: "0.1.0.0", version2: "0.1.0.1"},
			want: 1,
		},
		{
			name: "windows-styled progression 9",
			args: args{version1: "0.1.0.1", version2: "0.1.0.2"},
			want: 1,
		},
		{
			name: "windows-styled progression 10",
			args: args{version1: "0.1.0.2", version2: "0.1.1.0"},
			want: 1,
		},
		{
			name: "windows-styled progression 11",
			args: args{version1: "0.1.1.0", version2: "0.1.1.1"},
			want: 1,
		},
		{
			name: "windows-styled progression 12",
			args: args{version1: "0.1.1.1", version2: "0.1.1.2"},
			want: 1,
		},
		{
			name: "windows-styled progression 13",
			args: args{version1: "0.1.1.2", version2: "0.2.0.0"},
			want: 1,
		},
		{
			name: "windows-styled progression 14",
			args: args{version1: "0.2.0.0", version2: "1.0.0.0"},
			want: 1,
		},

		{
			name: "partial semver differences 1",
			args: args{version1: "1", version2: "1.1"},
			want: 1,
		},
		{
			name: "partial semver differences 2",
			args: args{version1: "1", version2: "1.0"},
			want: 1,
		},
		{
			name: "partial semver differences 3",
			args: args{version1: "1.1.0.0", version2: "1.1"},
			want: -1,
		},
		{
			name: "partial semver differences 4",
			args: args{version1: "1.4", version2: "1.3.0"},
			want: -1,
		},
		{
			name: "partial semver differences 5",
			args: args{version1: "1.4", version2: "1.3.255.255"},
			want: -1,
		},
		{
			name: "partial semver differences 6",
			args: args{version1: "1.4", version2: "1.4.4"},
			want: 1,
		},
		{
			name: "partial semver differences 7",
			args: args{version1: "1.1.1_8", version2: "1.1.1"},
			want: -1,
		},
		{
			name: "partial semver differences 8",
			args: args{version1: "1.1.1_8", version2: "1.1.1_9"},
			want: 1,
		},
		{
			name: "partial semver differences 9",
			args: args{version1: "1.1.1_10", version2: "1.1.1_9"},
			want: -1,
		},
		{
			name: "partial semver differences 10",
			args: args{version1: "1.1.1b", version2: "1.1.1a"},
			want: -1,
		},
		{
			name: "partial semver differences 11",
			args: args{version1: "1.1.1a", version2: "1.1.1b"},
			want: 1,
		},
		{
			name: "partial semver differences 12",
			args: args{version1: "1.1a2", version2: "1.1a3"},
			want: 1,
		},
		{
			name: "partial semver differences 13",
			args: args{version1: "1.1.1a10", version2: "1.1.1b1"},
			want: 1,
		},

		{
			name: "dash-style versions 1",
			args: args{version1: "1.8.9", version2: "1.8.5-1"},
			want: -1,
		},
		{
			name: "dash-style versions 2",
			args: args{version1: "7.0.4-9", version2: "7.0.4-10"},
			want: 1,
		},
		{
			name: "dash-style versions 3",
			args: args{version1: "7.0.4-9", version2: "7.0.4-8"},
			want: -1,
		},
		{
			name: "dash-style versions 4",
			args: args{version1: "2019-01-01", version2: "2019-01-02"},
			want: 1,
		},
		{
			name: "dash-style versions 5",
			args: args{version1: "2019-01-02", version2: "2019-01-01"},
			want: -1,
		},
		{
			name: "dash-style versions 6",
			args: args{version1: "2018-01-01", version2: "2019-01-01"},
			want: 1,
		},
		{
			name: "dash-style versions 7",
			args: args{version1: "2019-01-01", version2: "2018-01-01"},
			want: -1,
		},

		{
			name: "post-release tagging 1",
			args: args{version1: "1", version2: "1+hotfix.0"},
			want: 1,
		},
		{
			name: "post-release tagging 2",
			args: args{version1: "1.0.0", version2: "1.0.0+hotfix.0"},
			want: 1,
		},
		{
			name: "post-release tagging 3",
			args: args{version1: "1.0.0+hotfix.0", version2: "1.0.0+hotfix.1"},
			want: 1,
		},
		{
			name: "post-release tagging 4",
			args: args{version1: "1.0.0+hotfix.1", version2: "1.0.1"},
			want: 1,
		},
		{
			name: "post-release tagging 5",
			args: args{version1: "1.0.0+1.1", version2: "1.0.0+1"},
			want: -1,
		},

		{
			name: "plain text string 1",
			args: args{version1: "latest", version2: "20150405"},
			want: -1,
		},
		{
			name: "plain text string 2",
			args: args{version1: "0.5alpha", version2: "0.5"},
			want: 1,
		},
		{
			name: "plain text string 3",
			args: args{version1: "0.5", version2: "0.5Beta"},
			want: -1,
		},
		{
			name: "plain text string 4",
			args: args{version1: "0.4", version2: "0.5Beta"},
			want: 1,
		},

		{
			name: "empty string",
			args: args{version1: "7.0.4-9", version2: ""},
			want: -1,
		},

		{
			name: "equal versions 1",
			args: args{version1: "12.0", version2: "12.0"},
			want: 0,
		},
		{
			name: "equal versions 2",
			args: args{version1: "7.0.4-9", version2: "7.0.4-9"},
			want: 0,
		},
		{
			name: "equal versions 3",
			args: args{version1: "nightly-20190801", version2: "nightly"},
			want: 0,
		},
		{
			name: "equal versions 4",
			args: args{version1: "nightly-20190801", version2: "nightly-20200801"},
			want: 0,
		},
	}
	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			if got := CompareVersions(test.args.version1, test.args.version2); got != test.want {
				t.Errorf("CompareVersions() = %v, want %v", got, test.want)
			}
		})
	}
}
