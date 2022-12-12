# tigerd
Turn stuff into identity providers

## Protobuf

Tigerd exposes most of its API extension points in [Protobuf v3](https://developers.google.com/protocol-buffers/docs/proto3) 
format. Protobuf provides strongly typed API and models, and has good support across all major languages. However, it 
could be challenging when it comes to file distribution and code generation, as there's no standard solution.

Tigerd uses [buf.build](https://buf.build/) to centrally manage our protobuf definitions. SDKs are generated and distributed
automatically to avoid vendor-client discrepancies.

> As to writing, v1.10.0 is used for the buf cli.
