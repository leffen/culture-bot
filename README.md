### Culture Bot

This bot uses Facial Recognition to recognize people who come to the office and interact with them.

 - [User Guide](https://github.com/plutov/culture-bot/wiki/User-Guide)
 - [Training Process](https://github.com/plutov/culture-bot/wiki/Training-Process)
 - [Culture Bot Architecture](https://docs.google.com/drawings/u/1/d/1XgNSmMCMUwQ7xErBtiVWeICr4ZOjmmu0CbWOsTpWVMo/edit?usp=drive_web)

### Integrations

 - Slack
 - Google Calendar
 - Google Speech

### Components

 - Facebox: http://localhost:8080
 - Textbox: http://localhost:8082
 - Web Dashboard: http://localhost:8081
 - Web Client: http://localhost:8081/client
 - Raspberry Pi Client: `./culture-bot-rpi`

### Web Dashboard

The web dashboard is running on http://localhost:8081

### Calendar

To get API token we need to have a public redirect URI, for this we can use ngrok:
```
ngrok http 8081
```
Then visit this URL to get token and cache it: http://localhost:8081/calendar/token. Cache token will be expired, so we'll need to go to this URL again later.

### Person properties

 - Face training model
 - First Name
 - Last Name
 - Slack Display Name
 - Pronounced name

### Running locally

```
docker-compose build
docker-compose up -d
```

`ml-state` folder contains states of face detection model. We need to import this state on http://localhost:8080.

### Raspberry Pi 3 Client

Raspberry Pi 3 client uses V4L2 Linux framework to capture image from webcam.
Raspberry Pi 3 requirements:
 - wi-fi connection
 - ssh interface
 - USB webcam
 - USB speaker
 - USB microphone

To record audio we are using `rec(sox)` command line tool. Saved audio files are stored in `record/` folder.
```
sudo apt-get install sox
```

Build and run using Supervisord on RPi device:
```
make rpi
```

Client uses Google TTS to play sentences, senteces are cached in `audio/` folder.

### Contribution

Run tests:
```
dep ensure
go test -v ./...
```