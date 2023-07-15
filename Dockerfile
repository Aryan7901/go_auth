# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /auth_service

EXPOSE 3001
EXPOSE 587
CMD [ "/auth_service" ]