package memorystorage

import "github.com/LightAir/bas/internal/storage"

func (s *Storage) AddSDGroup(description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastSDGroupID++
	s.sdGroups[s.lastSDGroupID] = &storage.SDGroup{
		ID:          s.lastSDGroupID,
		Description: description,
	}

	s.logg.Debugf("add group: %d", s.lastSDGroupID)

	return nil
}

func (s *Storage) RemoveSDGroupByID(sdGroupID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.sdGroups[sdGroupID]; !isExist {
		return errRecordNotFound
	}

	delete(s.sdGroups, sdGroupID)
	return nil
}

func (s *Storage) GetAllGroups() ([]*storage.SDGroup, error) {
	groups := make([]*storage.SDGroup, 0)
	for _, group := range s.sdGroups {
		groups = append(groups, group)
	}

	return groups, nil
}
