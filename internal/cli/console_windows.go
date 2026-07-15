//go:build windows

package cli

import "syscall"

const utf8CodePage = 65001

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	setConsoleOutputCP = kernel32.NewProc("SetConsoleOutputCP")
	setConsoleCP       = kernel32.NewProc("SetConsoleCP")
)

func configureConsoleOutput() {
	// Go writes UTF-8 strings. On Windows, older PowerShell/cmd sessions may use
	// a non-UTF-8 console code page, which makes Chinese output unreadable.
	_, _, _ = setConsoleOutputCP.Call(uintptr(utf8CodePage))
	_, _, _ = setConsoleCP.Call(uintptr(utf8CodePage))
}
