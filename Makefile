.PHONY: all build test clean run help

# Default target
all: help

# Go lessons
LESSONS := s01-hello-agent s02-api-client s03-streaming s04-tool-interface \
           s05-agent-loop s06-multi-tools s07-config s08-tui \
           s09-prompt-system s10-coordinator s11-memory s12-mcp

# Build all lessons
build:
	@for lesson in $(LESSONS); do \
		echo "Building $$lesson..."; \
		cd go/$$lesson && go build -o bin/agent . && cd ../..; \
	done
	@echo "All lessons built!"

# Test all lessons
test:
	@for lesson in $(LESSONS); do \
		echo "Testing $$lesson..."; \
		cd go/$$lesson && go test -v ./... && cd ../..; \
	done

# Run a specific lesson (usage: make run LESSON=s01)
run:
	@if [ -z "$(LESSON)" ]; then \
		echo "Usage: make run LESSON=s01"; \
		exit 1; \
	fi
	cd go/$(LESSON)* && go run main.go

# Clean build artifacts
clean:
	@for lesson in $(LESSONS); do \
		rm -rf go/$$lesson/bin/; \
	done
	rm -rf web/.next/ web/out/ web/node_modules/
	@echo "Cleaned!"

# Install web dependencies
web-install:
	cd web && npm install

# Build web interface
web-build:
	cd web && npm run build

# Export web for GitHub Pages
web-export:
	cd web && npm run export

# Serve web locally
web-serve:
	cd web && npm run dev

# Help
help:
	@echo "Learn Go Agent Harness - Makefile Commands"
	@echo ""
	@echo "  make build        Build all Go lessons"
	@echo "  make test         Test all Go lessons"
	@echo "  make run LESSON=s01   Run specific lesson"
	@echo "  make clean        Remove build artifacts"
	@echo ""
	@echo "  make web-install  Install web dependencies"
	@echo "  make web-build    Build web interface"
	@echo "  make web-export   Export for GitHub Pages"
	@echo "  make web-serve    Serve web locally"
	@echo ""
	@echo "Examples:"
	@echo "  make run LESSON=s01   # Run s01-hello-agent"
	@echo "  make run LESSON=s08   # Run s08-tui"
