// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import "strings"

// fileBase gets the base of the file without the extension
// path.Base doesn't work because Windows
func fileBase(path string) string {
	if path == "" {
		return ""
	}
	// Strip trailing slashes.
	for len(path) > 0 && (path[len(path)-1] == '/' || path[len(path)-1] == '\'') {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	// Find the last element on Windows
	if i := strings.LastIndex(path, "\\"); i >= 0 {
		path = path[i+1:]
	}
	// If empty now, it had only slashes.
	if path == "" {
		return ""
	}
	// Now remove the extension
	if i := strings.Index(path, "."); i >= 0 {
		return path[:i]
	}
	return path
}
