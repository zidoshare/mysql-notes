FROM golang:1.14
WORKDIR /app
COPY go.mod go.sum ./
RUN  export GOPROXY=https://goproxy.cn && go mod download
COPY . .
ENTRYPOINT ["go","run","main.go"]