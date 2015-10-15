package main

import (
	"github.com/rbg/simplekv/api"
	"github.com/rbg/simplekv/store"
	"github.com/stackengine/selog"
)

func main() {
	selog.SetLevel("all", selog.Debug)
	be := store.NewRedis()
	kv_api := simplekv.NewServer(be)
	kv_api.Run()
}
