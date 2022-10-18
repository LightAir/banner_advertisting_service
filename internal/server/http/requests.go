package http

type BaseAdminDescriptionRequest struct {
	Description string
}

type BannerSlotRequest struct {
	BannerID int
	SlotID   int
}

type TrackRequest struct {
	SlotID    int
	BannerID  int
	SDGroupID int
}
