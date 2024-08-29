package filewriter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/VanGoghDev/practicum-metrics/internal/domain/models"
)

// Filewriter инкапсулирует логику работы с файлом os.
type Filewriter struct {
	file    io.ReadWriter
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

// New возвращает новый экземпляр FileWriter.
func New(file io.ReadWriter, permission fs.FileMode) (*Filewriter, error) {
	if file == nil {
		return nil, errors.New("reader is nil")
	}

	return &Filewriter{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (f *Filewriter) SaveMetrics(ctx context.Context, data []byte) error {
	if f.file == nil {
		return errors.New("file is nil")
	}
	_, err := f.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = f.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush data to file: %w", err)
	}

	return nil
}

func (f *Filewriter) ReadMetrics() ([]*models.Metrics, error) {
	metrics := make([]*models.Metrics, 0)
	for f.scanner.Scan() {
		metric := models.Metrics{}
		data := f.scanner.Bytes()
		if len(data) > 0 {
			err := json.Unmarshal(data, &metric)
			if err != nil {
				fmt.Println(f.scanner.Text())
				return nil, fmt.Errorf("failed to unmarshal metric: %w", err)
			}
			metrics = append(metrics, &metric)
		}
	}
	if err := f.scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return metrics, nil
}

func (f *Filewriter) FileIsNil() bool {
	return f.file == nil
}

// Close закрывает файл.
func (f *Filewriter) Close() error {
	return nil
}
