NAME=star-rail-fps-unlocker
OUTDIR=build
PKGNAME=github.com/gizmo-ds/star-rail-fps-unlocker
MAIN=main.go
FLAGS+=-trimpath
FLAGS+=-tags timetzdata
FLAGS+=-ldflags "-s -w"
export CGO_ENABLED=0

all: windows-amd64

generate:
	go generate ./...

windows-amd64: generate
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME).exe $(MAIN)

sha256sum:
	cd $(OUTDIR); for file in *; do sha256sum $$file > $$file.sha256; done