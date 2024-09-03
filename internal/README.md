# internal

Вычисление процента покрытия тестами:

```go
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
grep -v "mocks" profile.cov > cover.out
go tool cover -func cover.out
```
