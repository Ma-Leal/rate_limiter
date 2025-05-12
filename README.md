# Rate Limiter

Implementação de um rate limiter em Go com Redis.

## Para rodar o projeto

1. Subir o Redis com o Docker Compose:
   ```bash
   docker-compose up
   ```

2. Execute a aplicação Go:
   ```bash
   go run cmd/server/main.go
   ```

3. Serviço disponível na URL `http://localhost:8080`.

## Exemplos de Uso com `curl`

- **Rate Limiter por IP** (máximo de 5 requisições por IP):
  ```bash
  curl http://localhost:8080
  ```

- **Rate Limiter por Token** (máximo de 10 requisições por token):
  ```bash
  curl -H "API_KEY: abc123" http://localhost:8080
  ```

