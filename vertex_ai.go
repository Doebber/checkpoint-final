package function

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"google.golang.org/genai"
)

const (
	gcpProjectID = "project-ab4fa986-e441-43f4-9fe"
	gcpLocation  = "us-central1"
	vertexModel  = "gemini-2.5-flash"
)

var (
	vertexClient     *genai.Client
	vertexClientErr  error
	vertexClientOnce sync.Once
)

type CreditAnalysisResult struct {
	Resumo              string   `json:"resumo"`
	PontosPositivos     []string `json:"pontos_positivos"`
	PontosAtencao       []string `json:"pontos_atencao"`
	AvaliacaoPreliminar string   `json:"avaliacao_preliminar"`
	Recomendacao        string   `json:"recomendacao"`
}

func getVertexClient(ctx context.Context) (*genai.Client, error) {

	vertexClientOnce.Do(func() {

		vertexClient, vertexClientErr = genai.NewClient(ctx, &genai.ClientConfig{
			Project:  gcpProjectID,
			Location: gcpLocation,
			Backend:  genai.BackendVertexAI,
		})
	})

	if vertexClientErr != nil {
		return nil, fmt.Errorf(
			"erro ao criar cliente Vertex AI: %w",
			vertexClientErr,
		)
	}

	return vertexClient, nil
}

func AnalyzeCreditWithVertexAI(
	ctx context.Context,
	input CreditAnalysisInput,
) (*CreditAnalysisResult, error) {

	client, err := getVertexClient(ctx)

	if err != nil {
		return nil, err
	}

	prompt := BuildCreditAnalysisPrompt(input)

	response, err := client.Models.GenerateContent(
		ctx,
		vertexModel,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			Temperature:      genai.Ptr[float32](0.2),
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao executar análise no Vertex AI: %w",
			err,
		)
	}

	responseText := response.Text()

	if responseText == "" {
		return nil, fmt.Errorf("Vertex AI retornou uma resposta vazia")
	}

	var result CreditAnalysisResult

	if err := json.Unmarshal(
		[]byte(responseText),
		&result,
	); err != nil {

		return nil, fmt.Errorf(
			"erro ao interpretar resposta do Vertex AI: %w",
			err,
		)
	}

	return &result, nil
}