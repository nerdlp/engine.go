package engine

import "strings"

func ComputePath(path string) string {
	return "/" + strings.Trim(path, "/") + "/"
}
