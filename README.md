# camerabot
## Software
Telegram bot that makes a photo and sends it to chat. 

I use this bot to monitor kiln temperature and make sure workshop is not on fire yet.
Go part is responsible for interacting with Telegram API. Application uses long polling because in my case device is located behind two NATs. 
Uses `raspistill` (via os.exec) to make photos.
Parametrized commands for raspistill are stored in external bash scripts.

### Running bot
1. Setup Raspberry Pi and Pi camera.
2. Set environment variable `TOKEN` with your bot token (Botfather can provide you with the one).
3. Use systemd config to start as a service or simply run the app from the console.
4. Start direct conversation with bot or add bot to group chat if yu are interested in broadcasting your kiln paranoia.

### Commands
1. `/pi` sends ordinary photo.
2. `/zoom` sends zoomed and croped region of interest. Kiln controller in my case.

## Hardware
Currently runs on a Raspberry Pi. Using onboard V2 camera.
