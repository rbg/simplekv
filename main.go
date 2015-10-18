package main

import (
	"os"

	"github.com/rbg/simplekv/api"
	"github.com/rbg/simplekv/store"
	"github.com/stackengine/selog"

	"github.com/codegangsta/cli"
)

var (
	slog = selog.Register("simplekv", 0)
)

func startServer(c *cli.Context, be store.Store) {
	var kv_api *simplekv.ApiServer

	if ep := c.String("api-endpoint"); len(ep) > 0 {
		kv_api = simplekv.NewServer(be, &ep)
	} else {
		kv_api = simplekv.NewServer(be, nil)
	}

	if kv_api != nil {
		kv_api.Run()
	}
}

func MemKV(c *cli.Context) {
	if c.Bool("dbg") {
		selog.SetLevel("all", selog.Debug)
	} else {
		selog.SetLevel("all", selog.Info)
	}
	if be := store.NewMem(); be != nil {
		startServer(c, be)
	}
}

func RedisKV(c *cli.Context) {
	if c.Bool("dbg") {
		selog.SetLevel("all", selog.Debug)
	} else {
		selog.SetLevel("all", selog.Info)
	}

	if be := store.NewRedis(c.String("redis-write"), c.String("redis-read")); be != nil {
		startServer(c, be)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "simplekv"
	app.Usage = "starts up the simplekv API"
	app.Commands = []cli.Command{
		{
			Name:      "redis",
			ShortName: "r",
			Usage:     "start up a redis backed simplekv",
			Action:    RedisKV,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "dbg, d",
					Usage:  "turn on debugging",
					EnvVar: "",
				},

				cli.StringFlag{
					Name:   "api-endpoint, api",
					EnvVar: "API_ENDPOINT",
					Usage:  "endpoint of where to listen for requests.",
				},

				cli.StringFlag{
					Name:   "redis-write, rw",
					EnvVar: "REDIS_WRITE",
					Usage:  "required end-point where to read/write data.",
				},

				cli.StringFlag{
					Name:   "redis-read, rr",
					EnvVar: "REDIS_READ",
					Usage:  "optional end-point where to read data."},
			}},
		{
			Name:      "mem",
			ShortName: "m",
			Usage:     "start up a memory backed simplekv",
			Action:    MemKV,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dbg, d",
					Usage: "turn on debugging",
				},

				cli.StringFlag{
					Name:   "api-endpoint, api",
					EnvVar: "API_ENDPOINT",
					Usage:  "endpoint of where to listen for requests."},
			}},
	}

	app.Run(os.Args)
}
