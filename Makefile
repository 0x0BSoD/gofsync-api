build:
	gofmt main.go > main.go_f
	mv main.go_f main.go
	go build -o fGetter main.go
	mv fGetter ./bin
