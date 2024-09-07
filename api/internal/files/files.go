package files

import (
	"os"
	"time"

	"math/rand"
)

type FileManager struct {
}

func New() *FileManager {
	return &FileManager{}
}

func generateRandomString() string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}
	return string(b)
}

func (f *FileManager) UploadImage(file []byte, name string) (string, error) {
	filePath := os.Getenv("FILE_PATH") + "/images/" + name + generateRandomString() + ".png"
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer newFile.Close()

	_, err = newFile.Write(file)
	if err != nil {
		return "", err
	}

	return filePath, nil

}
