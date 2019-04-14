# keyval
in-memory key-value хранилище

# Принцип работы
* Данные хранятся в `map[string]interface{}`.
* Race condition разрешается с помощью `sync.RWMutex`.
* Для увеличения производительности используется партицирование данных.
* Партиция ключа определяется хэш-функцией.
* Очистка устаревших ключей производится в отдельном потоке. Map блокируется и проверяется каждый ключ.

# Запуск
## Флаги
* `-p 8000` - Порт API.
* `-c 100` - Количество партиций (инстансов map).
* `-i 1000` - Интервал очистки устаревших ключей в ms.

# REST API
* Формат ответа - JSON. Может возвращаться ответ с пустым телом.
* Формат ответа со значением `{"value":"foo"}`.
* TTL указывается в ms.

## Коды ошибок
* 400 - Некорректный запрос
* 404 - Запись не найдена

## Методы
* GET `/api/keys` - Возвращает массив строк со всеми ключами.
* GET `/api/get/{key}` - Возвращает значение по ключу.
* POST `/api/set/{key}` - Сохраняет строковое значение по ключу. Формат запроса: `{"value":"foo", "ttl":1000}`.
* POST `/api/remove/{key}` - Удаляет значение по ключу.
* POST `/api/expire/{key}` - Устанавливает ttl ключа. Формат запроса: `{"ttl":1000}`.
* GET `/api/hget/{key}/{field}` - Возвращает значение поля словаря.
* POST `/api/hset/{key}/{field}` - Устанавливает значение поля словаря. Формат запроса: `{"value":"foo"}`.
* POST `/api/hdel/{key}/{field}` - Удаляет значение поля из словаря.

# Benchmarks
```
BenchmarkStorage_Set-12                  1000000              1534 ns/op             384 B/op          2 allocs/op
BenchmarkCluster_Set10-12                5000000               588 ns/op             107 B/op          4 allocs/op
BenchmarkCluster_Set100-12               5000000               550 ns/op              72 B/op          4 allocs/op
BenchmarkCluster_Set1000-12              5000000               515 ns/op              78 B/op          4 allocs/op
BenchmarkStorage_SetGet-12               1000000              2520 ns/op             422 B/op          4 allocs/op
BenchmarkCluster_SetGet10-12             3000000               663 ns/op             136 B/op          7 allocs/op
BenchmarkCluster_SetGet100-12            2000000               668 ns/op             134 B/op          7 allocs/op
BenchmarkCluster_SetGet1000-12           2000000               686 ns/op             154 B/op          7 allocs/op
```

# TODO
* Реализация хранения списков
* Сохранение данных на диск
* Тесты
* Документация
* Graceful shutdown
* Улучшение алгоритма партицирования
* Оптимизации и бэнчмарки