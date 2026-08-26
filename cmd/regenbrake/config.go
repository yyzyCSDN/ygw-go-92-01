package main

import "flag"

type Config struct {
	Addr      string
	DataDir   string
	Threshold float64
}

func parseConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Addr, "addr", "127.0.0.1:8090", "listen address")
	flag.StringVar(&cfg.DataDir, "dir", "./data", "data directory for run records")
	flag.Float64Var(&cfg.Threshold, "threshold", 900.0, "voltage threshold to trigger absorption")
	flag.Parse()
	return cfg
}
