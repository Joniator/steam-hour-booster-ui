clean:
	rm -rf build/
	rm static/css/tailwind.css

setup:
	npm install
	go mod download

build:
	go generate
	go build -o build/

ci-build:
	go generate
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/steam-hour-booster-ui-$(GOOS)-$(GOARCH)

run: 
	go generate
	go run main.go -u test -p test

watch:
	find -name "*.go" \
			-or -name "*.css" \
			-not -name "tailwind.css" \
			-and -not -path "./node_modules/*" \
	| entr -r make run

build-image:
	docker build -t joniator/steam-hour-booster-ui:latest .

