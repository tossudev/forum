.PHONY: run reset

run:
	go run -tags sqlite_fts5 ./cmd

reset:
	go run -tags sqlite_fts5 ./cmd -reset
