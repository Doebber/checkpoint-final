# Projeto Final

## Visão geral

Este projeto apresenta uma arquitetura **event-driven** para processamento de solicitações de empréstimo pessoal, integrada com **Inteligência Artificial utilizando Vertex AI / Gemini**.

A solicitação de empréstimo é publicada em um tópico do **Google Cloud Pub/Sub**. A mensagem é processada de forma assíncrona por uma **Cloud Function de segunda geração**, acionada por **Eventarc**.

Durante o processamento, a função:

1. Recebe a solicitação de empréstimo.
2. Valida os dados recebidos.
3. Consulta o histórico do cliente em uma base interna simulada.
4. Prepara os dados para análise.
5. Envia o contexto para o Vertex AI / Gemini.
6. Obtém uma análise estruturada.
7. Registra o resultado no Google Cloud Logging.

A IA atua como **assistente de análise de crédito**. Ela não toma a decisão definitiva de aprovação ou reprovação. A análise e a recomendação geradas servem como apoio para um profissional responsável pela decisão final.

O projeto também utiliza **GitHub Actions** para automatizar o processo de CI/CD.

A autenticação entre GitHub Actions e Google Cloud utiliza **Workload Identity Federation (OIDC)**, sem armazenamento de chaves JSON ou credenciais permanentes no repositório.

---

# Arquitetura

A arquitetura principal do projeto é:

```text
                         GitHub
                           │
                           │ Push
                           ▼
                    GitHub Actions
                           │
                           │ OIDC / WIF
                           ▼
                     Google Cloud
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
      Cloud Function Gen2           Cloud Workflows
       LoanHandlerPubSub            loan-workflow
             │                           │
             │                           │
             │                           ▼
             │                       Pub/Sub
             │                        orders
             │                           │
             │                           ▼
             │                        Eventarc
             │                           │
             └──────────────┬────────────┘
                            ▼
                       LoanHandler
                            │
                  ┌─────────┴─────────┐
                  │                   │
                  ▼                   ▼
           customers.json        Vertex AI
           base interna           Gemini
                  │                   │
                  └─────────┬─────────┘
                            ▼
                    Análise de Crédito
                            │
                            ▼
                      Cloud Logging
                            │
                            ▼
                     Analista Humano
```

A comunicação entre os componentes de processamento é baseada em eventos, mantendo o desacoplamento entre a publicação da solicitação e o processamento da análise.

---

# Componentes utilizados

## Google Cloud

* Google Cloud Pub/Sub
* Google Cloud Eventarc
* Google Cloud Functions Gen2
* Google Cloud Workflows
* Vertex AI / Gemini
* Google Cloud Logging
* Google Cloud IAM
* Workload Identity Federation

## CI/CD

* GitHub
* GitHub Actions
* OpenID Connect (OIDC)
* Workload Identity Federation (WIF)

## Linguagem

* Go

## Modelo de IA

* Gemini 2.5 Flash através do Vertex AI

---

# Estrutura do projeto

```text
checkpoint-final/
│
├── .github/
│   └── workflows/
│       └── deploy.yml
│
├── data/
│   └── customers.json
│
├── main.go
├── function.go
├── customer.go
├── credit_analysis.go
├── vertex_ai.go
│
├── workflow.yaml
├── go.mod
└── go.sum
```

### `main.go`

Responsável pelo processamento principal do evento Pub/Sub.

Realiza:

* Recepção do CloudEvent.
* Decodificação da mensagem.
* Validação da solicitação.
* Consulta do cliente.
* Preparação da análise.
* Chamada do Vertex AI.
* Registro estruturado dos resultados.
* Tratamento de erros.

### `function.go`

Realiza o registro da função no Functions Framework:

```go
functions.CloudEvent("LoanHandler", LoanHandler)
```

### `customer.go`

Responsável pelo acesso aos dados dos clientes.

A implementação atual utiliza `customers.json` como uma base interna simulada.

Foi criada uma abstração `CustomerRepository` para permitir que futuramente a fonte de dados seja substituída por uma solução como Cloud SQL ou Firestore sem alterar a lógica principal da função.

### `credit_analysis.go`

Responsável pela preparação do contexto enviado ao modelo de IA.

