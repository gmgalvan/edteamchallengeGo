package transport

import (
	"gmgalvan/edChallenge2021/internal/schema"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

//go:generate mockgen -destination=./mocks/usecases_mock.go -package=mocks gmgalvan/edChallenge2021/internal/transport Usecases
type Usecases interface {
	RetrieveChart(ticker *schema.Ticker) (*schema.Chart, error)
}

type Handler struct {
	uc Usecases
}

func NewTrasport(uc Usecases) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) InitRoutes(e *echo.Echo) {
	// server health
	internal := e.Group("/internal")
	internal.GET("/heartbeat", h.heartBeat)

	// usecases
	nomicsData := e.Group("api/v1/nomics")
	nomicsData.GET("/chart", h.getChart)

}

func (h *Handler) heartBeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func (h *Handler) getChart(c echo.Context) error {
	// get query data
	id := c.QueryParam("id")
	convert := c.QueryParam("convert")
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	// format date interval 20060102 to time.RFC3339
	layout := "20060102"
	startInterval, err := time.Parse(layout, start)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not format start interval, it must be in 20060102 format")
	}
	endInterval, err := time.Parse(layout, end)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not format end interval, it must be in 20060102 format")
	}
	// id, start and end obligatory
	if id == "" || start == "" || end == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing id, or date interval start or end")
	}
	// default USD
	if convert == "" {
		convert = "USD"
	}

	t := &schema.Ticker{
		ID:      id,
		Convert: convert,
		Start:   startInterval.Format(time.RFC3339),
		End:     endInterval.Format(time.RFC3339), // Start time of the interval in RFC3339 (URI escaped)
	}

	data, err := h.uc.RetrieveChart(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Add("Content-Type", "image/png")
	c.Response().Write(data.Image.Bytes())
	return c.JSON(http.StatusOK, data.Image.Bytes())
}
