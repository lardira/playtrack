# playtrack
Compete with friends on randomly assigned games 
to see who's the fastest gamer across any genre.

## Quick Start

### Prerequisites

- Go 1.25+
- goose
- mockery
- docker

### Installation

1. Clone the repository
```bash
git clone https://github.com/lardira/playtrack
cd ./playtrack
```

2. Copy `.env.template` and set environment variables in `.env` files

3. Start docker services 
```
docker compose up -d
```

4. Install dependencies and run

```bash
go mod download; go run cmd/api/main.go;
```
