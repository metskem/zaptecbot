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


## TODO
* How to authenticate (give it some authorized telegram accountids in an envvar?)
* Do we want to persist schedules, so they survive a bot restart?
* Should the bot notify when a schedule starts/stops? (I think it will then need some data to persist, like the subscriber chatid's)
* When adding a schedule, it should check for overlap with existing schedules.
* 