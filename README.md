### zaptecbot

A Telegram bot to operate a zaptec charger (wallbox).  
Zaptac offers a [REST API](https://api.zaptec.com/help/index.html), that allows a ton of things, also starting and stopping/pausing the charger.  
This bot allows to:
* add a charging schedule (datetime+duration)
* list existing schedules
* delete a schedule
* list the current state of the charger

On startup, you provide the bot with the following envvars:
* BOT_TOKEN - The Telegram bot token
* ZAPTEC_CHARGER_ID - The id (guid) of your charger to operate
* ZAPTEC_USERNAME - Your Zaptec account (user name). This will be used to login and get an OAuth token for.  
* ZAPTEC_PASSWORD - The password for your Zaptec account
* DEBUG - True/false, turn bot debugging on (default is false)

On startup the bot will first initialize itself with Telegram, and after that it will try to login with the provided username/password, will get an OAuth token and will try to get some basic info from the charger (using the given charger_id).

### TimeZones
Since a Telegram client does not send its timezone information to the bot, and we don't want to ask the clients for their timezone, you might have to set the TZ envvar (i.e. "Europe/Amsterdam") in order to allow the clients their local timezone.

## Supported Bot Commands

The following commands are supported (and can/should be configured with BotFather (Edit Commands) for convenience):
```
state - show the current state of the charger
debug - [on|off] - dynamically turn Telegram Bot debugging on/off
sl - (Schedule List) list the current schedules
sa - (Schedule Add) HH:mm n - Add a schedule, H=Hours, m=minutes, n=duration in hours. When the given time is before current time, we add one day (assuming you wanted that time for tomorrow)
sd - (Schedule Delete) HH:mm n - Delete a schedule, H=Hours, m=minutes, n=duration in hours. 
``` 

### Testing

If you want to mess manually with the API, here are some handy curl's: 
To get an OAuth token:
```bash
export TOKEN=$(curl -s -X POST https://api.zaptec.com/oauth/token -H "content-type: application/x-www-form-urlencoded" --data-raw 'grant_type=password&username=<ENCODED-USER>>&password=<ENCODE-PASSWORD>' | jq -r .access_token)
```

## State of the charger
````bash
curl -s "https://api.zaptec.com/api/chargers/<charger-guid>/state" -H  "Accept: text/plain" -H  "Authorization: Bearer $TOKEN" |jq
````

This will give a response like:
````json
[
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": -2,
    "Timestamp": "2024-01-27T07:58:43.52",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": -1,
    "Timestamp": "2024-03-09T17:13:38.907"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 1,
    "Timestamp": "2023-05-23T16:09:28.847",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 100,
    "Timestamp": "2023-06-27T14:44:59.167",
    "ValueAsString": "{\"DeviceType\":\"Go\",\"SerialNumber\":\"ZAP025320\",\"MeterCalibrated\":false}"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 120,
    "Timestamp": "2023-05-22T17:57:06.05",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 145,
    "Timestamp": "2023-05-22T19:11:55.773",
    "ValueAsString": "3600"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 147,
    "Timestamp": "2024-03-09T17:03:38.123",
    "ValueAsString": "600"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 150,
    "Timestamp": "2023-05-23T16:10:51.387",
    "ValueAsString": "LTE"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 151,
    "Timestamp": "2024-03-02T16:14:50.473",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 153,
    "Timestamp": "2024-03-02T16:59:28.833",
    "ValueAsString": "0.106"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 201,
    "Timestamp": "2024-03-09T17:01:06.09",
    "ValueAsString": "26.089"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 202,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "16.448"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 204,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "16.306"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 205,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "17.241"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 206,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "15.420"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 207,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "16.279"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 270,
    "Timestamp": "2024-03-09T17:01:06.093",
    "ValueAsString": "21.178"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 501,
    "Timestamp": "2024-03-09T17:02:13.153",
    "ValueAsString": "215.570"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 502,
    "Timestamp": "2024-03-09T17:02:13.153",
    "ValueAsString": "233.088"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 503,
    "Timestamp": "2024-03-09T17:02:13.153",
    "ValueAsString": "230.044"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 507,
    "Timestamp": "2024-03-09T17:03:39.223",
    "ValueAsString": "0.696"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 508,
    "Timestamp": "2024-03-09T17:03:39.223",
    "ValueAsString": "0.028"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 509,
    "Timestamp": "2024-03-09T17:03:39.223",
    "ValueAsString": "0.024"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 510,
    "Timestamp": "2024-03-09T17:01:17.127",
    "ValueAsString": "0.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 511,
    "Timestamp": "2023-05-22T17:51:00.76",
    "ValueAsString": "6.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 513,
    "Timestamp": "2024-03-09T17:03:39.223",
    "ValueAsString": "0.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 519,
    "Timestamp": "2024-03-09T11:46:14.467",
    "ValueAsString": "4"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 520,
    "Timestamp": "2023-05-22T17:51:00.76",
    "ValueAsString": "3"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 522,
    "Timestamp": "2023-05-22T18:06:08.737",
    "ValueAsString": "4"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 523,
    "Timestamp": "2024-03-09T10:03:24.24",
    "ValueAsString": "0.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 544,
    "Timestamp": "2023-05-22T19:11:55.773",
    "ValueAsString": "2"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 545,
    "Timestamp": "2023-05-22T17:50:58.277",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 546,
    "Timestamp": "2023-05-22T17:55:13.98",
    "ValueAsString": "40.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 547,
    "Timestamp": "2024-03-09T11:26:41.54",
    "ValueAsString": "28.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 548,
    "Timestamp": "2023-05-22T17:51:00.35",
    "ValueAsString": "4"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 553,
    "Timestamp": "2024-03-09T17:03:37.213",
    "ValueAsString": "0.234"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 554,
    "Timestamp": "2024-03-09T12:00:00.09",
    "ValueAsString": "OCMF|{\"FV\":\"1.0\",\"GI\":\"ZAPTEC GO\",\"GS\":\"ZAP025320\",\"GV\":\"2.3.1.2\",\"PG\":\"F1\",\"RD\":[{\"TM\":\"2024-03-09T12:00:00,000+00:00 R\",\"RV\":1546.403,\"RI\":\"1-0:1.8.0\",\"RU\":\"kWh\",\"RT\":\"AC\",\"ST\":\"G\"}]}"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 702,
    "Timestamp": "2024-03-09T17:04:39.263",
    "ValueAsString": "9"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 708,
    "Timestamp": "2024-03-09T17:03:23.193",
    "ValueAsString": "25.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 710,
    "Timestamp": "2024-03-09T17:04:21.243",
    "ValueAsString": "2"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 711,
    "Timestamp": "2023-05-22T17:57:06.057",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 712,
    "Timestamp": "2024-03-09T17:03:22.183",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 714,
    "Timestamp": "2024-03-09T15:47:31.597",
    "ValueAsString": "32"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 715,
    "Timestamp": "2023-05-22T17:50:58.273",
    "ValueAsString": "4"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 717,
    "Timestamp": "2023-06-06T11:14:48.133",
    "ValueAsString": "\r\n 4: VG:1.27 L12:417.65"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 718,
    "Timestamp": "2024-03-09T17:04:20.233",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 721,
    "Timestamp": "2024-03-09T15:47:32.627",
    "ValueAsString": "8d845c4a-7331-47c5-b308-f64f942ad17a"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 722,
    "Timestamp": "2024-03-09T11:47:24.91"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 723,
    "Timestamp": "2024-03-09T11:47:25.317",
    "ValueAsString": "{\"SessionId\":\"b1975e75-67a1-4d0b-a6ed-60b0d23c579c\",\"Energy\":8.715,\"StartDateTime\":\"2024-03-09T09:54:55.147849Z\",\"EndDateTime\":\"2024-03-09T11:47:24.919002Z\",\"ReliableClock\":true,\"StoppedByRFID\":false,\"AuthenticationCode\":\"\",\"SignedSession\":\"OCMF|{\\\"FV\\\":\\\"1.0\\\",\\\"GI\\\":\\\"ZAPTEC GO\\\",\\\"GS\\\":\\\"ZAP025320\\\",\\\"GV\\\":\\\"2.3.1.2\\\",\\\"PG\\\":\\\"T1\\\",\\\"RD\\\":[{\\\"TM\\\":\\\"2024-03-09T09:54:55,000+00:00 R\\\",\\\"TX\\\":\\\"B\\\",\\\"RV\\\":1537.688,\\\"RI\\\":\\\"1-0:1.8.0\\\",\\\"RU\\\":\\\"kWh\\\",\\\"RT\\\":\\\"AC\\\",\\\"ST\\\":\\\"G\\\"},{\\\"TM\\\":\\\"2024-03-09T11:00:00,000+00:00 R\\\",\\\"TX\\\":\\\"T\\\",\\\"RV\\\":1542.951,\\\"RI\\\":\\\"1-0:1.8.0\\\",\\\"RU\\\":\\\"kWh\\\",\\\"RT\\\":\\\"AC\\\",\\\"ST\\\":\\\"G\\\"},{\\\"TM\\\":\\\"2024-03-09T11:47:24,000+00:00 R\\\",\\\"TX\\\":\\\"E\\\",\\\"RV\\\":1546.403,\\\"RI\\\":\\\"1-0:1.8.0\\\",\\\"RU\\\":\\\"kWh\\\",\\\"RT\\\":\\\"AC\\\",\\\"ST\\\":\\\"G\\\"}]}\"}"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 732,
    "Timestamp": "2023-11-25T12:07:14.757",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 733,
    "Timestamp": "2023-05-23T16:09:28.847",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 749,
    "Timestamp": "2023-05-22T18:06:08.75",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 751,
    "Timestamp": "2023-05-22T18:03:47.227",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 800,
    "Timestamp": "2023-05-22T17:57:06.057",
    "ValueAsString": "037cd0c2-5969-4002-a298-f63e3769752e"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 801,
    "Timestamp": "2023-05-22T17:51:00.773",
    "ValueAsString": "default"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 802,
    "Timestamp": "2023-05-22T17:57:06.06",
    "ValueAsString": "zaptecje"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 803,
    "Timestamp": "2023-05-22T18:57:56.057",
    "ValueAsString": "1179648"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 804,
    "Timestamp": "2023-11-25T12:09:48.113",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 805,
    "Timestamp": "2023-05-22T17:57:06.06",
    "ValueAsString": "0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 807,
    "Timestamp": "2024-03-02T16:48:48.263",
    "ValueAsString": "#5 BLE restart - OK"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 808,
    "Timestamp": "2024-03-09T17:01:06.097",
    "ValueAsString": "7d 00h03m44s T_EM: 16.45 16.31 17.24  T_M: 15.42 16.28   V: 197.02 201.46 200.32   I: 2.29 0.02 0.02  0.254kW 0.000kWh C6 CM3 MCnt:487896 Rs:21 Rc:0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 809,
    "Timestamp": "2024-03-02T16:56:09.893",
    "ValueAsString": "60.000"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 811,
    "Timestamp": "2024-03-02T16:57:41.85",
    "ValueAsString": "6"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 815,
    "Timestamp": "2024-03-02T16:57:41.85",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 823,
    "Timestamp": "2023-05-22T18:58:01.067",
    "ValueAsString": "8"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 830,
    "Timestamp": "2024-02-22T17:13:24.603",
    "ValueAsString": "[2024-02-22T17:13:24+0000] OTA             : Accepted"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 908,
    "Timestamp": "2024-02-22T17:14:48.277",
    "ValueAsString": "2.4.0.0"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 909,
    "Timestamp": "2023-12-16T13:40:52.8",
    "ValueAsString": "6"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 911,
    "Timestamp": "2024-02-22T17:14:48.277",
    "ValueAsString": "2.3.1.2"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 916,
    "Timestamp": "2024-02-22T17:14:48.277",
    "ValueAsString": "2.3.1.2"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 917,
    "Timestamp": "2023-05-23T16:09:31.697",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 918,
    "Timestamp": "2023-05-23T16:09:31.7",
    "ValueAsString": "1"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 960,
    "Timestamp": "2023-05-22T19:11:57.81",
    "ValueAsString": "242016001464605"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 962,
    "Timestamp": "2023-05-22T19:11:57.813",
    "ValueAsString": "89470060210810207713"
  },
  {
    "ChargerId": "5273dba1-4e80-49ae-ac44-9420ae596ad3",
    "StateId": 963,
    "Timestamp": "2023-05-22T19:11:57.813",
    "ValueAsString": "866642058845159"
  }
]
````

And all these StateId's have to be interpreted using https://api.zaptec.com/api/constants

## Manipulate the charger
````bash
curl -vsk -X POST -H "Accept: */*" -H "Authorization: bearer $TOKEN" https://api.zaptec.com/api/chargers/<charger-guid>/sendCommand/<cmdCode> 
````
Where cmdCodes are (https://api.zaptec.com/help/index.html#/Charger/post_api_chargers__id__sendCommand__commandId_):
102   - restart charger
200   - upgrade firmware
506   - stop/pause chargeer
507   - resume charging
10001 - deauthorize and stop charging

## api constants
````bash
curl -sk -H "Accept: */*" -H "Authorization: bearer $TOKEN" https://api.zaptec.com/api/constants | jq
````
