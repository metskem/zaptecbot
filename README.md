### zaptecbot

A Telegram bot to operate a zaptec charger (wallbox).  
Zaptac offers a [REST API](https://api.zaptec.com/help/index.html), that allows a ton of things, also starting and stopping/pausing the charger.  
This bot allows to:
* add a charging schedule (datetime+duration)
* list existing schedules
* delete a schedule
* start charging
* stop/pause charging

On startup, you provide the bot with the following envvars:
* BOT_TOKEN - The Telegram bot token
* ZAPTEC_CHARGER_ID - The id (guid) of your charger to operate
* ZAPTEC_USERNAME - Your Zaptec account (user name). This will be used to login and get an OAuth token for.  
* ZAPTEC_PASSWORD - The password for your Zaptec account
* DEBUG - True/false, turn bot debugging on (default is false)

On startup the bot will first initialize itself with Telegram, and after that it will try to login with the provided username/password, will get an OAuth token and will try to get some basic info from the charger (using the given charger_id).

### TimeZones
Since a Telegram client does not send it's timezone information to the bot, and we don't want to ask the clients for their timezone, you might have to set the TZ envvar (i.e. "Europe/Amsterdam") in order to allow the clients their local timezone.

## Supported Bot Commands

The following commands are supported (and can/should be configured with BotFather (Edit Commands) for convenience):
```
state - show the current state of the charger
debug - [on|off] - dynamically turn Telegram Bot debugging on/off
sl - (Schedule List) list the current schedules
sa - (Schedule Add) HH:mm n - Add a schedule, H=Hours, m=minutes, n=duration in hours. When the given time is before current time, we add one day (assuming you wanted that time for tomorrow)
sd - (Schedule Delete) HH:mm n - Delete a schedule, H=Hours, m=minutes, n=duration in hours. 
``` 
