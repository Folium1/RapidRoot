package rapidRoot

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getFunctionName(fn interface{}) string {
	fullName := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
	return fullName[len(fullName)-1]
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func cleanPath(path string) string {
	if path == "" {
		panic(fmt.Errorf("path can't be empty string"))
	}
	if path[0] != '/' {
		panic(fmt.Errorf("path must start with '/', wrong path: %s", path))
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}
