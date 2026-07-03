FROM golang:1.25

RUN apt-get update && apt-get install -y bash curl jq && rm -rf /var/lib/apt/lists/*

COPY . /tests
WORKDIR /tests

RUN go mod download

RUN go install github.com/tsenart/vegeta@latest
RUN chmod +x rate_limit_test.sh

ENTRYPOINT ["sh", "-c", "go test -v ./... && ./rate_limit_test.sh"]
