package filesystem

import (
	"path"
	"strconv"
)

// SequencePath ...
func (fs *FileSystem) SequencePath(name string) string {
	return path.Join(fs.Root, "Meta", name+"Sequence.txt")
}

// ObjectPath ...
func (fs *FileSystem) ObjectPath(name string, id int) string {
	return path.Join(fs.Root, "Object", name, strconv.Itoa(id)+".yml")
}
