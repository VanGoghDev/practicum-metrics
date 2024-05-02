package memstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fields struct {
	Gauges   map[string]float64
	Counters map[string]int64
}
type args struct {
	name  string
	value float64
}
type want struct {
	err         error
	metricValue float64
}

func TestSaveCount(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	type want struct {
		err         error
		metricValue int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "empty counters table",
			fields: fields{},
			args: args{
				name:  "test",
				value: 20,
			},
			want: want{
				err: ErrCountersTableNil,
			},
		},
		{
			name: "save new counter",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
			},
			args: args{
				name:  "test",
				value: 20,
			},
			want: want{
				err:         nil,
				metricValue: 20,
			},
		},
		{
			name: "add to existing counter",
			fields: fields{
				Gauges: map[string]float64{},
				Counters: map[string]int64{
					"test": 10,
				},
			},
			args: args{
				name:  "test",
				value: 10,
			},
			want: want{
				err:         nil,
				metricValue: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				GaugesM:   tt.fields.Gauges,
				CountersM: tt.fields.Counters,
			}
			err := s.SaveCount(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.metricValue, s.CountersM[tt.args.name])
		})
	}
}

func TestSaveGauge(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "empty gauge table",
			fields: fields{},
			args: args{
				name:  "test",
				value: 20,
			},
			want: want{
				err: ErrGaugesTableNil,
			},
		},
		{
			name: "update gauge",
			fields: fields{
				Gauges: map[string]float64{
					"test": 1,
				},
				Counters: map[string]int64{},
			},
			args: args{
				name:  "test",
				value: 20,
			},
			want: want{
				err:         nil,
				metricValue: 20,
			},
		},
		{
			name: "save new gauge",
			fields: fields{
				Gauges: map[string]float64{
					"test": 1,
				},
				Counters: map[string]int64{},
			},
			args: args{
				name:  "test2",
				value: 1,
			},
			want: want{
				err:         nil,
				metricValue: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				GaugesM:   tt.fields.Gauges,
				CountersM: tt.fields.Counters,
			}
			err := s.SaveGauge(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.metricValue, s.GaugesM[tt.args.name])
		})
	}
}
