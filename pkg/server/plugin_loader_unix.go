//go:build !windows
// +build !windows

package server

/*
#include <dlfcn.h>
#include <stdlib.h>

typedef void* (*CreateFilterFunc)();

void* loadLibrary(const char* path) {
    return dlopen(path, RTLD_LAZY);
}

void* getSymbol(void* handle, const char* symbol) {
    return dlsym(handle, symbol);
}

const char* getError() {
    return dlerror();
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func loadPluginFilter(pluginPath string) (FilterFactory, error) {
	cPath := C.CString(pluginPath)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.loadLibrary(cPath)
	if handle == nil {
		return nil, fmt.Errorf("failed to load plugin: %s", C.GoString(C.getError()))
	}

	symbol := C.CString("CreateFilter")
	defer C.free(unsafe.Pointer(symbol))

	createFunc := C.getSymbol(handle, symbol)
	if createFunc == nil {
		return nil, fmt.Errorf("failed to find CreateFilter symbol: %s", C.GoString(C.getError()))
	}

	createFilter := *(*func() FilterFactory)(unsafe.Pointer(&createFunc))
	return createFilter(), nil
}
