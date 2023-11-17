FROM golang:latest

WORKDIR /app

COPY . /app

RUN ls -la
RUN make deps

RUN cp .env.example .env
RUN make client-build

EXPOSE 8083

CMD ["/app/build/tcc-client-application"]