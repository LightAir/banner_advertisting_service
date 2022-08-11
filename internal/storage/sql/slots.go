package sqlstorage

func (s *Storage) AddSlot(description string) error {
	query := `insert into slots(description) values($1)`

	_, err := s.db.ExecContext(s.ctx, query, description)
	return err
}

func (s *Storage) RemoveSlotByID(slotID int) error {
	query := "delete from slots where id = $1"

	_, err := s.db.ExecContext(s.ctx, query, slotID)
	return err
}
