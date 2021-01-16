# Article Service
Medium Rare에서 아티클, 댓글, 대댓글 등을 관리하는 메인 도메인 서비스.

## Run
```
$ go run main.go
```

## Config
|환경변수|설명|기본값|
|:-|:-|:-|
|APP_ADDR|포트번호를 포함한 서비스 주소|0.0.0.0:8082|
|APP_JWTSECRET|JWT 토큰의 서명키|kkh|
|APP_DB_DRIVER|사용할 DB 드라이버|mysql|
|APP_DB_URL|DB 연결을 위한 URL|root:rootpw@tcp(127.0.0.1:3306)/articleDB?charset=utf8mb4&parseTime=True&loc=Local|
|APP_LOG_OUTPUTS|로그를 출력할 파일들|stdout|
|APP_LOG_LEVEL|로그 출력 레벨|debug|
|APP_LOG_ENCODING|로그 출력 형식|json|
|APP_KAFKA_BROKERS|카프카 브로커 주소들|127.0.0.1:9092|

