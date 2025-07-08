.PHONY: css cleancss templ build run clean

css:
	npx @tailwindcss/cli -i ./css/input.css -o ./assets/css/main.css --minify

cleancss:
	rm -rf ./assets/css/main.css

templ:
	go tool templ generate

build: css templ
	go build -o ./school ./main.go

run: css templ
	go run ./main.go

clean: cleancss
	rm -rf ./school