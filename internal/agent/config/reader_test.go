package config_test

import (
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/stretchr/testify/assert"
)

func TestFlagReader_ReadFlags(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "",
			cfg:     &config.Config{},
			wantErr: false,
		},
		{
			name:    "config is null should return error",
			cfg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &config.FlagReader{}
			if err := r.ReadFlags(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("FlagReader.ReadFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			assert.Equal(t, "localhost:8080", tt.cfg.Address)
		})
	}
}

func TestEnvReader_ReadEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "nil config, returns error",
			cfg:     nil,
			wantErr: true,
		},
		{
			name:    "config is not nil, returns valid config",
			cfg:     &config.Config{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &config.EnvReader{}
			if err := r.ReadEnvironment(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("EnvReader.ReadEnvironment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
