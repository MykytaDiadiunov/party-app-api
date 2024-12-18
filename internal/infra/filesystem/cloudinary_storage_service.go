package filesystem

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-rest-api/config"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	_ "github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloudinaryObject *cloudinary.Cloudinary
}

func NewCloudinaryService(config config.Configuration) *CloudinaryService {
	cloudinaryObject, err := cloudinary.NewFromParams(config.CloudinaryNameKey, config.CloudinaryApiKey, config.CloudinarySecretKey)
	if err != nil {
		log.Printf("cloudinaryService: failed to initialize cloudinary: %s", err)
	}
	return &CloudinaryService{cloudinaryObject: cloudinaryObject}
}

func (c *CloudinaryService) SaveImageToCloudinary(imageBase64 string, fileName string) (string, error) {
	ctx := context.Background()
	decodedImage, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return "", fmt.Errorf("cloudinaryService: failed to decode image: %w", err)
	}

	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", fmt.Errorf("cloudinaryService: failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(decodedImage)
	if err != nil {
		return "", fmt.Errorf("cloudinaryService: failed to write to temporary file: %w", err)
	}

	uploadResult, err := c.uploadImage(ctx, tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("cloudinaryService: failed to upload image to Cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

func (c *CloudinaryService) DeleteImage(imageUrl string) error {
	ctx := context.Background()
	publicId := c.imageUrlToImagePublicId(imageUrl)
	_, err := c.cloudinaryObject.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})
	if err != nil {
		log.Printf("cloudinaryServiece.DeleteImage: %s", err)
		return err
	}

	return nil
}

func (c *CloudinaryService) uploadImage(ctx context.Context, filePath string) (*uploader.UploadResult, error) {
	result, err := c.cloudinaryObject.Upload.Upload(ctx, filePath, uploader.UploadParams{})
	if err != nil {
		log.Printf("cloudinaryServiece.uploadImage: %s", err)
		return nil, err
	}
	return result, nil
}

func (c *CloudinaryService) imageUrlToImagePublicId(imageUrl string) string {
	filePath := path.Base(imageUrl)
	imagePublicId := strings.TrimSuffix(filePath, path.Ext(filePath))
	return imagePublicId
}
