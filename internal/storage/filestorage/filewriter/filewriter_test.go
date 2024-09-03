package filewriter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"testing"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		reader     io.ReadWriter
		permission fs.FileMode
	}
	tests := []struct {
		name    string
		args    args
		want    *Filewriter
		wantErr bool
	}{
		{
			name: "nil reader returns error",
			args: args{
				reader: nil,
			},
			wantErr: true,
		},
		{
			name: "reader ok returns new filewriter",
			args: args{
				reader: &bytes.Buffer{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.reader, tt.args.permission)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestFilewriter_SaveMetrics(t *testing.T) {
	type fields struct {
		file   io.ReadWriter
		writer *bufio.Writer
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid writer, returns ok",
			fields: fields{
				file:   &bytes.Buffer{},
				writer: bufio.NewWriter(&bytes.Buffer{}),
			},
		},
		{
			name:    "invalid writer, returns error",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			f := &Filewriter{
				file:   tt.fields.file,
				writer: tt.fields.writer,
			}
			err := f.SaveMetrics(ctx, tt.args.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestFilewriter_ReadMetrics(t *testing.T) {
	type args struct {
		metric *models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.Metrics
		wantErr bool
	}{
		{
			name: "valid metric in file, returns metric",
			args: args{
				metric: &models.Metrics{
					ID:    "test",
					MType: "gauge",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var file bytes.Buffer

			str, _ := json.Marshal(tt.args.metric)
			_, _ = file.Write(str)
			f, _ := New(&file, 0o666)
			got, err := f.ReadMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("Filewriter.ReadMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.NotEmpty(t, got)
			assert.Equal(t, tt.args.metric.ID, got[0].ID)
		})
	}
}

func TestFilewriter_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "close file reader",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filewriter{}
			err := f.Close()
			assert.Nil(t, err)
		})
	}
}

func TestFilewriter_FileIsNil(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "file is nill, returns true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filewriter{
				file: nil,
			}
			got := f.FileIsNil()
			assert.True(t, got)
		})
	}
}
