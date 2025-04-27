package supabase

import (
	"fmt"
	"itfest-2025/model"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
)

type Supabase struct {
	client storage_go.Client
}

type Interface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

func Init() Interface {
	url := fmt.Sprintf("%s/storage/v1", os.Getenv("SUPABASE_URL"))
	client := storage_go.NewClient(url, os.Getenv("SUPABASE_TOKEN"), nil)

	return Supabase{
		client: *client,
	}
}

func (s Supabase) UploadFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	path := uuid.NewString() + filepath.Ext(file.Filename)
	contentType, err := model.GetImageType(file)
	if err != nil {
		return "", err
	}

	_, err = s.client.UploadFile(
		os.Getenv("SUPABASE_BUCKET"),
		path,
		src,
		storage_go.FileOptions{
			ContentType: &contentType,
		},
	)

	if err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_BUCKET"),
		path,
	)

	return publicURL, nil
}
