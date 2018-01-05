# build stage
FROM golang:alpine AS build-env
COPY . /src
RUN  apk update && apk upgrade && \
     apk add --no-cache bash git openssh
RUN go get github.com/boltdb/bolt 
RUN go get gopkg.in/validator.v2
RUN cd /src && go build -o status-tracker

# Final container
FROM alpine

RUN mkdir /status-tracker
COPY html/ /status-tracker/html
COPY --from=build-env /src/status-tracker /status-tracker
EXPOSE 7080

WORKDIR /status-tracker

CMD ["./status-tracker"]
