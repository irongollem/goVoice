# goVoice

## Description

goVoice is a voice driven call agent based on golang. It is currently in development and nearing a first BETA release.

## Installation

### Prerequisites

* [Go](https://golang.org/doc/install)
* [Telnyx](https://telnyx.com/) account (for receiving calls and interacting with the callee)
* [ngrok](https://ngrok.com/download) (optional for local testing)
* [openAI](https://openai.com/) account (for validating answers)
* [gmail](https://mail.google.com/) account (for sending the email)

### Install

```sh
git clone git@github.com:irongollem/goVoice.git
go mod download
```

Setup your .env file (see [.env.example](.env.example)).

## Usage

Run locally by running `go run cmd/server/main.go` or build and run the binary.
For running this in vscode (included in the repo) you might want to setup a `.vscode/launch.json` file with the following:

```sh
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "cwd": "${workspaceFolder}",
      "envFile": "${workspaceFolder}/.env",
    }
  ]
}
```

Then you can simply run and debug it from the console. Since most of it's behavior is responsive to incoming calls, you might want to use [ngrok](https://ngrok.com/) to expose your local server to the internet. You can do this by running `ngrok http 8080` (assuming you are running the server on port 8080). Then using the ngrok url as the webhook url in your telnyx dashboard followed by `/call` (e.g. `https://12345678.ngrok.io/call`).

The app itself is deployed on appEngine for now and can be deployed by running `gcloud app deploy` (see [app.yaml](app.yaml)).

## Project structure

```sh
root/
├── cmd/
│   └── server/
│       └── main.go # Entrypoint for the server
│
├── internal/
│   ├── api/
│   │   ├── voiceApi.go # HTTP handlers for the voice API
│   │   ├── webClientApi.go # HTTP handlers for the web client API
│   ├── app/
│   │   ├── audioProcessor/ # Currently not used
│   │   └── conversation/ # Conversation Controller (brain of the app)
│   │         ├── controller.go # The actual controller
│   │         └── helpers.go # Helper functions dealing with the conversation
│   ├── config/
│   │   └── config.go # Configuration and environment variables
│   ├── email/
│   │   └── email.go # Email client
│   └── models/
│       └── models.go # Any internal models
│
├── pkg/
│   ├── ai/
│   │   ├── openAi/
│   │   │   └── openAi.go # OpenAI client
│   │   ├── ai.go # AI interface
│   │   └── models.go # AI specific response model
│   ├── audio/
│   │   ├── audio.go # Audio interface
│   │   └── telnyx/
│   │       ├── commands_test.go
│   │       ├── commands.go # Telnyx commands to trigger the telnyx API
│   │       ├── commandStructs.go # Models describing the required body for the telnyx API
│   │       ├── dto.go # Models derscribing the response from the telnyx API
│   │       ├── helpers.go # Helper functions for the telnyx API
│   │       ├── procedures.go # Hook handlers
│   │       └── telnyx.go # API router for the incomming hooks implementing the audio voiceApi interface
│   ├── db/
│   │   ├── db.go # DB interface
│   │   └── firestore/
│   │       ├── conversationHandlers.go # Firestore specific conversation handlers
│   │       ├── firestore.go # Firestore client implementing the DB interface
│   │       └── rulesetHandlers.go # Firestore specific ruleset handlers
│   ├── host/
│   │   └── host.go # File dealing with security if incoming requests
│   ├── storage/
│   │   ├── storage.go # Storage interface
│   │   └── googleStorage.go # Google Storage client implementing the storage interface (bucket)
├── .gcloudignore
├── .gitignore
├── app.yaml
├── cloudbuild.yaml (not yet working)
├── go.mod
├── go.sum
├── LICENSE.md
├── pilot.json # (temporary file describing the initial pilot's ruleset)
└── README.md
```

## License

Proprietary License (see [LICENSE.MD](LICENSE.MD)).
