build:
	go generate
	go build

run: 
	go generate
	go run main.go

watch:
	find -name "*.go" \
			-or -name "*.css" \
			-not -name "tailwind.css" \
			-and -not -path "./node_modules/*" \
		| entr -r make run

