package filesystem

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gofrs/flock"
	"github.com/tzvetkoff-go/logger"
)

// LoadSequence ...
func (fs *FileSystem) LoadSequence(name string) int {
	sequencePath := fs.SequencePath(name)
	_, err := os.Stat(sequencePath)
	if err != nil {
		logger.Warning("%s", err)
		return 0
	}

	b, err := os.ReadFile(sequencePath)
	if err != nil {
		logger.Warning("%s", err)
		return 0
	}

	result, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		logger.Warning("%s", err)
		return 0
	}

	return result
}

// StoreSequence ...
func (fs *FileSystem) StoreSequence(name string, value int) error {
	sequencePath := fs.SequencePath(name)

	sequenceDir := path.Dir(sequencePath)
	err := os.MkdirAll(sequenceDir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(sequencePath, []byte(strconv.Itoa(value)+"\n"), 0644)
	if err != nil {
		return err
	}

	return nil
}

// NextSequenceValue ...
func (fs *FileSystem) NextSequenceValue(name string) (int, error) {
	lockPath := fs.SequencePath(name) + ".lock"

	lockDir := path.Dir(lockPath)
	err := os.MkdirAll(lockDir, 0755)
	if err != nil {
		logger.Error("%s", err)
		return 0, err
	}

	os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL, 0644)

	sequenceLock := flock.New(fs.SequencePath(name) + ".lock")
	err = sequenceLock.Lock()
	if err != nil {
		logger.Error("%s", err)
		return 0, err
	}
	defer sequenceLock.Unlock()

	value := fs.LoadSequence(name) + 1
	err = fs.StoreSequence(name, value)
	return value, err
}
