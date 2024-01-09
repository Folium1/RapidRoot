package rapidroot

import (
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
	if fn == nil {
		return "nil"
	}
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
		return "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return path
}

func isDynamicSegment(segment string) bool {
	return strings.HasPrefix(segment, "$")
}

func isLastSegment(index int, segments []string) bool {
	return index == len(segments)-1
}
