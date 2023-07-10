# go-yadro-testtask
Тестовое задание от компании Yadro на обработку событий из входного файла

## Сборка и запуск
### Обычный способ
Для сборки проекта и запуска проекта следуйте инструкции:
  1. Загрузить репозиторий на свой компьютер. Для этого вам нужно выполнить команду:
   ```bash
   git clone https://github.com/SpeedCrash100/go-yadro-testtask
   ```

  2.  Собрать проект:
  ```bash
  go build -o program ./cmd
  ```
  3. Вы получите программу `program` запустите её указав входной файл:
  ```bash
  ./program <file_name>
  ```
### Сборка и запуск в Docker. Srly?
  1. После загрузки репозитория собрать контейнер с приложением можно используя следующую команду:
  ```bash
  docker build . --tag go-yadro-task:1
  ```
  2. Для внесения входных файлов в контейнер можно воспользоваться опцией `--mount` в `docker run` Указав в ключе `source` место с входными файлами и в `target` место где они будут располагаться в контейнере.
  3. Запуск программы, где тестовые примеры из [test_cases/input](https://github.com/SpeedCrash100/go-yadro-testtask/tree/85282a30f699be2fca6ed66f85740e189e994341/test_cases/input) примонтированы к папке `/app/data` происходит следующей командой:
  ```bash
  docker run --rm -it --mount "type=bind,source=./test_cases/input,target=/app/data" go-yadro-task:1 /app/program /app/data/stock.txt
  ```
  Последний аргумент задает входной файл. В данном случае, используется файл stock.txt, в котором представлен пример из самого задания.



## Тестирование
Для программы написаны базовые тест, которые можно запустить командной: 
```bash
go test ./pkg
```
### Тест TestApp
Данный тест может автоматизированно использовать входные данные из папки `test_cases/input` для запуска программы и затем проверять соответствие вывода `test_cases/output`. Для этого файлы в этой директории должны иметь одинаковые имена и расширения.

Что бы запустить только этот тест используйте:
```bash
go test ./pkg --run TestApp
```
Добавьте флаг `-v`, что бы удостоверится, что все тесты выполняются.

Пример вывода с флагом `-v`:
```
=== RUN   TestApp
=== RUN   TestApp/client_already_in_club.txt
    app_test.go:98: app process error: <nil>
=== RUN   TestApp/client_leaves_in_valid_order.txt
    app_test.go:98: app process error: <nil>
=== RUN   TestApp/client_unknown_symbols.txt
    app_test.go:98: app process error: invalid event format
=== RUN   TestApp/error_event_invalid_order.txt
    app_test.go:98: app process error: invalid order of events
=== RUN   TestApp/events_after_close.txt
    app_test.go:98: app process error: <nil>
=== RUN   TestApp/invalid_time_leading_zeros.txt
    app_test.go:98: app process error: invalid time in input files
=== RUN   TestApp/invalid_time_out_of_range.txt
    app_test.go:98: app process error: time format valid but values are out of range
=== RUN   TestApp/no_events.txt
    app_test.go:98: app process error: <nil>
=== RUN   TestApp/stock.txt
    app_test.go:98: app process error: <nil>
=== RUN   TestApp/table_out_of_range.txt
    app_test.go:98: app process error: invalid event format
=== RUN   TestApp/unknown_client.txt
    app_test.go:98: app process error: <nil>
--- PASS: TestApp (0.00s)
    --- PASS: TestApp/client_already_in_club.txt (0.00s)
    --- PASS: TestApp/client_leaves_in_valid_order.txt (0.00s)
    --- PASS: TestApp/client_unknown_symbols.txt (0.00s)
    --- PASS: TestApp/error_event_invalid_order.txt (0.00s)
    --- PASS: TestApp/events_after_close.txt (0.00s)
    --- PASS: TestApp/invalid_time_leading_zeros.txt (0.00s)
    --- PASS: TestApp/invalid_time_out_of_range.txt (0.00s)
    --- PASS: TestApp/no_events.txt (0.00s)
    --- PASS: TestApp/stock.txt (0.00s)
    --- PASS: TestApp/table_out_of_range.txt (0.00s)
    --- PASS: TestApp/unknown_client.txt (0.00s)
PASS
ok      github.com/speedcrash100/go-yadro-testtask/pkg  0.002s
```

## Как работает приложение
Приложение имеет состояние [State](https://github.com/SpeedCrash100/go-yadro-testtask/blob/main/pkg/state.go), которое может изменятся и дополняться согласно входным событиям реализующие [InputEvent](https://github.com/SpeedCrash100/go-yadro-testtask/blob/02f08ddc37cbb14c3e9a26a30bd99088c6ab2dcc/pkg/event.go#L104)
