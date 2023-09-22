# goVoice

## Description

goVoice is a simple voice driven call agent. It is currently in development and is not ready for use.

## Installation

### Prerequisites

* [Go](https://golang.org/doc/install)
* [Twilio](https://www.twilio.com/try-twilio) account
* [ngrok](https://ngrok.com/download) (optional for local testing)

### Install

```bash
go mod download
```

## Usage

TBD

## Project structure

```bash
root/
├── cmd/
│   └── server/
│       └── main.go # Entrypoint for the server
│
├── internal/
│   ├── api/
│   │   ├── handlers/ # HTTP handlers
│   │   ├── middleware/ # HTTP middleware
│   │   └── routes/ # HTTP routes
│
├── pkg/
│   ├── audio/
│   │   ├── twilio/
│   │   │   ├── client.go # Twilio client
│   │   │   ├── recorder.go # Twilio recorder
│   │   │   └── transcriber.go # Twilio transcriber
│   │   └── google/
│   │       ├── client.go # Google client
│   │       ├── recorder.go # Google recorder
│   │       └── transcriber.go # Google transcriber
│   ├── llt/
│   │   ├── client.go # LLT client
│   ├── storage/
│   │   ├── dynamodb/ # placeholder until we chose a db provider
│   │   │   ├── client.go # DynamoDB client
│   │   ├── cloud/
│   │   │   ├── client.go # Cloud client for storage
│
├── config/
│   ├── config.go # Configuration
│
├── tests/
│   ├── integration/
│   ├── unit/
│
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
└── README.md
└── Dockerfile
```

## License

None, as of yet this is a propriatary project owned by CroCode BV, the Netherlands.
