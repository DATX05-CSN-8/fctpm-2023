OUT_DIR=out

clean:
	rm -rf $(OUT_DIR)

$(OUT_DIR):
	mkdir -p $(OUT_DIR)

build: $(OUT_DIR)
	go build -o $(OUT_DIR) ./cmd/orchestrator

run:
	go run ./cmd/orchestrator $(ARGS)

build-demo: $(OUT_DIR)
	go build -o $(OUT_DIR) ./cmd/demo
	sudo setcap cap_sys_admin+ep ./$(OUT_DIR)/demo
	sudo chmod g+rw /dev/vtpmx