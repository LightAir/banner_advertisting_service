package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/LightAir/bas/docs"
	"github.com/LightAir/bas/internal/config"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Logger interface {
	Error(msg ...interface{})
	Info(msg ...interface{})
}

type Server struct {
	host   string
	port   string
	logger Logger
	app    Application
	server *http.Server
	cfg    *config.Config
}

type Application interface {
	AddBanner(description string) error
	RemoveBanner(id int) error
	AddSlot(description string) error
	RemoveSlot(id int) error
	AddSDGroup(description string) error
	RemoveSDGroup(id int) error
	AddBannerToSlot(bannerID, slotID int) error
	GetBanner(slotID, sdGroupID int) (int, error)
	RemoveBannerFromSlot(bannerID, slotID int) error
	Track(bannerID, slotID, sdGroupID int) error
}

func NewServer(logger Logger, app Application, cfg *config.Config) *Server {
	return &Server{
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		logger: logger,
		app:    app,
		cfg:    cfg,
	}
}

func (s *Server) message(status int, message string, w http.ResponseWriter) {
	res, err := json.Marshal(TypicalResponse{
		Status:  status,
		Message: message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}

	w.WriteHeader(status)

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error(err)
	}
}

func (s *Server) baseAdminRequest(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	return nil
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.host, s.port)

	r := mux.NewRouter()
	r.HandleFunc("/", s.pingHandler)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/banner", s.addBannerHandler).Methods("POST")
	api.HandleFunc("/banner/{id}", s.removeBannerHandler).Methods("DELETE")

	api.HandleFunc("/slot", s.addSlotHandler).Methods("POST")
	api.HandleFunc("/slot/{id}", s.removeSlotHandler).Methods("DELETE")

	api.HandleFunc("/group", s.addSDGroupHandler).Methods("POST")
	api.HandleFunc("/group/{id}", s.removeSDGroupHandler).Methods("DELETE")

	api.HandleFunc("/banner-slot", s.addBannerToRotationHandler).Methods("POST")
	api.HandleFunc("/banner-slot", s.removeBannerFromRotationHandler).Methods("DELETE")

	api.HandleFunc("/show-banner/{slot_id}/{sd_group_id}", s.showBannerHandler)

	api.HandleFunc("/track", s.trackHandler).Methods("POST")

	if s.cfg.Environment != config.EnvProd {
		host := s.host
		if host == "" {
			host = "localhost"
		}
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, s.port)
		r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	s.server = server

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			s.logger.Error(err)
		}
	}()

	s.logger.Info("BAS started on " + addr)

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
