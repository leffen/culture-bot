#host := 172.31.0.63
host := 192.168.1.49

rpi: build
	rsync --rsync-path="sudo rsync" rpi-client/culture-bot.conf pi@$(host):/etc/supervisor/conf.d/
	ssh pi@$(host) "sudo supervisorctl reread"
	ssh pi@$(host) "sudo supervisorctl update"
	ssh pi@$(host) "sudo supervisorctl restart culture-bot"

build:
	dep ensure
	sudo GOARCH=arm GOOS=linux go get github.com/blackjack/webcam
	GOARCH=arm GOOS=linux go build -o culture-bot-rpi rpi-client/*.go
	ssh pi@$(host) "mkdir -p ~/culture-bot"
	rsync culture-bot-rpi pi@$(host):~/culture-bot/
	rsync rpi-client/speech.json pi@$(host):~/culture-bot/