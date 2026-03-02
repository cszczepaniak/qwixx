.PHONY: templ tailwind generate run build

# Regenerate Go code from .templ files
templ:
	go tool templ generate

# Build Tailwind CSS into embedded static dir
tailwind:
	tailwindcss -c tailwind.config.js -i cmd/qwixx/static/input.css -o cmd/qwixx/static/dist.css

format:
	go tool templ fmt .

# Regenerate templ and tailwind (run before building the Go server)
generate: templ tailwind format

# Build the Go server binary
build: generate
	go build -o qwixx ./cmd/qwixx

# Run the Go server (builds assets first)
run: generate
	go run ./cmd/qwixx
