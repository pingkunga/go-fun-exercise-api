FROM golang:1.22.1-alpine3.19 AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o /bin/app

FROM alpine:3.19.1

COPY --from=build /bin/app /bin

EXPOSE 1323

RUN adduser -D -u 2024 appuser
USER appuser

CMD [ "/bin/app" ]