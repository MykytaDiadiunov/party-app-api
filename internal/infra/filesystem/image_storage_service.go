package filesystem

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

type ImageStorageService interface {
	GetImages() ([]string, error)
	FileIsExist(filename string) (bool, error)
	SaveImage(filename string, content []byte) error
	RemoveImage(filename string) error
}

type imageStorageService struct {
	loc string
}

func NewImageStorageService(location string) ImageStorageService {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	absLocation := path.Join(rootDir, location)

	err = os.MkdirAll(absLocation, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	return imageStorageService{
		loc: absLocation,
	}
}

func (s imageStorageService) GetImages() ([]string, error) {
	var images []string

	dir, err := os.Open(s.loc)
	if err != nil {
		log.Printf("Error opening directory: %v", err)
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Error reading directory contents: %v", err)
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			images = append(images, filepath.Join(s.loc, file.Name()))
		}
	}

	return images, nil
}

func (s imageStorageService) FileIsExist(filename string) (bool, error) {
	images, err := s.GetImages()
	if err != nil {
		return false, err
	}

	for _, image := range images {
		if image == "file_storage/"+filename {
			return true, nil
		}
	}
	return false, nil
}

func (s imageStorageService) SaveImage(filename string, content []byte) error {
	location := path.Join(s.loc, filename)
	err := writeFileToStorage(location, content)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s imageStorageService) RemoveImage(filename string) error {
	location := path.Join(s.loc, filename)
	err := os.Remove(location)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func writeFileToStorage(location string, file []byte) error {
	dirLocation := path.Dir(location)
	err := os.MkdirAll(dirLocation, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}

	err = os.WriteFile(location, file, os.ModePerm)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
