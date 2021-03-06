# Телеграм Бот для Bitbucket 

Данный бот отслеживает активность в организации битбакета и отправляет сообщения о пул реквестах, комментариях, апрувах и мержах в телеграм чат проекта 

[![Build Status](https://travis-ci.org/anatoliyfedorenko/bitbucketbot.svg?branch=master)](https://travis-ci.org/anatoliyfedorenko/bitbucketbot) [![Go Report Card](https://goreportcard.com/badge/github.com/anatoliyfedorenko/bitbucketbot)](https://goreportcard.com/report/github.com/anatoliyfedorenko/bitbucketbot) [![Coverage Status](https://coveralls.io/repos/github/anatoliyfedorenko/bitbucketbot/badge.svg)](https://coveralls.io/github/anatoliyfedorenko/bitbucketbot)

## Инструкция по установке:

### Шаг 1 Создать бота и добавить его в чат проекта

1. Откройте телеграм
2. Найдите аккаунт @BotFather
3. Следуйте инструкциям BotFather чтобы создать нового бота
4. Скопируйте бот токен 
5. Зайдите в настройки бота (BotSettings) и отключите Privacy mode во вкладке Group Privacy
6. Добавьте бота в чат проекта. 
7. Скопируйте ID чата. Его можно узнать в URL (последние цифры). 
Например #/im?p=g134321707 ID чата "-134321707"

### Шаг 2 Загрузить бота на сервер и добавить переменные окружения

Заргузите код бота на ваш рабочий сервер. 

Для того, чтобы бот работал, нужно создать переменные окружения, в которых необходимо указать свой телеграм токен, полученный от BotFather, а так же ID телеграм чата, куда в хотите получать сообщения

В папке проекта создайте файл ```.env```, откройте его и добавьте туда код ниже. Замените значения в <> на реальные: 
```
BOT_TELEGRAM_TOKEN=<your-telegram-bot-token>
BOT_CHAT=<your-telegram-chat>
```
Затем в терминале выполните следующие команды: 
```
set -a 
. .env
set +a
```

### Шаг 3 Настроить вебхуки (webhooks) в Битбакете
1. Зайдите в Битбакет и создайте репозиторий
2. Зайдите в настройки репозитория (Settings) во вкладку Webhooks
3. Создайте вебхуки. При создании отмечайте только соотвествующие действия во вкладке "Triggers", которые вы хотите отследить. Например для вебхука "PR Created" отметьте только "Pull Request > Created" и т.д.

Таблица вебхуков: 

| Title | URL | 
| --- | --- |
| PR Created | <адресСервераСБотом>:порт/pull_request_created | 
| PR Commented | <адресСервераСБотом>:порт/pull_request_commented |
| PR Approved | <адресСервераСБотом>:порт/pull_request_approved | 
| PR Merged | <адресСервераСБотом>:порт/pull_request_merged | 
| PR Declined | <адресСервераСБотом>:порт/pull_request_declined | 

### Шаг 4 Запустить бота 

Выполните команду ```make build``` 
Затем ```make build_docker```

И наконец запустите контенер командой ```docker run -d -p порт:порт bitbucketbot```
Если докер не видит ваши переменные окружения, передайте их внутрь контенера флагами ```-e```

Попробуйте отправить ПР в ваш репозиторий и получить оповещение об этом от бота! 