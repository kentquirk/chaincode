.PHONY: generate clean fuzz

fuzz:
	go test ./... --race -timeout 10s -short
	FUZZ_RUNS=50000 go test --race -v -timeout 30s ./pkg/vm -run "TestFuzzJunk"
	FUZZ_RUNS=50000 go test --race -v -timeout 30s ./pkg/vm -run "TestFuzzFunctions"
	FUZZ_RUNS=10000 go test --race -v -timeout 30s ./pkg/vm -run "TestFuzzValid"

generate: opcodes.md pkg/vm/opcodes.go pkg/vm/miniasmOpcodes.go pkg/vm/opcode_string.go \
		pkg/vm/extrabytes.go cmd/chasm/chasm.peggo pkg/vm/enabledopcodes.go

clean:
	rm cmd/opcodes/opcodes

opcodes.md: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --opcodes opcodes.md

pkg/vm/opcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --defs pkg/vm/opcodes.go

pkg/vm/miniasmOpcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --miniasm pkg/vm/miniasmOpcodes.go

pkg/vm/extrabytes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --extra pkg/vm/extrabytes.go

pkg/vm/enabledopcodes.go: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --enabled pkg/vm/enabledopcodes.go

cmd/chasm/chasm.peggo: cmd/opcodes/opcodes
	cmd/opcodes/opcodes --pigeon cmd/chasm/chasm.peggo

cmd/opcodes/opcodes: cmd/opcodes/*.go
	cd cmd/opcodes && go build

pkg/vm/opcode_string.go: pkg/vm/opcodes.go
	go generate ./pkg/vm

