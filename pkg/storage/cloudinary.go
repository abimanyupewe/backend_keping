package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadReceipt(ctx context.Context, fileBytes []byte, filename string) (string, error)
}

type cloudinaryService struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewCloudinaryService(cloudName, apiKey, apiSecret, folder string) (CloudinaryService, error) {
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return &noopCloudinaryService{}, nil
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary client: %w", err)
	}

	if folder == "" {
		folder = "receipts"
	}

	return &cloudinaryService{
		cld:    cld,
		folder: folder,
	}, nil
}

func (s *cloudinaryService) UploadReceipt(ctx context.Context, fileBytes []byte, filename string) (string, error) {
	reader := bytes.NewReader(fileBytes)
	resp, err := s.cld.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:   s.folder,
		PublicID: filename,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload receipt to Cloudinary: %w", err)
	}

	return resp.SecureURL, nil
}

// Noop service when Cloudinary credentials are not configured
type noopCloudinaryService struct{}

func (s *noopCloudinaryService) UploadReceipt(ctx context.Context, fileBytes []byte, filename string) (string, error) {
	return "", nil
}
