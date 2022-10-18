package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/LightAir/bas/internal/storage"
)

var (
	errRecordAlreadyExist = errors.New("record already exist")
	errRecordNotFound     = errors.New("record not found")
)

type Logger interface {
	Debugf(format string, args ...interface{})
}

type (
	BannerSlot map[string]*storage.BannerSlot
	Tracks     map[int]*storage.Tracker
	Banners    map[int]*storage.Banner
	Slots      map[int]*storage.Slot
	SDGroups   map[int]*storage.SDGroup
)

type Storage struct {
	logg          Logger
	tracks        Tracks
	lastTrackID   int
	banners       Banners
	lastBannerID  int
	slots         Slots
	lastSlotID    int
	sdGroups      SDGroups
	lastSDGroupID int
	bannerSlot    BannerSlot
	mu            sync.RWMutex
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tracks = nil
	s.lastTrackID = 0
	s.banners = nil
	s.lastBannerID = 0
	s.slots = nil
	s.lastSlotID = 0
	s.sdGroups = nil
	s.lastSDGroupID = 0
	s.bannerSlot = nil

	return nil
}

func New(logg Logger) *Storage {
	return &Storage{
		logg:          logg,
		tracks:        make(map[int]*storage.Tracker),
		lastTrackID:   0,
		banners:       make(map[int]*storage.Banner),
		lastBannerID:  0,
		slots:         make(map[int]*storage.Slot),
		lastSlotID:    0,
		sdGroups:      make(map[int]*storage.SDGroup),
		lastSDGroupID: 0,
		bannerSlot:    make(map[string]*storage.BannerSlot),
		mu:            sync.RWMutex{},
	}
}
