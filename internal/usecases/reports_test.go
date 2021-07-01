package usecases

import (
	"errors"
	"gmgalvan/edChallenge2021/internal/schema"
	"gmgalvan/edChallenge2021/internal/usecases/mocks"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestNomics_getNomicsDataSparkline(t *testing.T) {
	type fields struct {
		m2m func(m *mocks.MockM2MClientCall)
	}
	type args struct {
		ticker *schema.Ticker
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*schema.CurrencySparkline
		wantErr bool
	}{
		{
			name: "Success get nomics data",
			fields: fields{
				m2m: func(m *mocks.MockM2MClientCall) {
					m.EXPECT().GET(gomock.Any()).
						Return([]byte(`[{"currency":"BTC", "timestamps":["2018-04-03T16:00:00Z"], "prices":["7436.82313"]}]`), nil)
				},
			},
			args: args{
				ticker: &schema.Ticker{
					ID:      "btc",
					Start:   "2018-01-12T00:00:00Z",
					End:     "2018-07-12T00:00:00Z",
					Convert: "USD",
				},
			},
			want: []*schema.CurrencySparkline{
				{
					Currency:   "BTC",
					Timestamps: []string{"2018-04-03T16:00:00Z"},
					Prices:     []string{"7436.82313"},
				},
			},
			wantErr: false,
		},

		{
			name: "Failed get nomics data",
			fields: fields{
				m2m: func(m *mocks.MockM2MClientCall) {
					m.EXPECT().GET(gomock.Any()).Return(nil, errors.New("status code different to 200 on calling nomics"))
				},
			},
			args: args{
				ticker: &schema.Ticker{
					ID:      "btc",
					Start:   "2018-01-12T00:00:00Z",
					End:     "2018-07-12T00:00:00Z",
					Convert: "USD",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			m2mClientCallMock := mocks.NewMockM2MClientCall(mockCtl)
			if tt.fields.m2m != nil {
				tt.fields.m2m(m2mClientCallMock)
			}
			n := &Nomics{
				m2m: m2mClientCallMock,
			}
			got, err := n.getNomicsDataSparkline(tt.args.ticker)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nomics.getNomicsDataSparkline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nomics.getNomicsDataSparkline() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNomics_RetrieveChart(t *testing.T) {
	type fields struct {
		m2m func(m *mocks.MockM2MClientCall)
	}
	type args struct {
		ticker *schema.Ticker
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schema.Chart
		wantErr bool
	}{
		{
			name: "Success retrieve chart",
			fields: fields{
				m2m: func(m *mocks.MockM2MClientCall) {
					m.EXPECT().GET(gomock.Any()).
						Return([]byte(`[{"currency":"BTC", "timestamps":["2018-04-03T16:00:00Z","2018-05-03T16:00:00Z"], "prices":["7436.82313", "7836.82313"]}]`), nil)
				},
			},
			args: args{
				ticker: &schema.Ticker{
					ID:      "btc",
					Start:   "2018-01-12T00:00:00Z",
					End:     "2018-07-12T00:00:00Z",
					Convert: "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "Failed retrieve chart",
			fields: fields{
				m2m: func(m *mocks.MockM2MClientCall) {
					m.EXPECT().GET(gomock.Any()).Return(nil, errors.New("status code different to 200 on calling nomics"))
				},
			},
			args: args{
				ticker: &schema.Ticker{
					ID:      "btc",
					Start:   "2018-01-12T00:00:00Z",
					End:     "2018-07-12T00:00:00Z",
					Convert: "USD",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			m2mClientCallMock := mocks.NewMockM2MClientCall(mockCtl)
			if tt.fields.m2m != nil {
				tt.fields.m2m(m2mClientCallMock)
			}
			n := &Nomics{
				m2m: m2mClientCallMock,
			}
			_, err := n.RetrieveChart(tt.args.ticker)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nomics.RetrieveChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func Test_stringToTime(t *testing.T) {
	parsedExample, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatalf("Could not start test, bad formating err: %v", err)
	}
	type args struct {
		list []string
	}
	tests := []struct {
		name    string
		args    args
		want    []time.Time
		wantErr bool
	}{
		{
			name: "Success convert string to time",
			args: args{
				list: []string{"2006-01-02T15:04:05Z"},
			},
			want:    []time.Time{parsedExample},
			wantErr: false,
		},
		{
			name: "Failed convert string to time due bad format string",
			args: args{
				list: []string{"2006-01-02"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToTime(tt.args.list)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringArrayToFloat64Array(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name    string
		args    args
		want    []float64
		wantErr bool
	}{
		{
			name: "Success convert string array to float64 array",
			args: args{
				list: []string{"1.1", "1.2", "1.3", "1.4"},
			},
			want:    []float64{1.1, 1.2, 1.3, 1.4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringArrayToFloat64Array(tt.args.list)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringArrayToFloat64Array() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringArrayToFloat64Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_max(t *testing.T) {
	type args struct {
		array []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Success get max number from array of floats",
			args: args{
				array: []float64{1.32, 0.43, 8.1, 3.3},
			},
			want: 8.1,
		},
		{
			name: "Success get max number from array of floats when number repeated",
			args: args{
				array: []float64{1.32, 0.43, 8.1, 8.1, 3.3},
			},
			want: 8.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.array); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		array []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Success get min number from array of floats",
			args: args{
				array: []float64{1.32, 0.43, 8.1, 3.3},
			},
			want: 0.43,
		},
		{
			name: "Success get min number from array of floats when number repeated",
			args: args{
				array: []float64{1.32, 0.43, 0.43, 8.1, 3.3},
			},
			want: 0.43,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.array); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}
