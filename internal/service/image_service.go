package service

import (
	"context"
	"io"

	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

// ImageUploader is what EventService depends on — not ImageKit directly.
// EventService only knows "give me bytes + a filename, I get back a URL
// and an ID to delete it by later." Everything ImageKit-specific
// (the client, the context, the SDK's own param struct) stays inside
// imageKitUploader, on the other side of this interface.
type ImageUploader interface {
	// Upload reads file to the end and stores it, returning the public
	// URL and the provider's internal file ID (needed for Delete).
	Upload(file io.Reader, filename string) (url string, fileID string, err error)
	Delete(fileID string) error
}

type imageKitUploader struct {
	client *imagekit.Client
}

// NewImageKitUploader builds the real, ImageKit-backed implementation.
// This is what main.go wires up. A test wiring EventService instead
// injects a fake ImageUploader — no real network call, no real
// ImageKit account needed to test Create/Update's business logic.
func NewImageKitUploader(privateKey string) ImageUploader {
	client := imagekit.NewClient(option.WithPrivateKey(privateKey))
	return &imageKitUploader{client: &client}
}

func (u *imageKitUploader) Upload(file io.Reader, filename string) (string, string, error) {
	res, err := u.client.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: filename,
	})
	if err != nil {
		return "", "", err
	}
	return res.URL, res.FileID, nil
}

func (u *imageKitUploader) Delete(fileID string) error {
	return u.client.Files.Delete(context.Background(), fileID)
}
