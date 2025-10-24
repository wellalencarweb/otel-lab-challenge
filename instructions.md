# 🧩 Instruções para o Copilot

## 🎯 Objetivo

Desenvolver um **sistema em Go** composto por **dois serviços (A e B)** que, a partir de um **CEP**, identifiquem a cidade correspondente e retornem o **clima atual** (em **Celsius**, **Fahrenheit** e **Kelvin**), implementando **OpenTelemetry (OTEL)** e **Zipkin** para **tracing distribuído**.

---

## 🧱 Estrutura do Projeto

- **Serviço A:** responsável por receber o input e encaminhar ao Serviço B.  
- **Serviço B:** responsável por orquestrar as chamadas externas (ViaCEP e WeatherAPI) e retornar o resultado formatado.

---

## 🚀 Requisitos — Serviço A (Input Handler)

### Endpoint
`POST /cep`

### Request Body
```json
{
  "cep": "29902555"
}
```

### Regras de Negócio
1. Validar se o input:
   - É uma **string** de **8 dígitos** numéricos.  
2. Se for **válido**:
   - Encaminhar a requisição ao **Serviço B** via HTTP.  
3. Se for **inválido**:
   - Retornar:
     ```json
     {
       "message": "invalid zipcode"
     }
     ```
     **HTTP Status:** `422`

---

## 🔁 Requisitos — Serviço B (Orquestrador)

### Responsabilidades
1. Receber um **CEP válido (8 dígitos)** do Serviço A.
2. Consultar a **API ViaCEP** (ou similar) para obter o **nome da cidade**.
3. Consultar a **WeatherAPI** (ou similar) para obter a **temperatura em Celsius**.
4. Converter e retornar as temperaturas:
   - **Fahrenheit:** `F = C * 1.8 + 32`
   - **Kelvin:** `K = C + 273`
5. Retornar a resposta formatada.

---

### Respostas Esperadas

#### ✅ Sucesso
**HTTP Status:** `200`  
**Response Body:**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### ⚠️ Erro — CEP inválido (formato incorreto)
**HTTP Status:** `422`  
**Response Body:**
```json
{
  "message": "invalid zipcode"
}
```

#### ❌ Erro — CEP não encontrado
**HTTP Status:** `404`  
**Response Body:**
```json
{
  "message": "can not find zipcode"
}
```

---

## 🛰️ OTEL + Zipkin (Tracing Distribuído)

### Requisitos de Observabilidade
1. Implementar **OpenTelemetry** em ambos os serviços.
2. Configurar **tracing distribuído** entre **Serviço A** e **Serviço B**.
3. Utilizar **spans** para medir:
   - Tempo de resposta da consulta de CEP (ViaCEP)
   - Tempo de resposta da consulta de temperatura (WeatherAPI)
4. Enviar os traces para um **collector OTEL** configurado via **docker-compose**.
5. Integrar com **Zipkin** para visualização.

---

## 🔧 Dicas Técnicas

- **API de CEP:** [ViaCEP](https://viacep.com.br/)
- **API de Clima:** [WeatherAPI](https://www.weatherapi.com/)
- Fórmulas de conversão:
  - `F = C * 1.8 + 32`
  - `K = C + 273`

---

## 🐳 Entrega (Dev + Docker)

### O que deve ser entregue:
1. **Código-fonte completo** dos Serviços A e B.
2. **Documentação** explicando como rodar o projeto localmente.
3. **Dockerfile** e **docker-compose.yml** configurados para subir:
   - Serviço A
   - Serviço B
   - OTEL Collector
   - Zipkin

### Exemplo de Comando
```bash
docker-compose up --build
```

Após subir o ambiente:
- Teste com `curl` ou Postman:
  ```bash
  curl -X POST http://localhost:<porta-servico-A>/cep -d '{"cep":"01001000"}' -H "Content-Type: application/json"
  ```

---

## 🧭 Objetivo Final

- Serviço A e Serviço B se comunicando corretamente via HTTP.  
- Observabilidade implementada com OTEL + Zipkin.  
- Sistema retornando corretamente a cidade e as temperaturas em três unidades.  
- Ambiente completo rodando via Docker Compose.