FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN apk --no-cache add bash

COPY . .

ARG BIN_FILE
ARG CONFIG_FILE

RUN go build -o ${BIN_FILE} ./cmd/banner

RUN echo "sleep 13; ${BIN_FILE} --config ${CONFIG_FILE}" > /opt/executable.sh
RUN chmod +x /opt/executable.sh

CMD ["/bin/bash", "-c", "/opt/executable.sh"]