GOBIN =../../bin

build:
	go build

windows:
	$(GOBIN)/rsrc -manifest app.manifest -ico assets/StarRealms.ico -arch amd64
	GOOS=windows go build -ldflags "-H=windowsgui"

all: build windows

clean:
	rm starRealmsNotify* rsrc.syso

.PHONY: build windows all clean
