package appconfig

import (
	"context"
	"fmt"
)

var Config *AppConfig

func Init(configPath string) error {
	var err error
	Config, err = LoadFromPath(context.Background(), configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}
