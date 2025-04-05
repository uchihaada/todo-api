package storage

import (
	"ada/models"
	"encoding/json"
	"os"
	"sync"
)

type FileStorage struct {
	FilePath  string
	taskChan  chan []models.Task
	closechan chan struct{}
	wg        sync.WaitGroup
}

func NewFileStorage(filepath string) *FileStorage {
	fs := &FileStorage{
		FilePath:  filepath,
		taskChan:  make(chan []models.Task, 10),
		closechan: make(chan struct{}),
	}
	fs.wg.Add(1)
	go fs.listenForUpdates()
	return fs
}

func (fs *FileStorage) listenForUpdates() {
	defer fs.wg.Done()
	for {
		select {
		case task := <-fs.taskChan:
			file, err := os.Create(fs.FilePath)
			if err != nil {
				continue
			}
			encoder := json.NewEncoder(file)
			err = encoder.Encode(task)
			file.Close()
			if err != nil {
				continue
			}
		case <-fs.closechan:
			return
		}
	}
}

func (fs *FileStorage) LoadTasks() ([]models.Task, error) {
	file, err := os.Open(fs.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Task{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var tasks []models.Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return []models.Task{}, nil
	}
	return tasks, nil
}

func (fs *FileStorage) SaveTasks(tasks []models.Task) {
	select {
	case fs.taskChan <- tasks:
	default:
	}
}

func (fs *FileStorage) Close() {
	close(fs.taskChan)
	close(fs.closechan)
	fs.wg.Wait()
}
