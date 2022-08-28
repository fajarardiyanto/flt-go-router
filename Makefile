# Build

help:
	@echo "See Makefile"
tidy:
	@bash -c "go mod tidy"
run:
	@bash -c "go run ./example ."
scan:
	@script/gosec.sh