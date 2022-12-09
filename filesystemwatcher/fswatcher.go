package filesystemwatcher

import (
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sjsafranek/logger"

	"github.com/sjsafranek/gosimpleserver/utils"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.New()
	//log.SetName("FileSystemWatcher")
	log.SetName("FSWatcher")
}

type FileSystemWatcher struct {
	lock        sync.RWMutex
	watcher     *fsnotify.Watcher
	directories map[string]bool
	OnCreate    func(filename string) error
	OnChange    func(filename string) error
	OnRename    func(filename string) error
	OnDelete    func(filename string) error
}

func (self *FileSystemWatcher) add(directory string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	log.Debugf("adding '%v' to watch list", directory)

	if nil == self.directories {
		self.directories = make(map[string]bool)
	}

	if _, ok := self.directories[directory]; !ok {
		self.directories[directory] = true
		return self.watcher.Add(directory)
	}

	return nil
}

// Add directory to watch list along with all subdirectories
func (self *FileSystemWatcher) Add(directory string) error {
	err := self.add(directory)
	if nil != err {
		return err
	}

	subdirectories, err := utils.GetSubDirectories(directory)
	if nil != err {
		return err
	}

	for _, subdirectory := range subdirectories {
		err = self.add(subdirectory)
		if nil != err {
			return err
		}
	}

	return nil
}

func (self *FileSystemWatcher) Remove(directory string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	log.Debugf("removing '%v' from watch list", directory)

	if nil == self.directories {
		return nil
	}

	if _, ok := self.directories[directory]; ok {
		delete(self.directories, directory)
		return self.watcher.Remove(directory)
	}

	return nil
}

func (self *FileSystemWatcher) IsWatching(directory string) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if nil == self.directories {
		return false
	}

	_, ok := self.directories[directory]
	return ok
}

func (self *FileSystemWatcher) Close() error {
	return self.watcher.Close()
}

func New() (*FileSystemWatcher, error) {

	// Create fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if nil != err {
		return nil, err
	}

	fileSystemWatcher := FileSystemWatcher{watcher: watcher}

	// Start go routine to handle file system events
	go func() {
		// Watch file system
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				//logger.Tracef("%+v", event)

				if fsnotify.Remove == event.Op {
					if fileSystemWatcher.IsWatching(event.Name) {
						fileSystemWatcher.Remove(event.Name)
						continue
					}

					if nil != fileSystemWatcher.OnDelete {
						fileSystemWatcher.OnDelete(event.Name)
					}
					continue
				}

				// Check to see if the event info
				info, err := os.Stat(event.Name)
				if nil != err {
					continue
				}

				// Handle directories
				if info.IsDir() {
					switch event.Op {
					case fsnotify.Create, fsnotify.Write:
						err = fileSystemWatcher.Add(event.Name)
						if nil != err {
							log.Errorf("%+v", err)
						}
					case fsnotify.Rename:
						log.Warnf("TODO")
					}
					continue
				}

				// Handle files
				switch event.Op {
				case fsnotify.Write:
					if nil != fileSystemWatcher.OnChange {
						fileSystemWatcher.OnChange(event.Name)
					}
				case fsnotify.Rename:
					if nil != fileSystemWatcher.OnRename {
						fileSystemWatcher.OnRename(event.Name)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("error: %v", err)
			}
		}
	}()

	return &fileSystemWatcher, nil
}
