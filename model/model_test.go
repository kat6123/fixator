package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixationFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		f          FixationFloat
		want       string
		wantErrMsg string
	}{
		{
			name: "marshal 65.5",
			f:    65.5,
			want: "\"65,50\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.MarshalJSON()
			if err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
			}

			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestFixationFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		b          []byte
		f          FixationFloat
		wantErrMsg string
	}{
		{
			name: "unmarshal \"65,5\"",
			b:    []byte("65,5"),
			f:    65.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotF FixationFloat
			if err := gotF.UnmarshalJSON(tt.b); err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
			}

			assert.Equal(t, tt.f, gotF)
		})
	}
}

func TestFixationTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		time       string
		want       []byte
		wantErrMsg string
	}{
		{
			name: "marshal",
			time: "20.12.2019 14:31:25",
			want: []byte("\"20.12.2019 14:31:25\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTime, err := time.Parse(layout, tt.time)
			if err != nil {
				fmt.Printf("parse fixation time: %v", err)
				return
			}
			got, err := FixationTime(parsedTime).MarshalJSON()
			if err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFixationTime_UnmarshalJSON(t *testing.T) {
	const layout = "02.01.2006 15:04:05"
	tests := []struct {
		name       string
		b          []byte
		parsedTime string
		wantErrMsg string
	}{
		{
			name:       "unmarshal",
			b:          []byte("\"20.12.2019 14:31:25\""),
			parsedTime: "20.12.2019 14:31:25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotT FixationTime
			if err := gotT.UnmarshalJSON(tt.b); err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
			}

			assert.Equal(t, tt.parsedTime, time.Time(gotT).Format(layout))
		})
	}
}
