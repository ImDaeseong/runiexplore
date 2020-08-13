// runChrome
package main

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

type (
	HANDLE uintptr
	HWND   HANDLE
	BOOL   int32
	CSIDL  uint32
)

const (
	SW_HIDE = 0
	SW_SHOW = 5
)

const (
	MAX_PATH = 260

	CSIDL_PROGRAM_FILES    = 0x26
	CSIDL_PROGRAM_FILESX86 = 0x2A
)

var (
	URLS   []string
	nIndex int
)

var (
	shell32 = syscall.NewLazyDLL("shell32.dll")

	procShellExecuteW = shell32.NewProc("ShellExecuteW")

	procSHGetSpecialFolderPathW = shell32.NewProc("SHGetSpecialFolderPathW")
)

func BoolToBOOL(value bool) BOOL {
	if value {
		return 1
	}
	return 0
}

func SHGetSpecialFolderPath(hwndOwner HWND, lpszPath *uint16, csidl CSIDL, fCreate bool) bool {

	ret, _, _ := procSHGetSpecialFolderPathW.Call(uintptr(hwndOwner), uintptr(unsafe.Pointer(lpszPath)), uintptr(csidl), uintptr(BoolToBOOL(fCreate)))

	return ret != 0
}

func getProgramFilesDir() string {

	var buf [MAX_PATH]uint16

	if !SHGetSpecialFolderPath(0, &buf[0], CSIDL_PROGRAM_FILESX86, false) {
		return ""
	}

	return (syscall.UTF16ToString(buf[0:]))
}

func ShellExecute(hwnd HWND, lpOperation string, lpFile string, lpParameters string, lpDirectory string, nShowCmd int) error {

	var ptrlpOperation uintptr
	var ptrlpFile uintptr
	var ptrlpParameters uintptr
	var ptrlpDirectory uintptr

	if len(lpOperation) != 0 {
		ptrlpOperation = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpOperation)))
	} else {
		ptrlpOperation = uintptr(0)
	}

	if len(lpFile) != 0 {
		ptrlpFile = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpFile)))
	} else {
		ptrlpFile = uintptr(0)
	}

	if len(lpParameters) != 0 {
		ptrlpParameters = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpParameters)))
	} else {
		ptrlpParameters = uintptr(0)
	}

	if len(lpDirectory) != 0 {
		ptrlpDirectory = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpDirectory)))
	} else {
		ptrlpDirectory = uintptr(0)
	}

	ret, _, _ := procShellExecuteW.Call(uintptr(hwnd), ptrlpOperation, ptrlpFile, ptrlpParameters, ptrlpDirectory, uintptr(nShowCmd))

	errMsg := ""
	if ret != 0 && ret <= 32 {
		errMsg = "error"
	} else {
		return nil
	}
	return errors.New(errMsg)
}

func main() {

	exePath := fmt.Sprintf("%s\\Google\\Chrome\\Application\\chrome.exe", getProgramFilesDir())

	ShellExecute(0, "open", exePath, "https://github.com/ --no-sandbox", "", SW_SHOW)

}
