build:
	gofmt main.go > main.go_f
	mv main.go_f main.go
	gofmt dbActions.go > dbActions.go_f
	mv dbActions.go_f dbActions.go
	go build -o fGetter main.go dbActions.go
	mv fGetter ./bin
