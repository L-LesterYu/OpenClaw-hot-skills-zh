package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
)

// Compressor 图片压缩器
type Compressor struct {
	log          *zap.Logger
	maxWidth     int
	maxSize      int64
	quality      int // JPEG 质量 1-100
	enableResize bool
	enableShrink bool
}

// NewCompressor 创建压缩器
func NewCompressor(log *zap.Logger, maxWidth int, maxSize int64) *Compressor {
	return &Compressor{
		log:          log,
		maxWidth:     maxWidth,
		maxSize:      maxSize,
		quality:      85, // 默认 JPEG 质量
		enableResize: maxWidth > 0,
		enableShrink: maxSize > 0,
	}
}

// CompressImage 压缩图片
// 返回: 压缩后的文件路径, 是否进行了压缩, 错误
func (c *Compressor) CompressImage(filePath string) (string, bool, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", false, fmt.Errorf("stat file: %w", err)
	}

	// 检查文件大小
	if c.enableShrink && fileInfo.Size() <= c.maxSize {
		c.log.Debug("file size within limit, no compression needed",
			zap.Int64("size", fileInfo.Size()),
			zap.Int64("max", c.maxSize))
		return "", false, nil
	}

	// 打开图片文件
	img, err := imaging.Open(filePath)
	if err != nil {
		return "", false, fmt.Errorf("open image: %w", err)
	}

	// 获取原始尺寸
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	c.log.Debug("image loaded",
		zap.String("path", filePath),
		zap.Int("width", originalWidth),
		zap.Int("height", originalHeight),
		zap.Int64("size", fileInfo.Size()))

	// 判断是否需要调整尺寸
	var processedImg image.Image
	needsResize := c.enableResize && originalWidth > c.maxWidth

	if needsResize {
		// 计算新的高度，保持宽高比
		newHeight := int(float64(c.maxWidth) * float64(originalHeight) / float64(originalWidth))
		processedImg = imaging.Resize(img, c.maxWidth, newHeight, imaging.Lanczos)

		c.log.Info("image resized",
			zap.Int("original_width", originalWidth),
			zap.Int("original_height", originalHeight),
			zap.Int("new_width", c.maxWidth),
			zap.Int("new_height", newHeight))
	} else {
		processedImg = img
	}

	// 创建临时文件保存压缩结果
	ext := filepath.Ext(filePath)

	// 根据扩展名确定输出格式
	var outputFormat string
	switch strings.ToLower(ext) {
	case ".png":
		outputFormat = "png"
	case ".jpg", ".jpeg":
		outputFormat = "jpeg"
	default:
		// 默认使用 JPEG 格式以获得更好的压缩
		outputFormat = "jpeg"
		ext = ".jpg"
	}

	tmpFile, err := os.CreateTemp("", "md2wechat-compressed-*"+ext)
	if err != nil {
		return "", false, fmt.Errorf("create temp file: %w", err)
	}
	tempPath := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return "", false, fmt.Errorf("close temp file: %w", err)
	}

	// 保存压缩后的图片
	if err := c.saveImage(processedImg, tempPath, outputFormat); err != nil {
		return "", false, fmt.Errorf("save compressed image: %w", err)
	}

	// 检查压缩后的大小
	newFileInfo, err := os.Stat(tempPath)
	if err != nil {
		_ = os.Remove(tempPath)
		return "", false, fmt.Errorf("stat compressed file: %w", err)
	}

	compressionRatio := float64(newFileInfo.Size()) / float64(fileInfo.Size()) * 100

	c.log.Info("image compressed",
		zap.Int64("original_size", fileInfo.Size()),
		zap.Int64("compressed_size", newFileInfo.Size()),
		zap.Float64("ratio", compressionRatio),
		zap.String("output_path", tempPath))

	// 如果压缩后反而变大，删除临时文件并返回原路径
	if newFileInfo.Size() >= fileInfo.Size() {
		_ = os.Remove(tempPath)
		c.log.Debug("compressed image larger than original, using original")
		return "", false, nil
	}

	return tempPath, true, nil
}

// saveImage 保存图片到文件
func (c *Compressor) saveImage(img image.Image, filePath, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: c.quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(file, img)
	default:
		// 默认使用 JPEG
		return jpeg.Encode(file, img, &jpeg.Options{Quality: c.quality})
	}
}

// SetQuality 设置 JPEG 压缩质量 (1-100)
func (c *Compressor) SetQuality(quality int) {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	c.quality = quality
}

// GetImageDimensions 获取图片尺寸
func GetImageDimensions(filePath string) (width, height int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("decode config: %w", err)
	}

	return img.Width, img.Height, nil
}

// GetImageFormat 获取图片格式
func GetImageFormat(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", fmt.Errorf("decode config: %w", err)
	}

	return format, nil
}

// IsValidImageFormat 检查是否是有效的图片格式
func IsValidImageFormat(filePath string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	return validExts[ext]
}

// ImageInfo 图片信息
type ImageInfo struct {
	Width    int    // 宽度
	Height   int    // 高度
	Format   string // 格式
	Size     int64  // 文件大小
	FilePath string // 文件路径
}

// GetImageInfo 获取图片信息
func GetImageInfo(filePath string) (*ImageInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	width, height, err := GetImageDimensions(filePath)
	if err != nil {
		return nil, err
	}

	format, err := GetImageFormat(filePath)
	if err != nil {
		return nil, err
	}

	return &ImageInfo{
		Width:    width,
		Height:   height,
		Format:   format,
		Size:     fileInfo.Size(),
		FilePath: filePath,
	}, nil
}

// NeedsCompression 检查图片是否需要压缩
func NeedsCompression(info *ImageInfo, maxWidth int, maxSize int64) bool {
	if maxWidth > 0 && info.Width > maxWidth {
		return true
	}
	if maxSize > 0 && info.Size > maxSize {
		return true
	}
	return false
}
