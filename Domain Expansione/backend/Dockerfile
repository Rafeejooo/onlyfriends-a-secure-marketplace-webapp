FROM golang:1.21-alpine

WORKDIR /app

# Copy mod files first to use Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

RUN go build -o main .

EXPOSE 8000
CMD ["./main"]
