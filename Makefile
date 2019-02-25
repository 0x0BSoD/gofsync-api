lint:
	gofmt hostgroups.go > hostgroups.go_f
	mv hostgroups.go_f hostgroups.go
	gofmt puppetclasses.go > puppetclasses.go_f
	mv puppetclasses.go_f puppetclasses.go
	gofmt main.go > main.go_f
	mv main.go_f main.go
	gofmt dbActions.go > dbActions.go_f
	mv dbActions.go_f dbActions.go

build:
	go build -o fGetter main.go dbActions.go hostgroups.go puppetclasses.go
	mv fGetter ./bin
