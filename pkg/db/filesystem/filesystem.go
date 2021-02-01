package filesystem

import (
	"os"
	"path"
	"sync"

	"github.com/tzvetkoff-go/errors"
	"gopkg.in/yaml.v2"

	"github.com/tzvetkoff-go/pasteur/pkg/config"
	"github.com/tzvetkoff-go/pasteur/pkg/model"
)

// FileSystem ...
type FileSystem struct {
	Root string
	Mux  sync.RWMutex

	PasteSequence int
}

// New ...
func New(fsConfig *config.FileSystem) (*FileSystem, error) {
	result := &FileSystem{
		Root: fsConfig.Root,
	}
	return result, nil
}

// CreatePaste ...
func (fs *FileSystem) CreatePaste(paste *model.Paste) (*model.Paste, error) {
	pasteID, err := fs.NextSequenceValue("Paste")
	if err != nil {
		return nil, errors.Propagate(err, "cannot get next sequence value")
	}
	paste.ID = pasteID

	data, err := yaml.Marshal(paste)
	if err != nil {
		return nil, errors.Propagate(err, "cannot marshal paste")
	}

	pastePath := fs.ObjectPath("Paste", pasteID)
	pasteDir := path.Dir(pastePath)
	err = os.MkdirAll(pasteDir, 0755)
	if err != nil {
		return nil, errors.Propagate(err, "cannot create paste directory")
	}

	err = os.WriteFile(pastePath, data, 0644)
	if err != nil {
		return nil, errors.Propagate(err, "cannot write paste")
	}

	return paste, nil
}

// RetrievePasteByID ...
func (fs *FileSystem) RetrievePasteByID(id int) (*model.Paste, error) {
	pastePath := fs.ObjectPath("Paste", id)
	_, err := os.Stat(pastePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}

		return nil, errors.Propagate(err, "cannot stat paste")
	}

	data, err := os.ReadFile(pastePath)
	if err != nil {
		return nil, errors.Propagate(err, "cannot read paste")
	}

	result := &model.Paste{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.Propagate(err, "cannot unmarshal paste")
	}

	return result, nil
}
