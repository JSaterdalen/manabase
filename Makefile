build:
	@echo "Building..."
	@make templ
	@make tailwind
	@go build -o bin/manabase ./cmd/manabase

run:
	./manabase
watch:
	air; \
	echo "Watching...";

tailwind:
	./tailwindcss -i web/static/css/input.css -o cmd/web/static/css/output.css

tailwind-watch:
	./tailwindcss -i web/static/css/input.css -o cmd/web/static/css/output.css --watch

templ:
	templ generate

templ-watch:
	templ generate -watch

dev:
	wgo run ./cmd/api :: \
	wgo -file .templ templ generate :: \
	wgo -file templ -file html ./tailwindcss -i cmd/web/static/css/input.css -o cmd/web/static/css/output.css
