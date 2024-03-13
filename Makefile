clean:
	rm -rf build/
	rm -rf web/static/dist
	rm -rf web/node_modules
	rm -rf internal/testdata/working*

setup:
	npm install --prefix ./web
	go mod download

build:
	go generate ./web
	go build -o build/steam-hour-booster-ui cmd/main.go

ci-build:
	go generate ./web
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/steam-hour-booster-ui-$(GOOS)-$(GOARCH) cmd/main.go

test:
	go test ./internal

run: 
	go generate ./web
	go run ./cmd/main.go -u test -p test -P 8123

watch:
	find -name "*.go" \
			-or -name "*.html" \
			-or -name "*.css" \
			-not -name "tailwind.css" \
			-and -not -path "*/node_modules/*" \
	| entr -r make run

build-image:
	docker build -t joniator/steam-hour-booster-ui:latest .
