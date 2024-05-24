## Execução

### Local

#### Testes Unitários e de Integração

Para executar todos os testes localizados no diretório `test`, utilize o comando:

```bash
go test ./test -v
```

#### Subir o Docker Compose

Para iniciar o serviço com Docker Compose, utilize o comando:

```bash
docker-compose up -d
```

#### Fazer Requisição

Após iniciar o serviço, faça uma requisição para verificar o funcionamento:

```bash
curl --request GET --url 'http://localhost:8080/08210010'
```

### Requisições no Cloud Run

#### CEP válido

Para consultar um CEP válido:

```bash
curl --request GET --url 'https://cloud-run-challenge-go-sdalcvq5xa-uc.a.run.app/08210010'
```

**Status:** 200

**Resposta:**

```json
{ 
  "temp_C": 19, 
  "temp_F": 66.2, 
  "temp_K": 292.15
}
```

#### CEP inválido

Para consultar um CEP inválido:

```bash
curl --request GET --url 'https://cloud-run-challenge-go-sdalcvq5xa-uc.a.run.app/0821001'
```

**Status:** 422

**Resposta:**

```json
{
  "error": "invalid zipcode"
}
```

#### CEP inexistente

Para consultar um CEP inexistente:

```bash
curl --request GET --url 'https://cloud-run-challenge-go-sdalcvq5xa-uc.a.run.app/99999999'
```

**Status:** 404

**Resposta:**

```json
{
  "error": "can not found zipcode"
}
```
