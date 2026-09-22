# AdalQarau repeated-need screening

Учебный самостоятельный модуль на Go без внешних зависимостей. Он отмечает пары договоров одного заказчика с одинаковыми нормализованными предметом и единицей, точным количеством и суммой в заданном окне. Сигнал нужен для проверки и не доказывает нарушение.

## Запуск

```sh
go run ./cmd/screen screen -input testdata/contracts.csv -output /tmp/repeated.csv -window-days 30
go test ./...
go vet ./...
```

Вход содержит ровно поля `contract_id,customer_bin,contract_date,item_name,unit,quantity_units,amount_kzt`. Дата имеет формат ISO `YYYY-MM-DD`; количество и сумма являются неотрицательными целыми. Идентификаторы не повторяются. Выход стабилен и содержит `base_contract_id,repeated_contract_id,customer_bin,days_between,rule_version`.

Репозиторий создан для Assignment 3. Данные в `testdata` синтетические.
