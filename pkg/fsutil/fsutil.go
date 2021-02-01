package fsutil

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/tzvetkoff-go/errors"
)

// Extract ...
func Extract(fsys fs.FS, rootPath string) error {
	return Walk(fsys, func(filePath string) error {
		destPath := path.Join(rootPath, filePath)

		err := os.MkdirAll(path.Dir(destPath), 0755)
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(fsys, filePath[1:])
		if err != nil {
			return err
		}

		err = os.WriteFile(destPath, data, 0644)
		if err != nil {
			return err
		}

		fmt.Println(filePath, "=>", destPath)
		return nil
	})
}

// Walk ...
func Walk(fsys fs.FS, callback func(string) error) error {
	return walkInternal(fsys, callback, "/")
}

func walkInternal(fsys fs.FS, callback func(string) error, currentPath string) error {
	fsysReadDir, ok := fsys.(fs.ReadDirFS)
	if !ok {
		return errors.New("%s: filesystem does not support readdir", currentPath)
	}

	entries, err := fsysReadDir.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			sub, err := fs.Sub(fsys, entry.Name())
			if err != nil {
				continue
			}

			err = walkInternal(sub, callback, path.Join(currentPath, entry.Name()))
			if err != nil {
				return err
			}
		} else {
			err = callback(path.Join(currentPath, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
