### Simple Bitbucket Telegram Bot 

To configure bot, please insert env variables: 
```
BOT_TELEGRAM_TOKEN=<your-telegram-bot-token>
BOT_CHAT=<your-telegram-chat>
```

In Bitbucket settings add the following webhooks with appropriate triggers:

| Title | URL | 
| --- | --- |
| Pull Request Created | <yourIPAddr>/merge_created | 
| Pull Request Merged | <yourIPAddr>/merge_accepted | 

Launch server!