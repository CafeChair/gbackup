package main

import (
	"flag"
	"fmt"
	"gbackup/plugin"
)

func main() {
	cfg := flag.String("c", "backup.json", "config file")
	flag.Parse()
	plugin.ParseConfig(*cfg)
	scripts := plugin.ListScripts(plugin.Config().Script.Dir)
	for _, s := range scripts {
		str, err := plugin.RunScript(s)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(str)
	}
}

//后续加循环执行的功能
