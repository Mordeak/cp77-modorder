BINARY := CP77-modorder.exe
CC     := C:/msys64/ucrt64/bin/gcc.exe

export CGO_ENABLED := 1
export CC

.PHONY: all build dev generate clean tidy

all: build

build:
	wails build -platform windows/amd64 -o $(BINARY)

dev:
	wails dev

generate:
	wails generate module

clean:
	$(RM) $(BINARY)
	$(RM) -rf frontend/dist

tidy:
	go mod tidy
	cd frontend && npm install
