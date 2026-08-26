FROM golang:1.22
RUN apt-get update && apt-get install -y curl && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && apt-get install -y nodejs && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY web/package.json web/package-lock.json ./web/
RUN cd web && npm install
COPY . .
RUN go build ./... && go test -count=1 ./... || true
RUN cd web && npm run build
CMD ["bash"]
