all: zip

windows-amd64: gakujo-google-calendar_windows_amd64.exe
linux-amd64: gakujo-google-calendar_linux_amd64
darwin-amd64: gakujo-google-calendar_darwin_amd64
darwin-arm64: gakujo-google-calendar_darwin_arm64

gakujo-google-calendar_windows_amd64.exe: credentials.json
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build \
		-ldflags \
		-H=windowsgui \
		-o $@ \
		.

gakujo-google-calendar_linux_amd64: credentials.json Dockerfile-linux_amd64
	docker build -t gakujo-google-calendar-linux-amd64:latest -f Dockerfile-linux_amd64 .
	docker run \
		-it \
		-v $(CURDIR):/work \
		-e CGO_ENABLED=1 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		gakujo-google-calendar-linux-amd64:latest \
			go build \
				-o $@ \
				.

gakujo-google-calendar_darwin_amd64: credentials.json
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	CC=clang \
	go build \
		-o $@ \
		.

gakujo-google-calendar_darwin_arm64: credentials.json
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=arm64 \
	CC=clang \
	go build \
		-o $@ \
		.

%.zip: % README.md
	zip -r $@ $^

gakujo-google-calendar_windows_amd64.zip: gakujo-google-calendar_windows_amd64.exe README.md
	zip -r $@ $^

zip: gakujo-google-calendar_windows_amd64.zip gakujo-google-calendar_linux_amd64.zip gakujo-google-calendar_darwin_amd64.zip gakujo-google-calendar_darwin_arm64.zip

clean:
	rm -f gakujo-google-calendar_*_*