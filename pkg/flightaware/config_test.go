package flightaware_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jalavosus/stuffnotifier/pkg/flightaware"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		confPath string
	}

	var (
		wantConf1 = &flightaware.Config{
			PollInterval: 2 * time.Minute,
			Auth:         nil,
			Notifications: flightaware.NotificationsConfig{
				UseLocalTime:  true,
				GateDeparture: false,
				Takeoff:       true,
				Landing:       true,
				GateArrival:   true,
				BaggageClaim:  false,
				PreArrival: flightaware.PreEventConfig{
					Enabled:   true,
					Estimated: true,
					Scheduled: false,
					Offset:    30 * time.Minute,
				},
			},
		}

		wantConf2 = &flightaware.Config{
			PollInterval: 2 * time.Minute,
			Auth: &flightaware.AuthConfig{
				ApiKey: "1234567",
			},
			Notifications: flightaware.NotificationsConfig{
				UseLocalTime:  false,
				GateDeparture: true,
				GateArrival:   true,
				BaggageClaim:  true,
				PreArrival: flightaware.PreEventConfig{
					Enabled: false,
				},
			},
		}
	)

	tests := []struct {
		wantConf *flightaware.Config
		name     string
		args     args
		wantErr  bool
	}{
		{
			name:     "test config 1 (toml)",
			args:     args{"testdata/configs/flightaware_1.toml"},
			wantConf: wantConf1,
			wantErr:  false,
		},
		{
			name:     "test config 1 (yaml)",
			args:     args{"testdata/configs/flightaware_1.yaml"},
			wantConf: wantConf1,
			wantErr:  false,
		},
		{
			name:     "test config 2 (toml)",
			args:     args{"testdata/configs/flightaware_2.toml"},
			wantConf: wantConf2,
			wantErr:  false,
		},
		{
			name:     "test config 2 (yaml)",
			args:     args{"testdata/configs/flightaware_2.yaml"},
			wantConf: wantConf2,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flightaware.LoadConfig(tt.args.confPath)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)

			if tt.wantConf.Auth == nil {
				assert.Nil(t, got.Auth)
			} else {
				assert.NotNil(t, got.Auth)
				assert.Equal(t, tt.wantConf.Auth.ApiKey, got.Auth.ApiKey)
			}

			wantNotifsConf := tt.wantConf.Notifications
			gotNotifsConf := got.Notifications

			assert.Equalf(
				t, tt.wantConf.PollInterval, got.PollInterval,
				"expected PollInterval to be %[1]s seconds, got %[2]s seconds",
				formatDurationSeconds(t, tt.wantConf.PollInterval),
				formatDurationSeconds(t, got.PollInterval),
			)

			assert.Equalf(
				t, wantNotifsConf.PreArrival.Enabled, gotNotifsConf.PreArrival.Enabled,
				"expected Notifications.PreArrival.Enabled to be %[1]t, got %[2]t",
				wantNotifsConf.PreArrival.Enabled, gotNotifsConf.PreArrival.Enabled,
			)

			assert.Equalf(
				t, wantNotifsConf.BaggageClaim, gotNotifsConf.BaggageClaim,
				"expected Notifications.BaggageClaim to be %[1]t, got %[2]t",
				wantNotifsConf.BaggageClaim, gotNotifsConf.BaggageClaim,
			)

			assert.Equalf(
				t, wantNotifsConf.UseLocalTime, gotNotifsConf.UseLocalTime,
				"expected Notifications.UseLocalTime to be %[1]t, got %[2]t",
				wantNotifsConf.UseLocalTime, gotNotifsConf.UseLocalTime,
			)
		})
	}
}

func formatDurationSeconds(t *testing.T, d time.Duration) string {
	t.Helper()
	return strconv.FormatFloat(d.Seconds(), 'f', -1, 64)
}
