package bucket

import (
	"fmt"
	"strings"
)

type AppRef struct {
	Bucket  string
	Name    string
	Version string
}

func ParseAppString(appString string) (AppRef, error) {
	if appString == "" {
		return AppRef{}, fmt.Errorf("app string is empty")
	}

	appRef := AppRef{}
	unparsedPortion := appString

	if bucket, rest, found := strings.Cut(unparsedPortion, "/"); found {
		appRef.Bucket = bucket
		unparsedPortion = rest
	}

	if name, version, found := strings.Cut(unparsedPortion, "@"); found {
		appRef.Name = name
		appRef.Version = version
	} else {
		appRef.Name = unparsedPortion
	}

	if appRef.Name == "" {
		return AppRef{}, fmt.Errorf("app name not parsed correctly")
	} else if appRef.Bucket == "" && strings.Contains(appString, "/") {
		return AppRef{}, fmt.Errorf("app bucket name not parsed correctly")
	} else if appRef.Version == "" && strings.Contains(appString, "@") {
		return AppRef{}, fmt.Errorf("app version not parsed correctly")
	}

	return appRef, nil
}
