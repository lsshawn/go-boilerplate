.PHONY: run css

run:
	@templ generate
	@go run cmd/main.go

.PHONY: css
css:
	tailwindcss -i assets/css/tailwind.css -o static/css/tailwind.css --minify

.PHONY: css-watch
css-watch:
	tailwindcss -i assets/css/tailwind.css -o static/css/tailwind.css --watch
