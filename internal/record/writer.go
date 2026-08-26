package record

import "os"

type FileWriter struct {
	file *os.File
	path string
}

func (w *FileWriter) Open(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.path = path
	return nil
}

func (w *FileWriter) Write(data []byte) error {
	if w.file == nil {
		return nil
	}
	_, err := w.file.Write(data)
	return err
}

func (w *FileWriter) Close() error {
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	w.path = ""
	return err
}
