# channels-golang 
Задание:
Написать утилиту, которая будет вычислять "производительность". Утилита параллельно рассылает http-запросы на какой-либо ресурс или несколько ресурсов (например, https://google.com).
Программа принимает следующие параметры на вход:
- адрес ресурса(ов)
- количество запросов, которые необходимо выполнить
- таймаут (время ожидания, после которого мы ответа уже не ждём)

Программа собирает и выводит на экран следующие данные:
- время, за которое отработали все запросы
- среднее время на запрос
- максимальное/минимальное время возвращение ответа
- количество ответов, которых не дождались

Пример работы программы:
```sh
go run main.go -address "https://www.google.com/" -requestsNumber "10" -timeoutMilliseconds "100"
```
![alt text](https://github.com/bukhavtsov/channels-golang/blob/task-3/img/screenshots/screenshot-1.png)

