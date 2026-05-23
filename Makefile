CURRENT_DIR:=$(shell pwd)
LOCAL_BIN_INSTALL:=$(CURRENT_DIR)/.local/
LOCAL_BIN:=$(CURRENT_DIR)/.local/bin/
BACKEND_PATH:=$(CURRENT_DIR)/api/
FRONTEND_PATH:=$(CURRENT_DIR)/web-app/

.PHONY: run run\:% stop stop\:% go-deps go-mocks go-test go-run local-db-up local-db-down web-deps web-run

# DOCKER 

run:
	docker-compose up -d

run\:%:
	docker-compose up $* -d

stop:
	docker-compose down 

stop\:%:
	docker-compose down $* 

# BACKEND 

go-deps:
	@echo "installing go project binary dependencies..." &\
	GOPATH=$(LOCAL_BIN_INSTALL) go install github.com/pressly/goose/v3/cmd/goose@v3.27.0 &\
	GOPATH=$(LOCAL_BIN_INSTALL) go install github.com/vektra/mockery/v3@v3.7.0 &\
	wait

go-mocks: go-deps
	cd $(BACKEND_PATH) && \
		$(LOCAL_BIN)mockery

go-test: go-mocks
	cd $(BACKEND_PATH) && \
		go test ./... -v

go-run:
	cd $(BACKEND_PATH) && \
		go run ./cmd/api/main.go

local-db-up: go-deps 
	cd $(BACKEND_PATH) && \
		$(LOCAL_BIN)goose up

local-db-down: go-deps 
	cd $(BACKEND_PATH) && \
		$(LOCAL_BIN)goose down

# FRONTEND

web-deps:
	cd $(FRONTEND_PATH) && \
		npm ci

web-run: web-deps
	cd $(FRONTEND_PATH) && \
		npm run dev

