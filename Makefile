# Run this build target regardless of if file exists or not
.PHONY: bananaWaffleCookies backend frontend test

# Default 
all: bananaWaffleCookies

frontend:
	npm --prefix frontend/ run build

backend:
	go -C backend/ build -o ../bananaWaffleCookies 

# Build frontend and Go backend
bananaWaffleCookies: frontend backend

# Install Node.js and golang packages 
install:
	npm --prefix frontend/ install 
	go -C backend/ install 

# Install & Update Node.js and golang dependencies
update: install
	npm --prefix frontend/ update
	go -C backend/ get -u ./

# Clean project
clean: check_clean
	# Clean packages
	npm --prefix frontend/ clean-install # Removes node_modules and reinstalls without updating package.lock
	go -C backend/ clean
	# Remove potentially stale executables & build artifacts
	rm -f bananaWaffleCookies 
	rm -rf frontend/dist

# Run tests
test:
	npm --prefix frontend/ test run
	go -C backend/ test -v ./db/
	go -C backend/ test -v ./handlers/
