package main

import (
	"fmt"

	_ "github.com/varunamachi/vaali/vapp"
	_ "github.com/varunamachi/vaali/vcmn"
	_ "github.com/varunamachi/vaali/vevt"
	_ "github.com/varunamachi/vaali/vlog"
	_ "github.com/varunamachi/vaali/vmgo"
	_ "github.com/varunamachi/vaali/vnet"
	_ "github.com/varunamachi/vaali/vpg"
	_ "github.com/varunamachi/vaali/vsec"
	_ "github.com/varunamachi/vaali/vuman"
)

func main() {
	fmt.Println("Vaali...")
}
