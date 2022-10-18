package sqlstorage

func (s *Storage) AddBanner(description string) error {
	query := `insert into banners(description) values($1)`

	_, err := s.db.ExecContext(s.ctx, query, description)
	return err
}

func (s *Storage) RemoveBannerByID(bannerID int) error {
	query := "delete from banners where id = $1"

	_, err := s.db.ExecContext(s.ctx, query, bannerID)
	return err
}
