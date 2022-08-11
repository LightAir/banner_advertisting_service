package sqlstorage

import (
	"github.com/LightAir/bas/internal/storage"
)

func (s *Storage) TrackClick(bannerID, slotID, sdGroupID int) error {
	query := "update tracks set clicks = clicks + 1 where banner_id = $1 and slot_id = $2 and sd_group_id = $3"

	_, err := s.db.ExecContext(s.ctx, query, bannerID, slotID, sdGroupID)
	return err
}

func (s *Storage) TrackView(bannerID, slotID, sdGroupID int) error {
	query := "update tracks set views = tracks.views + 1 where banner_id = $1 and slot_id = $2 and sd_group_id = $3"

	_, err := s.db.ExecContext(s.ctx, query, bannerID, slotID, sdGroupID)
	return err
}

func (s *Storage) GetAllTracks(slotID, sdGroupID int) (map[int]*storage.Tracker, error) {
	query := `select 
    			id,
    			banner_id,
    			slot_id,
    			sd_group_id,
    			clicks,
    			views
			  from
			    tracks
			  where
			    slot_id = $1
			  and
			    sd_group_id = $2
			 `

	rows, err := s.db.QueryxContext(s.ctx, query, slotID, sdGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make(map[int]*storage.Tracker)

	for rows.Next() {
		var track storage.Tracker
		err := rows.StructScan(&track)
		if err != nil {
			return nil, err
		}
		tracks[track.ID] = &track
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s *Storage) AddTrack(bannerID, slotID, sdGroupID int) error {
	query := "insert into tracks(banner_id, slot_id, sd_group_id, views, clicks) values($1, $2, $3, 0, 0)"

	_, err := s.db.ExecContext(s.ctx, query, bannerID, slotID, sdGroupID)
	return err
}

func (s *Storage) RemoveTracks(bannerID, slotID int) error {
	query := "delete from tracks where banner_id = $1 and slot_id = $2"

	_, err := s.db.ExecContext(s.ctx, query, bannerID, slotID)
	return err
}
