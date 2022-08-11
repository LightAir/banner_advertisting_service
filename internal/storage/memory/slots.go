package memorystorage

import "github.com/LightAir/bas/internal/storage"

func (s *Storage) AddSlot(description string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastSlotID++
	s.slots[s.lastSlotID] = &storage.Slot{
		ID:          s.lastSlotID,
		Description: description,
	}

	s.logg.Debugf("add slot: %d", s.lastSlotID)

	return nil
}

func (s *Storage) RemoveSlotByID(slotID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, isExist := s.slots[slotID]; !isExist {
		return errRecordNotFound
	}

	delete(s.slots, slotID)
	return nil
}
