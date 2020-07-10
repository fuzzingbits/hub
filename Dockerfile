FROM node:12-buster as nodeBuilder
WORKDIR /project
COPY . .
RUN git clean -Xdf
RUN make full-ui

FROM golang:1.14-buster as goBuilder
WORKDIR /project
COPY . .
RUN git clean -Xdf
COPY --from=nodeBuilder /project/dist/ /project/dist/
RUN make full-go

FROM debian:buster
COPY --from=goBuilder /project/var/hub /usr/local/bin/
CMD ["hub"]
EXPOSE 2020
