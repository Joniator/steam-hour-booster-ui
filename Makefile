clean:
	rm -rf build/
	rm web/static/css/tailwind.css

setup:
	npm install --prefix ./web
	go mod download

build:
	go generate ./web
	go build -o build/

ci-build:
	go generate
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o build/steam-hour-booster-ui-$(GOOS)-$(GOARCH)

run: 
	go generate ./web
	go run ./cmd/steam-hour-booster-ui/main.go -u test -p test

watch:
	find -name "*.go" \
			-or -name "*.css" \
			-not -name "tailwind.css" \
			-and -not -path "*/node_modules/*" \
	| entr -r make run

build-image:
	docker build -t joniator/steam-hour-booster-ui:latest .

