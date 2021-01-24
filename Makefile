build-mips-74Kc:
	echo "Building for MIPS 74Kc";
	rm -rf owbot_install
	mkdir -p owbot_install
	env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o owbot_install/owbot -trimpath -ldflags="-s -w" cmd/main.go && upx owbot_install/owbot;
	mkdir -p owbot_install/resources
	cp resources/config.json owbot_install/resources/config.json
	cp -p bin/owbot.service owbot_install
	cp -p bin/install.sh owbot_install
	cp -p bin/uninstall.sh owbot_install



