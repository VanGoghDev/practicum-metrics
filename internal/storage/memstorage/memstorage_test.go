package memstorage

import (
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"
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

type test struct {
	name   string
	fields fields
	args   args
	want   want
}

func TestGauge(t *testing.T) {
	tests := []test{
		{
			name: "get existing gauge",
			fields: fields{
				Gauges: map[string]float64{
					"test": 10,
				},
				Counters: map[string]int64{},
			},
			args: args{
				name: "test",
			},
			want: want{
				err:         nil,
				metricValue: 10,
			},
		},
		{
			name: "get non existing gauge",
			fields: fields{
				Gauges: map[string]float64{},
				Counters: map[string]int64{
					"test": 10,
				},
			},
			args: args{
				name: "test",
			},
			want: want{
				err: serrors.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zlog, _ := logger.New("Info")
			s, _ := New(zlog)
			s.GaugesM = tt.fields.Gauges
			s.CountersM = tt.fields.Counters
			gauge, err := s.Gauge(tt.args.name)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.metricValue, gauge.Value)
		})
	}
}

func TestCounter(t *testing.T) {
	type args struct {
		name string
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
			name: "get existing counter",
			fields: fields{
				Gauges: map[string]float64{},
				Counters: map[string]int64{
					"test": 10,
				},
			},
			args: args{
				name: "test",
			},
			want: want{
				err:         nil,
				metricValue: 10,
			},
		},
		{
			name: "get non existing gauge",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]int64{},
			},
			args: args{
				name: "test",
			},
			want: want{
				err: serrors.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, _ := logger.New("Info")

			s := &MemStorage{
				GaugesM:   tt.fields.Gauges,
				CountersM: tt.fields.Counters,
				zlog:      log,
			}
			counter, err := s.Counter(tt.args.name)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.metricValue, counter.Value)
		})
	}
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
				err: serrors.ErrCountersTableNil,
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
			log, _ := logger.New("Info")

			s := &MemStorage{
				GaugesM:   tt.fields.Gauges,
				CountersM: tt.fields.Counters,
				zlog:      log,
			}
			err := s.SaveCount(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.metricValue, s.CountersM[tt.args.name])
		})
	}
}

func TestSaveGauge(t *testing.T) {
	tests := []test{
		{
			name:   "empty gauge table",
			fields: fields{},
			args: args{
				name:  "test",
				value: 20,
			},
			want: want{
				err: serrors.ErrGaugesTableNil,
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
		runTest(t, &tt)
	}
}

func runTest(t *testing.T, tt *test) func(name string, f func(t *testing.T)) bool {
	t.Helper()
	return func(name string, f func(t *testing.T)) bool {
		log, _ := logger.New("Info")

		s := &MemStorage{
			GaugesM:   tt.fields.Gauges,
			CountersM: tt.fields.Counters,
			zlog:      log,
		}
		err := s.SaveGauge(tt.args.name, tt.args.value)
		assert.Equal(t, tt.want.err, err)
		return assert.Equal(t, tt.want.err, err) && assert.Equal(t, tt.want.metricValue, s.GaugesM[tt.args.name])
	}
}
