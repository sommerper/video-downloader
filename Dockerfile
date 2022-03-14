FROM golang:1.17 as dev
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp
RUN apt-get update && apt-get install -y ffmpeg
RUN chmod a+rx /usr/bin/ffmpeg

RUN go build -v main.go
CMD ["/usr/src/app/main"]
# ENTRYPOINT ["tail", "-f", "/dev/null"]