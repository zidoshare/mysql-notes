FROM maven:3-openjdk-8
WORKDIR /app
COPY pom.xml ./
RUN  export GOPROXY=https://goproxy.cn && go mod download
COPY . .
ENTRYPOINT ["go","run","main.go"]