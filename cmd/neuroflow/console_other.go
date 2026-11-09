//go:build !windows

package main

// initWindowsConsole est un no-op sur les systèmes Unix/macOS/Linux.
func initWindowsConsole() {}
