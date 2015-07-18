package main

import (
	"github.com/rbg/simplekv/api"
	"github.com/stackengine/selog"
)

func main() {
	selog.SetLevel("all", selog.Debug)
	kv_api := simplekv.NewServer()
	kv_api.Run()
}
