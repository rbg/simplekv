package main

import (
	"github.com/rbg/simplekv/api"
	"github.com/stackengine/selog"
)

func main() {
	selog.setLevel("all", selog.Debug)
	be := NewMem()
	kv_api := simplekv.NewServer(be)
	kv_api.Run()
}
