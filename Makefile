clean:
	rm -rf build/
	rm static/css/tailwind.css

setup:
	npm install
	go mod download

build:
	go generate
	go build -o build/

run: 
	go generate
	go run main.go

watch:
	find -name "*.go" \
			-or -name "*.css" \
			-not -name "tailwind.css" \
			-and -not -path "./node_modules/*" \
		| entr -r make run

build-image:
	docker build -t joniator/steam-hour-booster-ui:latest .

