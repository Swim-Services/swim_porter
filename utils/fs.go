package utils

import (
	"fmt"
	"strings"
	"sync"
)

type MapFS struct {
	mut sync.RWMutex
	fs  map[string][]byte
}

func NewMapFS(in map[string][]byte) *MapFS {
	return &MapFS{fs: in}
}

func (m *MapFS) RawMap() map[string][]byte {
	return m.fs
}

func (m *MapFS) Write(data []byte, name string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.fs[strings.TrimPrefix(name, "/")] = data
}

func (m *MapFS) Read(name string) ([]byte, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	if data, ok := m.fs[name]; !ok {
		if data, ok := m.fs["/"+name]; !ok {
			return []byte{}, fmt.Errorf("file not found: %s", name)
		} else {
			return data, nil
		}
	} else {
		return data, nil
	}
}

func (m *MapFS) Dir(dirname string) map[string][]byte {
	m.mut.RLock()
	defer m.mut.RUnlock()
	out := make(map[string][]byte)
	for name, data := range m.fs {
		if strings.HasPrefix(name, dirname) {
			out[name[len(dirname):]] = data
		}
		if strings.HasPrefix("/"+name, dirname) {
			out[name[len(dirname)+1:]] = data
		}
	}
	return out
}

func (m *MapFS) DirExists(dirname string) bool {
	m.mut.RLock()
	defer m.mut.RUnlock()
	for name := range m.fs {
		if strings.HasPrefix(name, dirname) {
			return true
		}
	}
	return false
}

func (m *MapFS) Delete(name string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	delete(m.fs, name)
}

func (m *MapFS) Rename(oldName string, newName string) error {
	data, err := m.Read(oldName)
	if err != nil {
		return err
	}
	m.mut.Lock()
	delete(m.fs, oldName)
	m.mut.Unlock()
	m.Write(data, newName)
	return nil
}

func (m *MapFS) Copy(src string, dest string) error {
	data, err := m.Read(src)
	if err != nil {
		return err
	}
	m.Write(data, dest)
	return nil
}

func (m *MapFS) InterCopy(srcFS *MapFS, src string, dest string) error {
	data, err := srcFS.Read(src)
	if err != nil {
		return err
	}
	m.Write(data, dest)
	return nil
}

func (m *MapFS) CopyDir(src string, dest string) {
	for name, data := range m.Dir(src) {
		m.Write(data, dest+name)
	}
}

func (m *MapFS) InterCopyDir(srcFS *MapFS, src string, dest string) {
	for name, data := range srcFS.Dir(src) {
		m.Write(data, dest+name)
	}
}
