```sh
crawling-system/
│
├── cmd/
│   └── main.go        # Application entry point
│
│
├── internal/
│   ├── scrap/           # Application-specific package for crawling
│   │   ├── url-manager/
│   │   ├── policy/
│   │   └── scrap.go
│   │
│   ├── extractor/         # Package to extract links from crawled HTML
│   │   ├── parser/
│   │   └── extractor.go
│   │
│   ├── processor/         # Package to manage statement for crawler instance
│   │   └── process.go
│   │
│   └── db/                # Internal package for database interaction
│       ├── model.go       # Defines the structure of crawled data
│       ├── crud.go        # CRUD operations on crawled data
│       └── duplicate.go   # Search for duplicated URLs
│
├── pkg/
│   └── system.go          # Package for system-level operations, e.g. logging
│
├── go.mod
├── go.sum
│
├── config.yaml        # General configuration file
└── db_config.yaml     # Database configuration file
```

## Explanation
1. `cmd/`

   - main.go
   This is the entry point of your application where the application is started.

2. `internal/`
Houses the main application logic, divided into several packages.

   - `scrap/scrap.go`
   Contains functions and logic for crawling websites.

   - `extractor/extractor.go`
   Houses functions and logic for extracting links from the downloaded HTML.

   - `validator/validator.go`
   Contains functions for validating URLs.

3. `pkg/models/`
Contains any shared models or structures that are used across different packages.

4. `go.mod`
Defines the module’s module path, which is also the import path used for the root directory, and its dependency requirements.