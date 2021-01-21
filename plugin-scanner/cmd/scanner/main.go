package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"plugin-scanner/scanner"
)


const pluginsDir = "../../plugins/"

var (
	files []os.FileInfo
	err error
	p *plugin.Plugin
	n plugin.Symbol
	check scanner.Checker
	res *scanner.Result
)

func main() {
	if files, err = ioutil.ReadDir(pluginsDir); err != nil {
		log.Fatalln(err)
	}

	for k := range files {
		fmt.Println("Found plugin: " + files[k].Name())
		if p, err = plugin.Open(pluginsDir + files[k].Name()); err != nil {
			log.Fatalln(err)
		}
		if n, err = p.Lookup("New"); err != nil {
			log.Fatalln(err)
		}
		checkFunc, ok := n.(func() scanner.Checker)
		if !ok {
			log.Fatalln("Plugin entry point is no good. Expecting: func New() scanner.Checker{ ... }")
		}
		check = checkFunc()
		res = check.Check("127.0.0.1", 8080)
		if res.Vulnerable {
			fmt.Printf("Host is vulnerable: %s\n", res.Details)
		} else {
			fmt.Printf("Host is not vulnerable\n")
		}
	}
}