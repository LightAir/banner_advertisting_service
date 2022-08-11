//go:build integration
// +build integration

package tests

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/LightAir/bas/internal/config"
	http2 "github.com/LightAir/bas/internal/server/http"
	"github.com/stretchr/testify/suite"
)

type httpTestSuite struct {
	suite.Suite
	host string
	port string
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/banner/config-test.yaml", "Path to configuration file")
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, &httpTestSuite{})
}

func (s *httpTestSuite) SetupSuite() {
	flag.Parse()

	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	s.host = cfg.Server.Host
	s.port = cfg.Server.Port
}

func (s *httpTestSuite) req(method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%s/%s", s.host, s.port, path), body)
	s.NoError(err)

	client := http.Client{}

	response, err := client.Do(req)
	s.NoError(err)

	return response
}

func (s *httpTestSuite) getBody(response *http.Response) string {
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	return strings.Trim(string(byteBody), "\n")
}

func (s *httpTestSuite) TestPing() {
	response := s.req(http.MethodGet, "", nil)
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"Pong"}`, s.getBody(response))
}

func (s *httpTestSuite) addBanner(description string) {
	str := fmt.Sprintf(`{"description": "%s"}`, description)
	response := s.req(http.MethodPost, "api/v1/banner", strings.NewReader(str))
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"Banner added"}`, s.getBody(response))
}

func (s *httpTestSuite) addBannerToSlot(bannerID, slotID int) {
	str := fmt.Sprintf(`{"bannerId": %d, "slotId": %d}`, bannerID, slotID)
	response := s.req(http.MethodPost, "api/v1/banner-slot", strings.NewReader(str))
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"banner added to slot"}`, s.getBody(response))
}

func (s *httpTestSuite) trackBanner(bannerID, slotID, groupID int) {
	str := fmt.Sprintf(`{"bannerId": %d, "slotId": %d, "sdgroupID": %d}`, bannerID, slotID, groupID)
	response := s.req(http.MethodPost, "api/v1/track", strings.NewReader(str))
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"tracked"}`, s.getBody(response))
}

func (s *httpTestSuite) viewBanner(groupID, slotID int) int {
	path := fmt.Sprintf(`api/v1/show-banner/%d/%d`, groupID, slotID)
	response := s.req(http.MethodGet, path, nil)
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	data := http2.BannerResponse{}
	byteBody, err := ioutil.ReadAll(response.Body)
	s.NoError(err)

	err = json.Unmarshal(byteBody, &data)
	s.NoError(err)

	return data.BannerID
}

func (s *httpTestSuite) clickBanner(bannerID, slotID, groupID int) {
	str := fmt.Sprintf(`{"bannerId": %d, "slotId": %d}`, bannerID, slotID)
	response := s.req(http.MethodPost, "api/v1/banner-slot", strings.NewReader(str))
	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"banner added to slot"}`, s.getBody(response))
}

func (s *httpTestSuite) TestFull() {
	// add banners
	s.addBanner("first banner")
	s.addBanner("second banner")
	s.addBanner("third banner")

	// add group
	response := s.req(http.MethodPost, "api/v1/group", strings.NewReader(`{"description": "my social group"}`))
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"Group added"}`, s.getBody(response))
	response.Body.Close()

	// add slot
	response = s.req(http.MethodPost, "api/v1/slot", strings.NewReader(`{"description": "first slot"}`))
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(`{"Status":200,"Message":"Slot added"}`, s.getBody(response))
	response.Body.Close()

	// add banners to slot
	s.addBannerToSlot(1, 1)
	s.addBannerToSlot(2, 1)
	s.addBannerToSlot(3, 1)

	for i := 0; i < 1000; i++ {
		s.viewBanner(1, 1)
	}

	for i := 0; i < 100; i++ {
		s.trackBanner(3, 1, 1)
	}
}
