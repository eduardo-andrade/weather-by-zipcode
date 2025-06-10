# Projeto Weather by ZipCode

Este projeto é composto por dois serviços em Go que, trabalhando em conjunto, recebem um CEP e retornam informações climáticas da cidade correspondente. A arquitetura é distribuída, com uso de containers Docker e suporte a rastreamento distribuído via OpenTelemetry e Zipkin.

## Estrutura de Diretórios

```
.
├── docker-compose.yml
├── go.mod
├── go.sum
├── LICENSE
├── otel-config.yaml
├── service-a
│   ├── Dockerfile
│   ├── handler.go
│   └── main.go
├── service-b
│   ├── Dockerfile
│   ├── handler.go
│   ├── main.go
│   ├── utils.go
│   ├── viacep
│   └── weather
└── shared
    └── tracer
```

## Requisitos

- Go 1.20 ou superior
- Docker e Docker Compose
- Conexão com a internet para acessar as APIs externas:
  - [ViaCEP](https://viacep.com.br/)
  - [WeatherAPI](https://www.weatherapi.com/)
- Chave de API da WeatherAPI (setada na variável de ambiente `WEATHER_API_KEY`)

## Instalação

1. **Clone o repositório**:
   ```bash
   git clone <url-do-repositorio>
   cd weather-by-zipcode
   ```

2. **Configure sua chave da WeatherAPI** (exemplo usando `.env` ou diretamente no `docker-compose.yml`)

3. **Construa os containers**:
   ```bash
   docker-compose build
   ```

4. **Suba o ambiente**:
   ```bash
   docker-compose up
   ```

## Execução

1. Faça uma requisição HTTP POST para o `service-a` com um CEP no corpo:
   ```bash
   curl -X POST http://localhost:8080 -d '{"cep": "01001-000"}' -H "Content-Type: application/json"
   ```

2. A resposta será semelhante a:
   ```json
   {
     "city": "São Paulo",
     "temperature": {
       "celsius": 22.0,
       "fahrenheit": 71.6,
       "kelvin": 295.15
     }
   }
   ```

## Observabilidade

- O projeto inclui suporte a tracing distribuído usando OpenTelemetry e Zipkin.
- Acesse o painel do Zipkin:
  - [http://localhost:9411](http://localhost:9411)
- Verifique os traces ao fazer requisições para `service-a`.

## Problemas Comuns

- **Zipkin não exibe traces**:
  - Certifique-se de que a URL do exportador esteja usando `http://zipkin:9411` dentro dos containers.
  - Verifique a configuração em `shared/tracer`.

- **Erro de rede ao chamar APIs externas**:
  - Verifique se o container possui acesso à internet.

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
