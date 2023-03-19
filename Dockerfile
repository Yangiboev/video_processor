FROM golang:1.19

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y ffmpeg

RUN go build -o video_processor

ENTRYPOINT ["/app/video_processor"]
