# Build

help:
	@echo "See Makefile"
tidy:
	@bash -c "go mod tidy"
run:
	@bash -c "go run ./example ."
run-race:
	@bash -c "go run -race ./example ."
scan:
	@script/gosec.sh