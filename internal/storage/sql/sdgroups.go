package sqlstorage

import "github.com/LightAir/bas/internal/storage"

func (s *Storage) AddSDGroup(description string) error {
	query := `insert into sd_groups(description) values($1)`

	_, err := s.db.ExecContext(s.ctx, query, description)
	return err
}

func (s *Storage) GetAllGroups() ([]*storage.SDGroup, error) {
	query := `select id, description from sd_groups`

	rows, err := s.db.QueryxContext(s.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*storage.SDGroup, 0)

	for rows.Next() {
		var group storage.SDGroup
		err := rows.StructScan(&group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *Storage) RemoveSDGroupByID(sdGroupID int) error {
	query := "delete from sd_groups where id = $1"

	_, err := s.db.ExecContext(s.ctx, query, sdGroupID)
	return err
}
