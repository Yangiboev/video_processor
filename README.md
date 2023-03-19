# Video Processor

Video Processor is a command-line tool for processing videos in a zip archive. It scales and speeds up the videos using the specified width, height, and speed factor.

## Prerequisites

- Go 1.17 or later
- FFmpeg

## Installation

1. Clone the repository:

```sh
git clone https://github.com/your-username/video-processor.git
cd video-processor
go build -o video_processor
```

## Usage
```sh
./video_processor -zip videos.zip -output output_dir -width 540 -height 410 -speed 0.5
```

- zip: Path to the zip file containing video files.
- output: Output directory for processed videos.
- width: Output video width (default: 540).
- height: Output video height (default: 410).
- speed: Speed factor (default: 0.5).

## Docker Usage
You can also run the video processor using Docker:

## Build the Docker image:
```sh
docker build -t video_processor .
```
Run the Docker container with the required command-line arguments:
```sh
docker run -v /path/to/your/videos.zip:/videos.zip -v /path/to/output/dir:/output video_processor -zip /videos.zip -output /output -width 540 
-height 410 -speed 0.5
```
Replace /path/to/your/videos.zip with the absolute path to the zip file on your local machine, and /path/to/output/dir with the absolute path to the directory where you want to save the processed videos.