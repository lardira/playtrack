# playtrack
Compete with friends on randomly assigned games 
to see who's the fastest gamer across any genre.

## Quick Start

### Prerequisites

- Go 1.26+
- docker
- npm

### Installation

1. Clone the repository
```bash
git clone https://github.com/lardira/playtrack 
cd ./playtrack
```

2. Copy `.env.template` and set environment variables in `.env` file.
```bash
cp .env.template .env
```

3. Start services
```bash
make run
```
or 
```bash
docker compose up -d
```

### Local development

#### Backend API
1. Copy `.env.template` and set environment variables in `.env` file.
```bash
cd ./api 
cp .env.template .env
```

2. Run the project
```bash
make go-run
```
or
```bash
go run ./cmd/api/main.go
```

#### Frontend app
1. Copy `.env.template` and set environment variables in `.env` file.
```bash
cd ./web-app 
cp .env.template .env
```

2. Run the project
```bash
make web-run
```
or
```bash
npm run dev
```