package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2/event"
)

type PubSubMessage struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type LoanRequest struct {
	Name   string  `json:"name"`
	CPF    string  `json:"cpf"`
	Amount float64 `json:"amount"`
	Term   int     `json:"term"`
}

func structuredLog(event string, data interface{}) {

	entry := map[string]interface{}{
		"event": event,
		"data":  data,
	}

	jsonData, err := json.Marshal(entry)

	if err != nil {

		log.Printf(
			`{"event":"logging_error","data":{"error":"failed to serialize log"}}`,
		)

		return
	}

	log.Println(string(jsonData))
}

func LoanHandler(
	ctx context.Context,
	event cloudevents.Event,
) error {

	startTime := time.Now()

	// ---------------------------------------------------------
	// 1. Recebimento do evento
	// ---------------------------------------------------------

	structuredLog(
		"loan_request_received",
		map[string]interface{}{
			"event_id": event.ID(),
			"source":   event.Source(),
			"type":     event.Type(),
		},
	)

	// ---------------------------------------------------------
	// 2. Decodificação do evento Pub/Sub
	// ---------------------------------------------------------

	var msg PubSubMessage

	if err := event.DataAs(&msg); err != nil {

		structuredLog(
			"loan_request_error",
			map[string]interface{}{
				"stage":       "event_data",
				"error":       err.Error(),
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		return err
	}

	data, err := base64.StdEncoding.DecodeString(
		msg.Message.Data,
	)

	if err != nil {

		structuredLog(
			"loan_request_error",
			map[string]interface{}{
				"stage":       "base64_decode",
				"error":       err.Error(),
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		return err
	}

	var req LoanRequest

	if err := json.Unmarshal(data, &req); err != nil {

		structuredLog(
			"loan_request_error",
			map[string]interface{}{
				"stage":       "json_decode",
				"error":       err.Error(),
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		return err
	}

	// ---------------------------------------------------------
	// 3. Validação da solicitação
	// ---------------------------------------------------------

	if req.Amount <= 0 {

		structuredLog(
			"loan_request_invalid",
			map[string]interface{}{
				"reason":      "invalid_amount",
				"amount":      req.Amount,
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		return nil
	}

	if req.Term <= 0 {

		structuredLog(
			"loan_request_invalid",
			map[string]interface{}{
				"reason":      "invalid_term",
				"term":        req.Term,
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		return nil
	}

	// ---------------------------------------------------------
	// 4. Busca do cliente
	// ---------------------------------------------------------

	repository := JSONCustomerRepository{}

	customer, err := repository.GetByCPF(
		ctx,
		req.CPF,
	)

	if err != nil {

		structuredLog(
			"loan_request_invalid",
			map[string]interface{}{
				"reason":      "customer_not_found",
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		// Cliente inexistente é um erro de negócio.
		// Não faz sentido ficar tentando novamente.
		return nil
	}

	// ---------------------------------------------------------
	// 5. Preparação da análise
	// ---------------------------------------------------------

	analysisInput := BuildCreditAnalysisInput(
		req,
		customer,
	)

	// ---------------------------------------------------------
	// 6. Análise utilizando Vertex AI
	// ---------------------------------------------------------

	analysisResult, err := AnalyzeCreditWithVertexAI(
		ctx,
		analysisInput,
	)

	if err != nil {

		structuredLog(
			"loan_request_error",
			map[string]interface{}{
				"stage":       "vertex_ai",
				"error":       err.Error(),
				"duration_ms": time.Since(startTime).Milliseconds(),
			},
		)

		// Erros da IA são retornados para permitir
		// retry da infraestrutura de eventos.
		return err
	}

	// ---------------------------------------------------------
	// 7. Resultado da análise
	// ---------------------------------------------------------

	structuredLog(
		"loan_credit_analysis_completed",
		map[string]interface{}{
			"event_id": event.ID(),
			"customer_id":     customer.ClienteID,
			"name":     customer.Nome,
			"amount":   req.Amount,
			"term":     req.Term,

			"analysis": analysisResult,

			"duration_ms": time.Since(startTime).Milliseconds(),
		},
	)

	// ---------------------------------------------------------
	// 8. Processamento concluído
	// ---------------------------------------------------------

	structuredLog(
		"loan_request_processed",
		map[string]interface{}{
			"event_id":    event.ID(),
			"customer_id":     customer.ClienteID,
			"name":        customer.Nome,
			"amount":      req.Amount,
			"term":        req.Term,
			"duration_ms": time.Since(startTime).Milliseconds(),
		},
	)

	return nil
}