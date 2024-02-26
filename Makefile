build:
	go generate
	go build

run: 
	go generate
	go run main.go

watch:
	ls styles/*.css templates/*.html | entr -r make run
