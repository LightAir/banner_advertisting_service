package memorystorage

import "github.com/LightAir/bas/internal/storage"

func (s *Storage) filterTracks(bannerID, slotID, sdGroupID int) map[int]*storage.Tracker {
	allTracks := make(map[int]*storage.Tracker)
	s.mu.Lock()
	for key, value := range s.tracks {
		allTracks[key] = value
	}
	s.mu.Unlock()

	if slotID != 0 {
		for id, track := range s.tracks {
			if track.SlotID != slotID {
				delete(allTracks, id)
			}
		}
	}

	if bannerID != 0 {
		for id, track := range s.tracks {
			if track.BannerID != bannerID {
				delete(allTracks, id)
			}
		}
	}

	if sdGroupID != 0 {
		for id, track := range s.tracks {
			if track.SDGroupID != sdGroupID {
				delete(allTracks, id)
			}
		}
	}

	return allTracks
}

func (s *Storage) getTrack(bannerID, slotID, sdGroupID int) *storage.Tracker {
	tracks := s.filterTracks(bannerID, slotID, sdGroupID)

	for _, track := range tracks {
		return track
	}

	return nil
}

func (s *Storage) TrackClick(bannerID, slotID, sdGroupID int) error {
	track := s.getTrack(bannerID, slotID, sdGroupID)

	if track != nil {
		s.mu.Lock()
		t := s.tracks[track.ID]
		t.Clicks = track.Clicks + 1
		s.mu.Unlock()

		return nil
	}

	return errRecordNotFound
}

func (s *Storage) TrackView(bannerID, slotID, sdGroupID int) error {
	track := s.getTrack(bannerID, slotID, sdGroupID)

	if track != nil {
		s.mu.Lock()
		t := s.tracks[track.ID]
		t.Views = track.Views + 1
		s.mu.Unlock()

		return nil
	}

	return errRecordNotFound
}

func (s *Storage) AddTrack(bannerID, slotID, sdGroupID int) error {
	track := s.getTrack(bannerID, slotID, sdGroupID)

	if track == nil {
		s.mu.Lock()
		s.lastTrackID++
		s.tracks[s.lastTrackID] = &storage.Tracker{
			ID:        s.lastTrackID,
			BannerID:  bannerID,
			SlotID:    slotID,
			SDGroupID: sdGroupID,
			Views:     0,
			Clicks:    0,
		}
		s.mu.Unlock()

		return nil
	}

	return errRecordAlreadyExist
}

func (s *Storage) RemoveTracks(bannerID, slotID int) error {
	for sdGroupID := range s.sdGroups {
		track := s.getTrack(bannerID, slotID, sdGroupID)
		if track == nil {
			s.mu.Lock()
			delete(s.tracks, track.ID)
			s.mu.Unlock()

			return nil
		}
	}

	return errRecordNotFound
}

func (s *Storage) GetAllTracks(slotID, sdGroupID int) (map[int]*storage.Tracker, error) {
	track := s.filterTracks(0, slotID, sdGroupID)

	if len(track) == 0 {
		return nil, errRecordNotFound
	}

	return track, nil
}
