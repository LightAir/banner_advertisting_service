FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN apk --no-cache add bash

COPY . .

ARG BIN_TEST_FILE
ARG CONFIG_FILE

RUN CGO_ENABLED=0 go test -c -v --tags=integration -o ${BIN_TEST_FILE} /app/tests

RUN echo "sleep 15; ${BIN_TEST_FILE} --config ${CONFIG_FILE}" > /opt/executable.sh
RUN chmod +x /opt/executable.sh

CMD ["/bin/bash", "-c", "/opt/executable.sh"]