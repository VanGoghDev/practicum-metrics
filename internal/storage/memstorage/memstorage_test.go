package memstorage

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/serrors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
			gauge, err := s.Gauge(context.Background(), tt.args.name)
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
			counter, err := s.Counter(context.Background(), tt.args.name)
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
			err := s.SaveCount(context.Background(), tt.args.name, tt.args.value)
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
		err := s.SaveGauge(context.Background(), tt.args.name, tt.args.value)
		assert.Equal(t, tt.want.err, err)
		return assert.Equal(t, tt.want.err, err) && assert.Equal(t, tt.want.metricValue, s.GaugesM[tt.args.name])
	}
}

func BenchmarkSaveGauge(b *testing.B) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	type testMetric struct {
		name  string
		value float64
	}

	log, _ := logger.New("Info")
	s, _ := New(log)
	metrics := make([]*testMetric, 0, 100)

	for i := range 100 {
		_ = i
		metric := &testMetric{
			name:  randSeq(10),
			value: 10,
		}
		metrics = append(metrics, metric)
	}
	b.ResetTimer()
	for i := range metrics {
		err := s.SaveGauge(ctx, metrics[i].name, metrics[i].value)
		if err != nil {
			log.Sugar().Errorf("%w: failed to save metric", err)
		}
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func TestMemStorage_Gauges(t *testing.T) {
	type fields struct {
		zlog      *zap.Logger
		GaugesM   map[string]float64
		CountersM map[string]int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantGauges []models.Gauge
		wantErr    bool
	}{
		{
			name: "gauegesM is nil, should return error",
			fields: fields{
				GaugesM: nil,
			},
			wantErr: true,
		},
		{
			name: "gauegesM is ok, should return slice of metrics",
			fields: fields{
				GaugesM: map[string]float64{
					"test": 3.14,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := &MemStorage{
				zlog:      tt.fields.zlog,
				GaugesM:   tt.fields.GaugesM,
				CountersM: tt.fields.CountersM,
			}
			gotGauges, err := s.Gauges(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Gauges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, gotGauges)
		})
	}
}

func TestMemStorage_Counter(t *testing.T) {
	type fields struct {
		zlog      *zap.Logger
		GaugesM   map[string]float64
		CountersM map[string]int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantGauges []models.Gauge
		wantErr    bool
	}{
		{
			name: "gauegesM is nil, should return error",
			fields: fields{
				CountersM: nil,
			},
			wantErr: true,
		},
		{
			name: "gauegesM is ok, should return slice of metrics",
			fields: fields{
				CountersM: map[string]int64{
					"test": 3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := &MemStorage{
				zlog:      tt.fields.zlog,
				GaugesM:   tt.fields.GaugesM,
				CountersM: tt.fields.CountersM,
			}
			gotGauges, err := s.Counters(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Gauges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, gotGauges)
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	type fields struct {
		zlog      *zap.Logger
		GaugesM   map[string]float64
		CountersM map[string]int64
	}
	type args struct {
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*models.Metrics
		wantCount int
		wantErr   bool
	}{
		{
			name: "gauegesM is ok, should return slice of metrics",
			fields: fields{
				CountersM: map[string]int64{
					"test": 3,
				},
				GaugesM: map[string]float64{
					"test2": 4,
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := &MemStorage{
				zlog:      tt.fields.zlog,
				GaugesM:   tt.fields.GaugesM,
				CountersM: tt.fields.CountersM,
			}
			got, err := s.GetMetrics(ctx)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Equal(t, len(got), tt.wantCount)
		})
	}
}

func TestMemStorage_SaveMetrics(t *testing.T) {
	type fields struct {
		zlog      *zap.Logger
		GaugesM   map[string]float64
		CountersM map[string]int64
	}
	type args struct {
		metrics []*models.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "mem storage ok, returns result",
			fields: fields{
				CountersM: map[string]int64{
					"test": 3,
				},
				GaugesM: map[string]float64{
					"test2": 4,
				},
			},
			args: args{
				metrics: []*models.Metrics{
					{
						MType: "gauge",
						ID:    "test",
					},
					{
						MType: "counter",
						ID:    "test",
					},
				},
			},
		},
		{
			name: "CountersM nil, returns error",
			fields: fields{
				CountersM: nil,
				GaugesM: map[string]float64{
					"test2": 4,
				},
			},
			args: args{
				metrics: []*models.Metrics{
					{
						MType: "gauge",
						ID:    "test",
					},
					{
						MType: "counter",
						ID:    "test",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "GaugesM nil, returns error",
			fields: fields{
				CountersM: map[string]int64{
					"test": 3,
				},
				GaugesM: nil,
			},
			args: args{
				metrics: []*models.Metrics{
					{
						MType: "gauge",
						ID:    "test",
					},
					{
						MType: "counter",
						ID:    "test",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				zlog:      tt.fields.zlog,
				GaugesM:   tt.fields.GaugesM,
				CountersM: tt.fields.CountersM,
			}
			for _, v := range tt.args.metrics {
				switch v.MType {
				case handlers.Counter:
					pollCount := int64(3)
					v.Delta = &pollCount
				case handlers.Gauge:
					val := float64(13.4)
					v.Value = &val
				}
			}
			ctx := context.Background()
			if err := s.SaveMetrics(ctx, tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.SaveMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Ping(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "ping",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{}
			ctx := context.Background()
			if err := s.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Close(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "close"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{}
			ctx := context.Background()
			if err := s.Close(ctx); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
