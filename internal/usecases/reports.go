package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmgalvan/edChallenge2021/internal/schema"
	"os"
	"strconv"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

var NomicsAPIKey = os.Getenv("NOMICS_KEY")

type M2MClientCall interface {
	GET(url string) ([]byte, error)
}

type Nomics struct {
	m2m M2MClientCall
}

func NewNomicsReport(m2m M2MClientCall) *Nomics {
	return &Nomics{
		m2m: m2m,
	}
}

func (n *Nomics) getNomicsDataSparkline(ticker *schema.Ticker) ([]*schema.CurrencySparkline, error) {
	result, err := n.m2m.GET(`https://api.nomics.com/v1/currencies/sparkline?key=` + NomicsAPIKey + `&ids=` + ticker.ID + `&start=` + ticker.Start + `&end=` + ticker.End + `&convert=` + ticker.Convert)
	if err != nil {
		return nil, err
	}
	var sparkLine []*schema.CurrencySparkline
	err = json.Unmarshal(result, &sparkLine)
	if err != nil {
		return nil, err
	}
	return sparkLine, nil
}

func (n *Nomics) RetrieveChart(ticker *schema.Ticker) (*schema.Chart, error) {
	data, err := n.getNomicsDataSparkline(ticker)
	if err != nil {
		return nil, err
	}

	pricesFloat, err := stringArrayToFloat64Array(data[0].Prices)
	if err != nil {
		return nil, err
	}
	timestamps, err := stringToTime(data[0].Timestamps)
	if err != nil {
		return nil, err
	}

	priceSeries := chart.TimeSeries{
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: timestamps,
		YValues: pricesFloat,
	}

	graph := chart.Chart{
		Title: fmt.Sprintf("%v/%v", ticker.ID, ticker.Convert),
		TitleStyle: chart.Style{
			Show: true,
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Range: &chart.ContinuousRange{
				Max: max(pricesFloat) + (max(pricesFloat) * .15),
				Min: min(pricesFloat),
			},
		},
		Series: []chart.Series{
			priceSeries,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return &schema.Chart{
		Image: buffer,
	}, nil
}

func stringArrayToFloat64Array(list []string) ([]float64, error) {
	var listFloat64 []float64
	for _, e := range list {
		var floatElement float64
		floatElement, err := strconv.ParseFloat(e, 64)
		if err != nil {
			return nil, err
		}
		listFloat64 = append(listFloat64, floatElement)
	}
	return listFloat64, nil
}

func stringToTime(list []string) ([]time.Time, error) {
	var dates []time.Time
	for _, ts := range list {
		parsed, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return nil, err
		}
		dates = append(dates, parsed)
	}
	return dates, nil
}

func min(array []float64) float64 {
	var min float64 = array[0]
	for _, value := range array {
		if min > value {
			min = value
		}
	}
	return min
}
func max(array []float64) float64 {
	var max float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}
