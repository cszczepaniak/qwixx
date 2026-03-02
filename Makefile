.PHONY: templ tailwind generate run build

# Regenerate Go code from .templ files
templ:
	templ generate

# Build Tailwind CSS into embedded static dir
tailwind:
	npx tailwindcss -c tailwind.server.config.js -i cmd/qwixx/static/input.css -o cmd/qwixx/static/dist.css

# Regenerate templ and tailwind (run before building the Go server)
generate: templ tailwind

# Build the Go server binary
build: generate
	go build -o qwixx ./cmd/qwixx

# Run the Go server (builds assets first)
run: generate
	go run ./cmd/qwixx