Também define as informações utilizadas para orientar a análise de crédito.

### `vertex_ai.go`

Responsável pela integração com o Vertex AI / Gemini.

O arquivo:

* Inicializa o cliente Vertex AI.
* Envia o contexto da análise.
* Solicita uma resposta estruturada em JSON.
* Valida a resposta recebida.
* Converte o resultado para a estrutura utilizada pela aplicação.

### `data/customers.json`

Contém uma base simulada de clientes utilizada para o projeto acadêmico.

Os dados representam:

* CPF
* Nome
* Empréstimos ativos
* Saldo devedor
* Parcelas totais
* Parcelas pagas
* Parcelas em atraso
* Maior atraso em dias

Os dados são fictícios e não representam clientes reais.

---

# Fluxo de processamento

Uma solicitação de empréstimo possui a seguinte estrutura:

```json
{
  "name": "Joao",
  "cpf": "11111111111",
  "amount": 15000,
  "term": 24
}
```

O Workflow publica a solicitação no tópico:

```text
orders
```

A mensagem é processada pelo seguinte fluxo:

```text
Workflow
    ↓
Pub/Sub
    ↓
Eventarc
    ↓
LoanHandlerPubSub
    ↓
Validação da solicitação
    ↓
Busca do cliente
    ↓
Preparação da análise
    ↓
Vertex AI / Gemini
    ↓
Análise estruturada
    ↓
Cloud Logging
    ↓
Analista Humano
```

---

# Análise com Inteligência Artificial

A IA recebe exclusivamente os dados internos disponibilizados pela aplicação para realizar a análise.

O modelo é orientado a produzir:

* Resumo do histórico.
* Pontos positivos.
* Pontos de atenção.
* Avaliação preliminar de risco.
* Recomendação para o analista.

A resposta possui formato estruturado:

```json
{
  "resumo": "...",
  "pontos_positivos": [
    "..."
  ],
  "pontos_atencao": [
    "..."
  ],
  "avaliacao_preliminar": "...",
  "recomendacao": "..."
}
```

A aplicação utiliza um schema de resposta para orientar o modelo e também valida se a resposta retornada contém os campos esperados.

Uma resposta incompleta é tratada como erro.

## Decisão humana

A IA não realiza automaticamente a aprovação ou reprovação do empréstimo.

O fluxo é:

```text
Dados do cliente
       ↓
      IA
       ↓
Análise preliminar
       ↓
 Recomendação
       ↓
Analista humano
       ↓
 Decisão final
```

Dessa forma, a IA funciona como uma ferramenta de apoio à análise, mantendo a decisão final sob responsabilidade humana.

---

# Exemplo de análise

Durante um teste com um dos clientes da base simulada, foi utilizada uma solicitação de empréstimo contendo:

```json
{
  "name": "Joao",
  "cpf": "11111111111",
  "amount": 15000,
  "term": 24
}
```

A IA identificou que o cliente possuía histórico de pagamentos sem atrasos, mas também identificou uma inconsistência nos dados simulados:

* Empréstimo informado como ativo.
* Existência de saldo devedor.
* Todas as parcelas informadas como pagas.
* Nenhum atraso registrado.

A recomendação gerada orientou o analista humano a verificar essa inconsistência antes de prosseguir com a análise.

Esse cenário demonstra que o modelo não apenas resume os dados, mas também pode identificar informações que merecem validação.

---

# Como executar

## Pré-requisitos

É necessário possuir:

* Google Cloud CLI (`gcloud`)
* Go
* Uma conta com permissões no projeto Google Cloud
* Um projeto Google Cloud configurado
* APIs necessárias habilitadas

Clone o repositório:

```bash
git clone https://github.com/Doebber/checkpoint-final.git
```

Entre na pasta:

```bash
cd checkpoint-final
```

Instale e organize as dependências:

```bash
go mod tidy
```

---

# Configuração do Google Cloud

O projeto utiliza os seguintes recursos:

```text
Região:
us-central1

Tópico Pub/Sub:
orders

Cloud Function:
LoanHandlerPubSub

Entry point:
LoanHandler

Workflow:
loan-workflow
```

O projeto deve estar configurado no ambiente do Google Cloud CLI.

Para verificar o projeto atualmente selecionado:

