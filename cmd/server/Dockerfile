FROM golang:latest

WORKDIR /app

COPY . /app

RUN ls -la
RUN make deps

# TODO: adjust to make build correctly
RUN cp .env.example .env
RUN make server-build

# TODO: verify which por will run in Postman/Locust
EXPOSE 8081

# Passing flag to server command and setting http as default
# CMD ["/app/build/tcc-server-application", "-type=default-type"]

# ENTRYPOINT definindo o comando padrão
ENTRYPOINT ["/app/build/tcc-server-application"]

# CMD fornecendo argumentos padrão
CMD ["-type=default-type"]