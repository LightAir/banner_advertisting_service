package memorystorage

import "github.com/LightAir/bas/internal/storage"

func (s *Storage) AddBanner(description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastBannerID++
	s.banners[s.lastBannerID] = &storage.Banner{
		ID:          s.lastBannerID,
		Description: description,
	}

	s.logg.Debugf("add banner: %d", s.lastBannerID)

	return nil
}

func (s *Storage) RemoveBannerByID(bannerID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.banners[bannerID]; !isExist {
		return errRecordNotFound
	}

	delete(s.banners, bannerID)
	return nil
}
