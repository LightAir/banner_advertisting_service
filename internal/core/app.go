package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LightAir/bas/internal/config"
	"github.com/LightAir/bas/internal/queue"
	"github.com/LightAir/bas/internal/storage"
)

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	AddBanner(description string) error
	RemoveBannerByID(id int) error
	AddSlot(description string) error
	RemoveSlotByID(id int) error
	AddSDGroup(description string) error
	RemoveSDGroupByID(id int) error
	GetAllGroups() ([]*storage.SDGroup, error)
	GetAllTracks(slotID, sdGroupID int) (map[int]*storage.Tracker, error)
	TrackClick(bannerID, slotID, sdGroupID int) error
	TrackView(bannerID, slotID, sdGroupID int) error
	AddTrack(bannerID, slotID, sdGroupID int) error
	RemoveTracks(bannerID, slotID int) error
}

type App struct {
	config  config.Config
	storage Storage
	queue   queue.Queue
	logg    Logger
}

type Event struct {
	Type      string
	SlotID    int
	BannerID  int
	SDGroupID int
	DateTime  time.Time
}

type BannerWithWeight struct {
	bannerID int
	weight   float64
}

type Logger interface {
	Errorf(format string, args ...interface{})
}

func NewApp(storage Storage, cfg *config.Config, queue queue.Queue, logg Logger) *App {
	return &App{
		config:  *cfg,
		storage: storage,
		queue:   queue,
		logg:    logg,
	}
}

func (a App) AddBanner(description string) error {
	return a.storage.AddBanner(description)
}

func (a App) RemoveBanner(id int) error {
	return a.storage.RemoveBannerByID(id)
}

func (a App) AddSlot(description string) error {
	return a.storage.AddSlot(description)
}

func (a App) RemoveSlot(id int) error {
	return a.storage.RemoveSlotByID(id)
}

func (a App) AddSDGroup(description string) error {
	return a.storage.AddSDGroup(description)
}

func (a App) RemoveSDGroup(id int) error {
	return a.storage.RemoveSDGroupByID(id)
}

func (a App) AddBannerToSlot(bannerID, slotID int) error {
	allGroups, err := a.storage.GetAllGroups()
	if err != nil {
		return err
	}

	for _, group := range allGroups {
		if err := a.storage.AddTrack(bannerID, slotID, group.ID); err != nil {
			return err
		}
	}

	return nil
}

func (a App) sendEventToQueue(eventType string, bannerID, slotID, sdGroupID int) {
	e := Event{
		Type:      eventType,
		SlotID:    slotID,
		BannerID:  bannerID,
		SDGroupID: sdGroupID,
		DateTime:  time.Now(),
	}

	body, err := json.Marshal(e)
	if err != nil {
		a.logg.Errorf("unmarshal track error: %w", err)
	}

	if err := a.queue.Sent(body, "tracker"); err != nil {
		a.logg.Errorf("sent error: %w", err)
	}
}

func (a App) GetBanner(slotID, sdGroupID int) (int, error) {
	allTracks, err := a.storage.GetAllTracks(slotID, sdGroupID)
	if err != nil {
		return 0, err
	}

	bannerID := a.SelectBanner(allTracks)

	err = a.storage.TrackView(bannerID, slotID, sdGroupID)
	if err != nil {
		return 0, err
	}

	a.sendEventToQueue("view", bannerID, slotID, sdGroupID)

	return bannerID, nil
}

func (a App) RemoveBannerFromSlot(bannerID, slotID int) error {
	return a.storage.RemoveTracks(bannerID, slotID)
}

func (a App) Track(bannerID, slotID, sdGroupID int) error {
	a.sendEventToQueue("click", bannerID, slotID, sdGroupID)

	return a.storage.TrackClick(bannerID, slotID, sdGroupID)
}
