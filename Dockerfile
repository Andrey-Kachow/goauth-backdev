FROM golang:1.23.1

WORKDIR /app

COPY go.mod ./

# Conditionally copy go.sum if it exists
# COPY go.sum ./
# RUN go mod download

COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-goauth-exe cmd/mainapp
RUN go build -o ./docker-goauth-exe ./cmd/mainapp

EXPOSE 8080

CMD ["./docker-goauth-exe"]