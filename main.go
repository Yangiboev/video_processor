package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
		if !file.FileInfo().IsDir() && strings.HasSuffix(strings.ToLower(file.Name), ".mp4") {
			input, err := file.Open()
			if err != nil {
				return err
			}
			defer input.Close()
			output := fmt.Sprintf("output_%s", file.FileInfo().Name())
			inputFile := filepath.Join(outputDir, file.FileInfo().Name())
			outputFile := filepath.Join(outputDir, output)

			err = saveFile(input, inputFile)
			if err != nil {
				return err
			}

			err = scaleAndSpeedUpVideo(inputFile, outputFile, width, height, speedFactor)
			if err != nil {
				return err
			}

			err = os.Remove(inputFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func saveFile(reader io.Reader, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
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