```bash
gcloud config get-value project
```

---

# Habilitação do Vertex AI

A API do Vertex AI deve estar habilitada no projeto:

```bash
gcloud services enable aiplatform.googleapis.com
```

A Service Account utilizada pela execução da Cloud Function deve possuir permissão para utilizar o Vertex AI.

A permissão necessária é:

```text
roles/aiplatform.user
```

---

# Deploy e CI/CD

O projeto utiliza **GitHub Actions** como mecanismo principal de CI/CD.

O workflow está localizado em:

```text
.github/workflows/deploy.yml
```

O pipeline é acionado automaticamente quando ocorre um `push` na branch:

```text
main
```

O processo realiza:

1. Checkout do código.
2. Autenticação no Google Cloud.
3. Autenticação utilizando OIDC.
4. Workload Identity Federation.
5. Configuração do Google Cloud CLI.
6. Deploy da Cloud Function.
7. Deploy do Google Cloud Workflow.

Dessa forma, tanto a aplicação quanto o Workflow são atualizados automaticamente a partir do repositório.

---

# Deploy da Cloud Function

O deploy da Cloud Function é realizado automaticamente pelo GitHub Actions.

O comando utilizado pelo pipeline é equivalente a:

```bash
gcloud functions deploy LoanHandlerPubSub \
  --gen2 \
  --runtime=go124 \
  --region=us-central1 \
  --entry-point=LoanHandler \
  --trigger-topic=orders
```

O deploy manual pode ser utilizado para testes ou manutenção, mas não é o fluxo principal do projeto.

---

# Deploy do Workflow

O arquivo do Workflow está localizado na raiz do repositório:

```text
workflow.yaml
```

O GitHub Actions realiza o deploy utilizando:

```bash
gcloud workflows deploy loan-workflow \
  --location=us-central1 \
  --source=workflow.yaml
```

Dessa forma, uma alteração realizada no `workflow.yaml` e enviada para a branch `main` também pode ser publicada automaticamente no Google Cloud.

---

# Testando o fluxo

Uma solicitação pode ser enviada através do Workflow:

```bash
gcloud workflows run loan-workflow \
  --location=us-central1 \
  --data='{
    "name":"Joao",
    "cpf":"11111111111",
    "amount":15000,
    "term":24
  }'
```

Uma execução bem-sucedida apresenta:

```text
state: SUCCEEDED
```

O Workflow também retorna o identificador da mensagem publicada no Pub/Sub.

Exemplo:

```json
{
  "messageId": "...",
  "status": "SUCCESS"
}
```

---

# Verificando os logs

Os logs da Cloud Function podem ser consultados utilizando:

```bash
gcloud functions logs read LoanHandlerPubSub \
  --gen2 \
  --region=us-central1 \
  --limit=30
```

Durante o processamento são registrados eventos estruturados, como:

```text
loan_request_received
```

```text
loan_credit_analysis_completed
```

```text
loan_request_processed
```

Em caso de erro:

```text
loan_request_error
```

O CPF não é registrado nos logs, reduzindo a exposição desnecessária de dados pessoais.

---

# Observabilidade

O projeto utiliza **structured logging** para facilitar a identificação e consulta dos eventos.

Os registros incluem informações como:

* ID do evento.
* Nome do cliente.
* Valor solicitado.
* Prazo.
* Tempo de processamento.
* Resultado da análise.

O tempo de processamento também é registrado através do campo:

```text
duration_ms
```

Os eventos estruturados permitem acompanhar as principais etapas da execução e identificar erros durante o processamento.

---

# CI/CD

O projeto utiliza **GitHub Actions** para automatizar o processo de implantação.

O workflow está localizado em:

```text
.github/workflows/deploy.yml
```

O pipeline é executado automaticamente quando ocorre um `push` na branch:

```text
main
```

O processo executa:

1. Checkout do código.
2. Autenticação no Google Cloud.
3. Autenticação utilizando OIDC.
4. Workload Identity Federation.
5. Configuração do Google Cloud CLI.
6. Deploy da Cloud Function.
7. Deploy do Workflow.

Fluxo simplificado:

```text
GitHub Push
     ↓
GitHub Actions
     ↓
OIDC
     ↓
Workload Identity Federation
     ↓
Google Cloud
     ├── Cloud Function
     └── Workflow
```

