package main

import (
	"flag"
	"fmt"
	"github.com/chai2010/webp"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// go run iamge.go -url="https://picx.zhimg.com/80/v2-827365267ea950b8039af0f1259f25ac_1440w.webp?source=1940ef5c"
func main() {

	imageURL := flag.String("url", "https://picx.zhimg.com/80/v2-827365267ea950b8039af0f1259f25ac_1440w.webp?source=1940ef5c", "URL or file path of the image to convert")
	// 默认jpeg
	outputFormat := flag.String("format", "jpeg", "Output image format (jpg or jpeg)")

	flag.Parse()

	if *imageURL == "" {
		fmt.Println("Please provide the URL or file path of the image to convert")
		return
	}
	imageData, err := getImageData(*imageURL)
	if err != nil {
		fmt.Println("Failed to get image data:", err)
		return
	}
	outputImageData, err := convertImageFormat(imageData, *outputFormat)
	if err != nil {
		fmt.Println("Failed to convert image format: ", err)
		return
	}
	outputFilePath := getOutputFilePath(*imageURL, *outputFormat)
	err = saveImage(outputImageData, outputFilePath)
	if err != nil {
		fmt.Println("Failed to save image:", err)
		return
	}

	fmt.Println("Image converted successfully. Output file:", outputFilePath)

}

/*
*
获取图片信息
*/
func getImageData(imageURL string) ([]byte, error) {
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		resp, err := http.Get(imageURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}

	//todo 暂时不支持本地文件上传
	return nil, nil

}

func convertImageFormat(imageData []byte, outputFormat string) ([]byte, error) {
	imageReader := strings.NewReader(string(imageData))
	image, err := webp.Decode(imageReader)
	//image, _, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	outputImageData := new(strings.Builder)
	err = jpeg.Encode(outputImageData, image, nil)
	if err != nil {
		return nil, err

	}
	return []byte(outputImageData.String()), nil

}

func getOutputFilePath(imageUrl, outputFormat string) string {
	fileName := filepath.Base(imageUrl)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return fileNameWithoutExt + "." + outputFormat

}

func saveImage(imageData []byte, filepath string) error {
	outputFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(imageData)
	return err

}
