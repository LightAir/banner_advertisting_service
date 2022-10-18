package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Доступность сервиса
// @Description Проверит доступность сервиса по показу баннерной рекламы
// @Tags Base
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router / [get].
func (s *Server) pingHandler(w http.ResponseWriter, _ *http.Request) {
	s.message(http.StatusOK, "Pong", w)
}

// @Summary Добавить баннер в ротацию
// @Description Добавляет новый баннер в ротацию в данном слоте
// @Tags Admin
// @Produce json
// @Param data body BannerSlotRequest true "BannerID and SlotID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/banner-slot [post].
func (s *Server) addBannerToRotationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to read request body", w)
		return
	}

	data := &BannerSlotRequest{}

	err = json.Unmarshal(body, data)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to unmarshal request body", w)
		return
	}

	err = s.app.AddBannerToSlot(data.BannerID, data.SlotID)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)
		return
	}

	s.message(http.StatusOK, "banner added to slot", w)
}

// @Summary Выбрать баннер для показа
// @Description Возвращает баннер, для показать в указанном слоте для указанной соц-дем. группы
// @Tags Base
// @Param slot_id path int true "Slot ID"
// @Param sd_group_id path int true "Group ID"
// @Success 200 {object} BannerResponse
// @Failure 400
// @Router /api/v1/show-banner/{slot_id}/{sd_group_id} [get].
func (s *Server) showBannerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	slotID, err := strconv.Atoi(vars["slot_id"])
	if err != nil {
		s.logger.Info(err)
		s.message(http.StatusBadRequest, "Bad Slot ID", w)
		return
	}

	sdGroupID, err := strconv.Atoi(vars["sd_group_id"])
	if err != nil {
		s.message(http.StatusBadRequest, "Bad SDGroupID ID", w)
		return
	}

	bannerID, err := s.app.GetBanner(slotID, sdGroupID)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)
		return
	}

	res, err := json.Marshal(BannerResponse{BannerID: bannerID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}
}

// @Summary Удалить баннер из ротации
// @Description Удаляет баннер из ротации в данном слоте
// @Tags Admin
// @Produce json
// @Param data body BannerSlotRequest true "SlotID and BannerID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/banner-slot [delete].
func (s *Server) removeBannerFromRotationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to read request body", w)
		return
	}

	data := &BannerSlotRequest{}

	err = json.Unmarshal(body, data)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to unmarshal request body", w)
		return
	}

	err = s.app.RemoveBannerFromSlot(data.BannerID, data.SlotID)
	if err != nil {
		s.message(http.StatusInternalServerError, "internal server error", w)
		s.logger.Error(err)
		return
	}

	s.message(http.StatusOK, "banner from slot removed", w)
}

// @Summary Засчитать переход
// @Description Увеличивает счетчик переходов на 1 для указанного баннера в данном слоте в указанной группе.
// @Tags Tracker
// @Produce json
// @Param data body TrackRequest true "SlotID, BannerID and SDGroupID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/track [post].
func (s *Server) trackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to read request body", w)
		return
	}

	data := &TrackRequest{}

	err = json.Unmarshal(body, data)
	if err != nil {
		s.message(http.StatusBadRequest, "failed to unmarshal request body", w)
		return
	}

	err = s.app.Track(data.BannerID, data.SlotID, data.SDGroupID)
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.message(http.StatusOK, "tracked", w)
}

// @Summary Добавить баннер
// @Description Добавляет новый баннер
// @Tags Admin Banners
// @Produce json
// @Param data body BaseAdminDescriptionRequest true "Banner description"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/banner [post].
func (s *Server) addBannerHandler(w http.ResponseWriter, r *http.Request) {
	data := &BaseAdminDescriptionRequest{}
	err := s.baseAdminRequest(r, data)
	if err != nil {
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	err = s.app.AddBanner(data.Description)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Banner added", w)
}

// @Summary Удалить баннер
// @Description Удаляет баннер
// @Tags Admin Banners
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/banner/{id} [delete].
func (s *Server) removeBannerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.message(http.StatusBadRequest, "Bad banner ID", w)
		return
	}

	err = s.app.RemoveBanner(id)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Banner removed", w)
}

// @Summary Добавить слот
// @Description Добавляет новый слот
// @Tags Admin Slots
// @Produce json
// @Param data body BaseAdminDescriptionRequest true "Slot description"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/slot [post].
func (s *Server) addSlotHandler(w http.ResponseWriter, r *http.Request) {
	data := &BaseAdminDescriptionRequest{}
	err := s.baseAdminRequest(r, data)
	if err != nil {
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	err = s.app.AddSlot(data.Description)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Slot added", w)
}

// @Summary Удалить слот
// @Description Удаляет слот
// @Tags Admin Slots
// @Produce json
// @Param id path int true "Slot ID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/slot/{id} [delete].
func (s *Server) removeSlotHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.message(http.StatusBadRequest, "Bad banner ID", w)
		return
	}

	err = s.app.RemoveSlot(id)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Slot removed", w)
}

// @Summary Добавить группу
// @Description Добавляет новую Социал-демографическую группу
// @Tags Admin Groups
// @Produce json
// @Param data body BaseAdminDescriptionRequest true "SDGroupID description"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/group [post].
func (s *Server) addSDGroupHandler(w http.ResponseWriter, r *http.Request) {
	data := &BaseAdminDescriptionRequest{}
	err := s.baseAdminRequest(r, data)
	if err != nil {
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	err = s.app.AddSDGroup(data.Description)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Group added", w)
}

// @Summary Удалить группу
// @Description Удаляет Социал-демографическую группу
// @Tags Admin Groups
// @Produce json
// @Param id path int true "SDGroupID ID"
// @Success 200 {object} TypicalResponse
// @Failure 400
// @Router /api/v1/group/{id} [delete].
func (s *Server) removeSDGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		s.message(http.StatusBadRequest, "Bad banner ID", w)
		return
	}

	err = s.app.RemoveSlot(id)
	if err != nil {
		s.logger.Error(err)
		s.message(http.StatusBadRequest, "internal server error", w)
		return
	}

	s.message(http.StatusOK, "Slot removed", w)
}
