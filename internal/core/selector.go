package core

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/LightAir/bas/internal/core/algs"
	"github.com/LightAir/bas/internal/storage"
)

func randInt(min, max int) int {
	return int(math.Round(float64(min) + rand.Float64()*(float64(max-min)))) //nolint
}

func getTotalShows(allTracks map[int]*storage.Tracker) int {
	totalShows := 0
	for _, track := range allTracks {
		totalShows += track.Views
	}

	return totalShows
}

func (a App) SelectBanner(allTracks map[int]*storage.Tracker) int {
	totalTracks := len(allTracks)

	if totalTracks <= 0 {
		return 0
	}

	rand.Seed(time.Now().UnixNano())

	totalShows := getTotalShows(allTracks)

	banners := make([]BannerWithWeight, 0)

	weightSum := 0.0
	for _, track := range allTracks {
		weight := algs.Ucb1(track.Clicks, totalShows, track.Views)
		banners = append(banners, BannerWithWeight{
			bannerID: track.BannerID,
			weight:   weight,
		})
		weightSum += weight
	}

	sort.Slice(banners, func(i, j int) bool {
		return banners[i].weight < banners[j].weight
	})

	result := 0

	rndExcept := rand.Intn(99) + 1 //nolint
	if rndExcept <= a.config.PercentExclude {
		rndBanner := rand.Intn(totalTracks) //nolint
		result = banners[rndBanner].bannerID
	} else {
		randomWeight := float64(randInt(1, int(math.Floor(weightSum))))
		estimatedWeight := 0.0
		for id, banner := range banners {
			estimatedWeight += banner.weight
			if estimatedWeight >= randomWeight {
				result = banners[id].bannerID
				break
			}
		}
	}

	return result
}
