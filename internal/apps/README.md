## :bookmark_tabs: Applications codes

The directory `apps` contains all applications defined and divided as *client* and *server* and by communication method name. 

:bulb: There's also some files acting as helpers on *client* folder and gRPC related on server folder.

### *client*

- `bff.go`: BFF application code with interaction endpoints.
- `grpc.go`: helper functions to handle interaction with server using gRPC.
- `http.go`: helper functions to handle interaction with server using HTTP (REST).
- `rabbitmq.go`: helper functions to handle interaction with server using RabbitMQ.

### *server*

- `consumer.go`: RabbitMQ queue message consumer code.
- `grpc.go`: gRPC server code.
- `http.go`: HTTP (REST) server code.
- `message.proto`: Protocol Buffers (Protobuff) file.
- `message_grpc.pb.go` and `message.pb.go`: gRPC package generated code based on Probuffs defined at `message.proto`.
- `commom.go`: common interfaces for use cases.
- `setup.go`: *Setup Application* code to interact with database and provide some facilities for local testing and analysis.

:pill: All the servers and RabbitMQ message consumer call the use cases in  almost same conditions since use cases are agnostic and it's used dependency inversion too.