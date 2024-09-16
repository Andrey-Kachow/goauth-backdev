# goauth-backdev
Часть сервиса аутентификации. Тестовое задание на позицию Junior Backend Developer

## Необходимо установить дла работы:

- [Docker](https://www.docker.com/get-started).
- [Docker Compose](https://docs.docker.com/compose/install/).

## Инструкции по запуску:
Приложение запускается в двух контейнерах при помощи `docker-compose`. Один контейнер отвечает за Postgres базу данных, другой содержит в себе само приложение

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Andrey-Kachow/goauth-backdev.git
cd goauth-backend
```

### 2. Установите значения переменных среды (по желанию)
`docker-compose` ищет файл `.env` и бере от туда переменные среды. Небходимо создать этот файл.
Создайте файл .env и скопируйте в него содержимое файла.env.template, который содержит все необходимые значения по умолчанию.
```
cp .env.template .env
```

Некоторые переменные имеют значение CHANGEME, которое по желанию можно поменять.
```
# Database configuration
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=myuser
POSTGRES_PASSWORD=CHANGEME
POSTGRES_DB=mydb

# Application Environment [development|production|debug]
GOAUTH_BACKDEV_MODE="development"

# Application Secrets: Notification Service configuration: 
GOAUTH_BACKDEV_SMTP_HOST=CHANGEME
GOAUTH_BACKDEV_EMAIL_USERNAME=CHANGEME
GOAUTH_BACKDEV_EMAIL_PASSWORD=CHANGEME

```
В частности, переменные,  `GOAUTH_BACKDEV_SMTP_HOST`, `GOAUTH_BACKDEV_EMAIL_USERNAME` и `GOAUTH_BACKDEV_EMAIL_PASSWORD` содержат данные SMTP сервера, которые можно поменять для отправки email.
Если эти переменные не указаны или оставлены в значении по умолчанию, тогда будет использован серфис-пустышка, который реализует тот-же интерфейс, но email не отправляет.
Их можно поменять на данные GMAIL или Mailtrap или иных серверов.

### Запуск
Запускается прилодение командой
```
docker-compose up --build
```

## Проверка работоспособности.
Помимо двух маршрутов `api/access` и `api/refresh`, приложение реализует маршрут `/`, в котором возвращает небольшую HTML страницу мини-клиент, в котором можно потестировать приложение.

