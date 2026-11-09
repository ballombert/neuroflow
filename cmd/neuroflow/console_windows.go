//go:build windows

package main

import (
	"syscall"
)

// initWindowsConsole active la page de code UTF-8 (65001) sur la console Windows.
// Cela garantit un rendu parfait des emojis (🔥, 🧠, 💼, ⚡) et des blocs Unicode (█, ░)
// dans Windows Terminal, PowerShell et CMD.
func initWindowsConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	setConsoleCP := kernel32.NewProc("SetConsoleCP")

	if setConsoleOutputCP.Find() == nil {
		_, _, _ = setConsoleOutputCP.Call(uintptr(65001)) // CP_UTF8
	}
	if setConsoleCP.Find() == nil {
		_, _, _ = setConsoleCP.Call(uintptr(65001)) // CP_UTF8
	}
}
