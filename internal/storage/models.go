package storage

type Banner struct {
	ID          int
	Description string
}

type Slot struct {
	ID          int
	Description string
}

type SDGroup struct {
	ID          int
	Description string
}

type BannerSlot struct {
	ID       int
	BannerID int `db:"banner_id"`
	SlotID   int `db:"slot_id"`
}

type Tracker struct {
	ID        int
	BannerID  int `db:"banner_id"`
	SlotID    int `db:"slot_id"`
	SDGroupID int `db:"sd_group_id"`
	Clicks    int
	Views     int
}
