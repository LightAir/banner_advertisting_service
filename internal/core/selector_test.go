package core

import (
	"context"
	"testing"

	"github.com/LightAir/bas/internal/config"
	"github.com/LightAir/bas/internal/storage"
	memorystorage "github.com/LightAir/bas/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type Log struct{}

func (l Log) Debugf(format string, args ...interface{}) {}
func (l Log) Errorf(format string, args ...interface{}) {}

type QueueStab struct{}

func (q QueueStab) Connect(ctx context.Context) error {
	return nil
}

func (q QueueStab) Receive(name string, callback func(body []byte)) error {
	return nil
}

func (q QueueStab) Sent(body []byte, name string) error {
	return nil
}

func (q QueueStab) Close() error {
	return nil
}

func TestStorage(t *testing.T) {
	t.Run("selector tests", func(t *testing.T) {
		cfg := &config.Config{}
		logg := &Log{}
		stor := memorystorage.New(logg)
		q := &QueueStab{}

		app := NewApp(stor, cfg, q, logg)

		tracks := map[int]*storage.Tracker{
			1: {
				ID:        1,
				BannerID:  1,
				SlotID:    1,
				SDGroupID: 1,
				Clicks:    0,
				Views:     10,
			},
			2: {
				ID:        1,
				BannerID:  2,
				SlotID:    1,
				SDGroupID: 1,
				Clicks:    0,
				Views:     150,
			},
			3: {
				ID:        1,
				BannerID:  3,
				SlotID:    1,
				SDGroupID: 1,
				Clicks:    5,
				Views:     5,
			},
		}

		stat := make(map[int]int)

		for i := 0; i < 1000; i++ {
			bannerID := app.SelectBanner(tracks)
			if val, isExist := stat[bannerID]; isExist {
				stat[bannerID] = val + 1
			} else {
				stat[bannerID] = 1
			}
		}

		t.Logf("Banners showed. #1: %d, #2: %d, #3: %d", stat[1], stat[2], stat[3])

		require.Greater(t, stat[1], stat[2])
		require.Greater(t, stat[3], stat[2]*6)
		require.Greater(t, stat[3], stat[1]*3)
	})
}
