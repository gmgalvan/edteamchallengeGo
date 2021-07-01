package transport

import (
	"bytes"
	"errors"
	"gmgalvan/edChallenge2021/internal/schema"
	"gmgalvan/edChallenge2021/internal/transport/mocks"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
)

func TestHandler_getChart(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/chart?id=ETH&start=20210112&end=20211220&convert=MXN", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type fields struct {
		uc func(m *mocks.MockUsecases)
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success chart requested",
			fields: fields{
				uc: func(m *mocks.MockUsecases) {
					m.EXPECT().RetrieveChart(gomock.Any()).
						Return(&schema.Chart{
							Image: bytes.NewBuffer([]byte{}),
						}, nil)
				},
			},
			args: args{
				c: c,
			},
			wantErr: false,
		},
		{
			name: "Failed chart requested",
			fields: fields{
				uc: func(m *mocks.MockUsecases) {
					m.EXPECT().RetrieveChart(gomock.Any()).
						Return(nil, errors.New("Error on getting the chart"))
				},
			},
			args: args{
				c: c,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			ucMock := mocks.NewMockUsecases(mockCtl)
			if tt.fields.uc != nil {
				tt.fields.uc(ucMock)
			}
			h := &Handler{
				uc: ucMock,
			}
			if err := h.getChart(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Handler.getChart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
