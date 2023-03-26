package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	zipFile := flag.String("zip", "", "Zip file containing video files")
	outputDir := flag.String("output", "", "Output directory for processed videos")
	width := flag.Int("width", 540, "Output video width")
	height := flag.Int("height", 410, "Output video height")
	speedFactor := flag.Float64("speed", 0.5, "Speed factor")

	flag.Parse()

	if *zipFile == "" || *outputDir == "" {
		flag.Usage()
		log.Fatal("Error: Both zip file and output directory must be provided.")
	}

	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	err = processVideosInZip(*zipFile, *outputDir, *width, *height, *speedFactor)
	if err != nil {
		log.Fatalf("Error processing videos: %v", err)
	} else {
		fmt.Println("All videos processed successfully!")
	}
}

func processVideosInZip(zipFile, outputDir string, width, height int, speedFactor float64) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := file.Name
		if file.FileInfo().IsDir() {
			if _, err = os.Stat(filePath); os.IsNotExist(err) {
				err = os.MkdirAll(filePath, 0755)
				if err != nil {
					return err
				}
			}
			continue
		}

		input, err := file.Open()
		if err != nil {
			return err
		}
		defer input.Close()
		if err = saveFile(input, file.Name); err != nil {
			return err
		}

		if !file.FileInfo().IsDir() && strings.HasSuffix(strings.ToLower(file.Name), ".mp4") {
			lastIndex := strings.LastIndex(filePath, "/")
			fileDir := filePath[:lastIndex]
			fileBaseName := strings.ReplaceAll(file.FileInfo().Name(), ".mp4", "")

			output := fmt.Sprintf("%s/%s_%d_%d_%.2f.mp4", fileDir, fileBaseName, width, height, speedFactor)
			if err = scaleAndSpeedUpVideo(file.Name, output, width, height, speedFactor); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveFile(reader io.Reader, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func scaleAndSpeedUpVideo(inputFile, outputFile string, width, height int, speedFactor float64) error {
	size := []string{fmt.Sprintf("%d:%d", width, height)}
	speed := []string{fmt.Sprintf("%.2f*PTS", speedFactor)}

	args := ffmpeg.Input(inputFile).
		Filter("scale", size).
		Filter("setpts", speed).
		Output(outputFile, ffmpeg.KwArgs{"an": ""}).OverWriteOutput().ErrorToStdOut()

	return args.Run()
}
