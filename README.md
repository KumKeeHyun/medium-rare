# Medium Rare
![image](https://user-images.githubusercontent.com/44857109/104845699-2170e280-591a-11eb-8150-c3db687cb3fa.png)

마이크로서비스 아키텍처로 만들어보는 예시 프로젝트

## Coupling
- user <-> reading-list (kafka event)
- article <-> reading-list (kafka event, RESTful API)
- article <-> trend (kafka event, RESTful API)