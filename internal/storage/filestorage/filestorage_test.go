package filestorage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage"
	mock_filestorage "github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage/mocks"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	type args struct {
		restore    bool
		fileIsNil  bool
		memstorage *memstorage.MemStorage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "file is nil returns error",
			args: args{
				restore:    false,
				fileIsNil:  true,
				memstorage: &memstorage.MemStorage{},
			},
			wantErr: true,
		},
		{
			name: "memstorage is nil return error",
			args: args{
				restore:    true,
				fileIsNil:  false,
				memstorage: nil,
			},
			wantErr: true,
		},
		{
			name: "restore false returns file storage",
			args: args{
				restore:    false,
				fileIsNil:  false,
				memstorage: &memstorage.MemStorage{},
			},
			wantErr: false,
		},
		{
			name: "restore true",
			args: args{
				restore:    true,
				fileIsNil:  false,
				memstorage: &memstorage.MemStorage{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zlog, _ := zap.NewDevelopment()
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filewriterMock := mock_filestorage.NewMockFilewriter(ctrl)
			filewriterMock.EXPECT().FileIsNil().Return(tt.args.fileIsNil).AnyTimes()
			filewriterMock.EXPECT().ReadMetrics().Return([]*models.Metrics{}, nil).AnyTimes()

			got, err := filestorage.New(ctx, zlog, tt.args.restore, filewriterMock, tt.args.memstorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestFileStorage_SaveMetrics(t *testing.T) {
	type args struct {
		failedToSaveMetricToFile bool
		metrics                  []*models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no errors",
			args: args{
				failedToSaveMetricToFile: false,
				metrics: []*models.Metrics{
					{
						MType: "gauge",
						ID:    "test",
					},
					{
						MType: "counter",
						ID:    "test",
					}},
			},
			wantErr: false,
		},
		{
			name: "filewriter returns error, should return error",
			args: args{
				failedToSaveMetricToFile: true,
				metrics: []*models.Metrics{
					{
						MType: "gauge",
						ID:    "test",
					},
					{
						MType: "counter",
						ID:    "test",
					}},
			},

			wantErr: true,
		},
		{
			name: "unknown metric type",
			args: args{
				failedToSaveMetricToFile: false,
				metrics: []*models.Metrics{
					{
						MType: "cccc",
						ID:    "test",
					}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// preset
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filewriterMock := getFileWriterMock(t, tt.args.failedToSaveMetricToFile)
			zlog, _ := zap.NewDevelopment()
			mstrg, _ := memstorage.New(zlog)

			f := newTestFileStorage(t, filewriterMock, mstrg)

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

			// execute
			err := f.SaveMetrics(ctx, tt.args.metrics)

			// assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestFileStorage_SaveGauge(t *testing.T) {
	type args struct {
		name                     string
		value                    float64
		mstrgInit                bool
		failedToSaveMetricToFile bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "failed to save to file, returns error",
			args: args{
				name:                     "test1",
				value:                    1.1,
				mstrgInit:                true,
				failedToSaveMetricToFile: true,
			},
			wantErr: true,
		},
		{
			name: "memstorage gauges map is ok, returns no error",
			args: args{
				name:      "test1",
				value:     1.1,
				mstrgInit: true,
			},
			wantErr: false,
		},
		{
			name: "memstorage gauges map is nil, returns error",
			args: args{
				name:      "test1",
				value:     1.1,
				mstrgInit: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// preset
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filewriterMock := getFileWriterMock(t, tt.args.failedToSaveMetricToFile)
			mstrg := &memstorage.MemStorage{}

			if tt.args.mstrgInit {
				zlog, _ := zap.NewDevelopment()
				mstrg, _ = memstorage.New(zlog)
			}

			f := newTestFileStorage(t, filewriterMock, mstrg)

			// execute
			err := f.SaveGauge(ctx, tt.args.name, tt.args.value)

			// assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestFileStorage_SaveCount(t *testing.T) {
	type args struct {
		name                     string
		value                    int64
		mstrgInit                bool
		failedToSaveMetricToFile bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "failed to save to file, returns error",
			args: args{
				name:                     "test1",
				value:                    1,
				mstrgInit:                true,
				failedToSaveMetricToFile: true,
			},
			wantErr: true,
		},
		{
			name: "memstorage gauges map is ok, returns no error",
			args: args{
				name:      "test1",
				value:     1,
				mstrgInit: true,
			},
			wantErr: false,
		},
		{
			name: "memstorage gauges map is nil, returns error",
			args: args{
				name:      "test1",
				value:     1,
				mstrgInit: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// preset
			ctx := context.Background()

			mstrg := &memstorage.MemStorage{}

			if tt.args.mstrgInit {
				zlog, _ := zap.NewDevelopment()
				mstrg, _ = memstorage.New(zlog)
			}
			filewriterMock := getFileWriterMock(t, tt.args.failedToSaveMetricToFile)
			f := newTestFileStorage(t, filewriterMock, mstrg)

			// execute
			err := f.SaveCount(ctx, tt.args.name, tt.args.value)

			// assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestFileStorage_Ping(t *testing.T) {
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
			ctx := context.Background()
			filewriterMock := getFileWriterMock(t, false)
			mstrg := &memstorage.MemStorage{}
			f := newTestFileStorage(t, filewriterMock, mstrg)
			if err := f.Ping(ctx); (err != nil) != tt.wantErr {
				t.Errorf("FileStorage.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileStorage_Close(t *testing.T) {
	tests := []struct {
		name                   string
		filewriterReturnsError bool
		wantErr                bool
	}{
		{
			name:                   "close ok",
			filewriterReturnsError: false,
			wantErr:                false,
		},
		{
			name:                   "close returns error",
			filewriterReturnsError: true,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			filewriterMock := getFileWriterMock(t, false)
			if tt.filewriterReturnsError {
				filewriterMock.EXPECT().Close().Return(errors.New("test")).AnyTimes()
			} else {
				filewriterMock.EXPECT().Close().Return(nil).AnyTimes()
			}
			mstrg := &memstorage.MemStorage{}
			f := newTestFileStorage(t, filewriterMock, mstrg)
			err := f.Close(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStorage.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func getFileWriterMock(t *testing.T, failedToSaveToFile bool) *mock_filestorage.MockFilewriter {
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	filewriterMock := mock_filestorage.NewMockFilewriter(ctrl)
	filewriterMock.EXPECT().FileIsNil().Return(false).AnyTimes()
	filewriterMock.EXPECT().ReadMetrics().Return([]*models.Metrics{}, nil).AnyTimes()
	if failedToSaveToFile {
		filewriterMock.EXPECT().SaveMetrics(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
	} else {
		filewriterMock.EXPECT().SaveMetrics(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	}
	return filewriterMock
}

func newTestFileStorage(
	t *testing.T,
	fileWriterMock *mock_filestorage.MockFilewriter,
	mstrg *memstorage.MemStorage,
) *filestorage.FileStorage {
	t.Helper()
	zlog, _ := zap.NewDevelopment()
	ctx := context.Background()

	fstrg, _ := filestorage.New(ctx, zlog, false, fileWriterMock, mstrg)
	return fstrg
}
