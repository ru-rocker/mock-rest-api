# builder image
FROM golang:1.17-alpine as builder
RUN mkdir /build && apk --no-cache add build-base gcc
ADD *.go /build/
ADD parser/*.go /build/parser/
ADD go.mod /build/
WORKDIR /build
RUN go get . && go build -a -o mock-rest-api .


# generate clean, final image for end users
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /build/mock-rest-api .
ADD config/mock.yaml config/
RUN chmod +x mock-rest-api

EXPOSE 8080

# executable
ENTRYPOINT [ "./mock-rest-api" ]
# CMD ["3", "300"]