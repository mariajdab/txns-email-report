# Build
FROM golang:1.18.1-alpine AS build
RUN apk --no-cache add git

WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./reporter

# Final container
FROM alpine

WORKDIR /

COPY --from=build  /usr/src/reporter /reporter
COPY --from=build /usr/src/.env /
COPY --from=build /usr/src/email.html /

RUN mkdir /input
RUN mkdir /output

ENTRYPOINT [ "./reporter" ]