---

# Workload Identity Federation

A autenticação do GitHub Actions não utiliza chaves JSON.

O fluxo utilizado é:

```text
GitHub
   ↓
OIDC
   ↓
Workload Identity Pool
   ↓
OIDC Provider
   ↓
Service Account
   ↓
Google Cloud
```

Essa configuração permite que o GitHub Actions obtenha credenciais temporárias para realizar os deployments sem armazenar uma chave privada no GitHub.

---

# Segurança

O projeto não utiliza:

* Chaves JSON.
* Tokens permanentes.
* Senhas armazenadas no repositório.
* Credenciais privadas versionadas.
* CPF nos logs da aplicação.

A autenticação do CI/CD utiliza:

* OpenID Connect (OIDC).
* Workload Identity Federation.
* Service Account com permissões específicas.

Além disso, a IA recebe somente os dados necessários para realizar a análise e não recebe o CPF como parte do contexto enviado ao modelo.

A base `customers.json` contém somente dados fictícios para fins acadêmicos.

---

# Decisões arquiteturais

## Arquitetura event-driven

Foi mantida uma arquitetura baseada em eventos utilizando Pub/Sub e Eventarc.

Isso permite que a publicação da solicitação seja desacoplada do processamento da análise.

```text
Publicação
    ↓
Pub/Sub
    ↓
Eventarc
    ↓
Processamento assíncrono
```

## Uso do Google Cloud Workflows

O Workflow foi utilizado como ponto de entrada para os testes e para a publicação das solicitações no Pub/Sub.

Isso permite iniciar o fluxo de forma controlada sem a necessidade de criar uma API HTTP adicional.

```text
Workflows
    ↓
Pub/Sub
    ↓
Eventarc
    ↓
Cloud Function
```

## Uma única Cloud Function para o processamento

A lógica foi dividida em diferentes arquivos Go, mas todos fazem parte da mesma Cloud Function:

```text
main.go
customer.go
credit_analysis.go
vertex_ai.go
```

Essa organização permite separar responsabilidades sem criar várias funções para cada etapa do processamento.

## Abstração do acesso aos clientes

Foi criada uma interface de repositório para que o `customers.json` possa futuramente ser substituído por uma base de dados real sem alterar a lógica principal da aplicação.

Possíveis implementações futuras:

```text
customers.json
      ↓
CustomerRepository
      ↓
Cloud SQL / Firestore / outro banco
```

## IA como apoio à decisão

A IA foi posicionada como uma camada de apoio ao analista humano.

Ela identifica informações relevantes, inconsistências e pontos de atenção, mas não toma a decisão definitiva sobre o crédito.

---

# Tratamento de erros e retry

O processamento diferencia erros de negócio de erros técnicos.

Erros relacionados a dados inválidos ou cliente não encontrado podem ser tratados pela aplicação sem provocar novas tentativas desnecessárias.

Falhas na comunicação com serviços externos, como o Vertex AI, são retornadas pela função para permitir o mecanismo de retry da arquitetura orientada a eventos.

O Workflow também possui mecanismo de retry para a publicação no Pub/Sub.

Essa abordagem evita que falhas temporárias sejam tratadas como erros permanentes.

---

# Conclusão

O projeto consolida os principais conceitos trabalhados durante os checkpoints:

* Arquitetura serverless.
* Processamento orientado a eventos.
* Google Cloud Pub/Sub.
* Eventarc.
* Cloud Functions Gen2.
* Workflows.
* Retry e tratamento de erros.
* Structured Logging.
* Observabilidade.
* Integração com Vertex AI / Gemini.
* CI/CD.
* GitHub Actions.
* OpenID Connect.
* Workload Identity Federation.
* Separação de responsabilidades.
* Apoio à decisão utilizando IA.

A solução demonstra como uma solicitação de empréstimo pode ser processada de forma assíncrona utilizando serviços gerenciados do Google Cloud e Inteligência Artificial para gerar uma análise estruturada.

A arquitetura mantém o processamento desacoplado, utiliza autenticação sem chaves permanentes no CI/CD e mantém a decisão final de crédito sob responsabilidade de um analista humano.