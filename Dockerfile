FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simulationgame

RUN export SIM_GAME_DB_PW=valmet865

CMD ["/simulationgame"]
