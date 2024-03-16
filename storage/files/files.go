package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/vltvdnl/Adviser-Bot/lib/e"
	"github.com/vltvdnl/Adviser-Bot/storage"
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

const defaultPerm = 0774

var ErrNoSaved = errors.New("no saved pages")

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.Wrap("can't save", err) }()
	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}
	fName, err := FileName(page)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("can't pick random", err) }()

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoSaved
	}
	// rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.DecodePage(filepath.Join(path, file.Name()))
}
func (s Storage) Remove(p *storage.Page) error {
	// defer func() {err = e.Wrap("can't remove page", err)}()
	fileName, err := FileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("can't remove file %s", path), err)
	}
	return nil
}
func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := FileName(p)
	if err != nil {
		return false, e.Wrap("can't find file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)
	_, err = os.Stat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if %s exist", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) DecodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() {
		_ = f.Close()
	}()
	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decodge page", err)
	}
	return &p, nil
}

func FileName(p *storage.Page) (string, error) {
	return p.Hash()
}
