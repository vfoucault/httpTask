package main

import (
	"taskReceiver/internal/config"
	"taskReceiver/internal/manager"
)

func main() {
	cfg := config.New()

	//logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mgr := manager.NewManager(cfg)
	err := mgr.Start()
	if err != nil {
		panic(err)
	}
}
