OUT_DIR=out

clean:
	rm -rf $(OUT_DIR)

build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR) ./cmd/orchestrator

run:
	go run ./cmd/orchestrator $(ARGS)