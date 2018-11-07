default: config.go server.go db.go
	go build -o ctfExec
deps:
	go get -u github.com/gorilla/mux
	go get -u github.com/flosch/pongo2
	go get -u gopkg.in/ini.v1
	go get -u github.com/lib/pq
clean:
	rm -f ctfExec
run: ctfExec
	./ctfExec
