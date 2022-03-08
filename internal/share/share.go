package share

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type Share struct {
	UUID           string
	FullPath       string
	ExpirationDate time.Time
	MaxAccessCount int
}

type FS map[string]Share

func (fs FS) GetPath(filename string) string {
	return fs[filename].FullPath
}

func (fs FS) GetUUID(filename string) string {
	return fs[filename].UUID
}

func (fs FS) GetExpirationDate(filename string) time.Time {
	return fs[filename].ExpirationDate
}

func (fs FS) GetMaxAccessCount(filename string) int {
	return fs[filename].MaxAccessCount
}

func (fs FS) IsShareExist(uuid string) bool {
	if err := isValidUUID(uuid); err != nil {
		return false
	}

	for _, v := range fs {
		if v.UUID == uuid {
			return true
		}
	}
	return false
}

func (fs FS) GetShareFromUUID(uuid string) (string, error) {
	if err := isValidUUID(uuid); err != nil {
		return "", err
	}
	for k, v := range fs {
		if v.UUID == uuid {
			return k, nil
		}
	}
	return "", errors.New("share not found")
}

func (fs FS) GetPathFromUUID(link string) (string, error) {
	if err := isValidUUID(link); err != nil {
		return "", err
	}

	for _, v := range fs {
		if v.UUID == link {
			return v.FullPath, nil
		}
	}
	return "", errors.New("share not found")
}

func isValidUUID(link string) error {
	_, err := uuid.Parse(link)
	if err != nil {
		return err
	}
	return nil
}

func (fs FS) DeleteShare(filename string) {
	delete(fs, filename)
}

func (fs FS) ProcessShare(filename string) (string, error) {
	if v, ok := fs[filename]; ok {
		if v.MaxAccessCount == 0 {
			fs.DeleteShare(filename)
			return "", errors.New("limit of access reached")
		}

		v.MaxAccessCount--
		fs[filename] = v
		return v.FullPath, nil
	}
	return "", errors.New("UUID not found")
}

func (fs FS) SaveFile(file *multipart.FileHeader, filename, dst string, count int) error {
	err := createFile(file, dst)
	if err != nil {
		return err
	}

	fs[filename] = Share{
		UUID:           uuid.New().String(),
		FullPath:       dst,
		MaxAccessCount: count,
	}
	return nil
}

func createFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)

	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
