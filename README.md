# Описание проекта
Простейший планировщик задач с функциями добавления, удаления и редактирования параметров задачи.

# Список выполенных заданий со звёздочкой
1. Реализована возможность определять извне порт при запуске сервера. Если существует переменная окружения TODO_PORT, сервер при старте должен слушать порт со значением этой переменной. 
2. Реализована возможность определять путь к файлу базы данных через переменную окружения. Для этого сервер должен получать значение переменной окружения TODO_DBFILE и использовать его в качестве пути к базе данных, если это не пустая строка.
3. Поддержка всех вариантов правил повторения.
4. В браузере рядом с кнопкой Добавить задачу есть поле для поиска. Добавлена возможность выбрать задачи через строку поиска (по заголовку или комментарию к задач или по дате).
5. Создан Dockerfile
Пример запуска:
docker build --tag go_final_project:latest .
docker run -it --env-file .env  -d go_final_project:latest

# Файл .env 
Заведены переменные окружения TODO_PORT, TODO_DBFILE, CGO_ENABLED, GOOS, GOARCH

# Запуск тестов 
В файле tests/settings.go следует указывать следующие параметры:
var DBFile = "../internal/database/scheduler.db"
var FullNextDate = true
var Search = true

Локально проект можно запускать через 
go build -o main cmd/api/main.go 
./main