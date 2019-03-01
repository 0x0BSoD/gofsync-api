lint:
	gofmt apiActions.go > apiActions.go_f
	mv apiActions.go_f apiActions.go
	gofmt hostgroups.go > hostgroups.go_f
	mv hostgroups.go_f hostgroups.go
	gofmt puppetclasses.go > puppetclasses.go_f
	mv puppetclasses.go_f puppetclasses.go
	gofmt main.go > main.go_f
	mv main.go_f main.go
	gofmt dbActions.go > dbActions.go_f
	mv dbActions.go_f dbActions.go

build:
	go build -x -o fGetter *.go
	mv fGetter ./bin
