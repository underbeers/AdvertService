FROM golang:1.18-alpine AS build

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o /build/advert_service /build/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY ./conf ./conf
COPY ./service.json ./service.json
COPY --from=build /build/advert_service .

ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_HOST=$POSTGRES_HOST
ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PORT=$POSTGRES_PORT
ENV POSTGRES_DB_NAME=$POSTGRES_DB_NAME
ENV GATEWAY_PORT=$GATEWAY_PORT
ENV GATEWAY_IP=$GATEWAY_IP
ENV ADVERTSERVICE_PORT=$ADVERTSERVICE_PORT
ENV ADVERTSERVICE_IP=$ADVERTSERVICE_IP

EXPOSE $ADVERTSERVICE_PORT

CMD [ "./advert_service" ]