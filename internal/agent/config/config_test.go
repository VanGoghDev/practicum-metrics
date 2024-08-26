package config_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	mock_config "github.com/VanGoghDev/practicum-metrics/internal/agent/config/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		wantConfig  *config.Config
		wantFlagErr bool
		wantEnvErr  bool
	}{
		{
			name:        "environment reader returns error",
			wantEnvErr:  true,
			wantFlagErr: false,
		},
		{
			name:        "flag reader returns error",
			wantEnvErr:  false,
			wantFlagErr: true,
		},
		{
			name:        "reader returns config",
			wantFlagErr: false,
			wantEnvErr:  false,
			wantConfig:  &config.Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			er := mock_config.NewMockEnvironmentReader(ctrl)
			fr := mock_config.NewMockFlagsReader(ctrl)

			if tt.wantEnvErr {
				er.EXPECT().ReadEnvironment(gomock.Any()).Return(errors.New("test error")).AnyTimes()
			} else {
				er.EXPECT().ReadEnvironment(gomock.Any()).Return(nil).AnyTimes()
			}

			if tt.wantFlagErr {
				fr.EXPECT().ReadFlags(gomock.Any()).Return(errors.New("test error")).AnyTimes()
			} else {
				fr.EXPECT().ReadFlags(gomock.Any()).Return(nil).AnyTimes()
			}

			gotConfig, err := config.Load(er, fr)
			if (err != nil) != (tt.wantEnvErr || tt.wantFlagErr) {
				t.Errorf("Load() error = %v", err)
				return
			}
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("Load() = %v, want %v", gotConfig, tt.wantConfig)
			}
			if !tt.wantEnvErr && !tt.wantFlagErr {
				assert.Equal(t, tt.wantConfig.Address, gotConfig.Address)
			}
		})
	}
}
