# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN  apk update && apk upgrade && \
     apk add --no-cache bash git openssh
RUN go get github.com/boltdb/bolt 
RUN go get gopkg.in/validator.v2
RUN cd /src && go build -o status-tracker

# Final container
FROM alpine
MAINTAINER <jim@jimturpin.com>

COPY --from=build-env /src/status-tracker /
EXPOSE 7080

CMD ["/status-tracker"]
