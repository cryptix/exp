//go:generate go-bindata -prefix=static static/...
package main

import (
	"fmt"

	"github.com/cryptix/go/logging"
	"github.com/olebedev/go-duktape"
)

func main() {
	vm := duktape.NewContext()
	const stubs = `var self = {}, console = {log:print,warn:print,error:print,info:print}`
	if vm.PevalString(stubs) != 0 {
		panic(vm.SafeToString(-1))
	}
	rt, err := Asset("react-tools.js")
	logging.CheckFatal(err)

	vm.PushString(string(rt))
	vm.PushString("jsx")

	if vm.Pcompile(0) != 0 {
		panic(vm.SafeToString(-1))
	}
	if vm.Pcall(0) != 0 {
		panic(vm.SafeToString(-1))
	}

	// vm.PushGlobalObject()
	if vm.PevalString(`print("fromDuktape");Object.keys(this);Object.keys(self)`) != 0 {
		panic(vm.SafeToString(-1))
	}
	fmt.Println(vm.SafeToString(-1))
}
