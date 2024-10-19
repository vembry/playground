package loadtest

import "time"

func setDefault(cfg *Config) {
	// default type to count when undefined
	if cfg.Type == 0 {
		cfg.Type = LoadType_Time
	}

	// default duration
	if cfg.Duration == 0 {
		cfg.Duration = 5 * time.Second
	}

	// default counter
	if cfg.MaxCounter == 0 {
		cfg.MaxCounter = 1000
	}

	// when 'WorkerCount' not defined,
	// then defaults to 1
	if cfg.WorkerCount == 0 {
		cfg.WorkerCount = 1
	}
}
