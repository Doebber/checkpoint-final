package function

import (
	"fmt"
)

type CreditAnalysisInput struct {
	Nome                  string
	ValorSolicitado       float64
	PrazoSolicitado       int
	EmprestimosAtivos     int
	SaldoDevedor          float64
	ParcelasTotais        int
	ParcelasPagas         int
	ParcelasEmAtraso     int
	MaiorAtrasoDias       int
}

func BuildCreditAnalysisInput(
	request LoanRequest,
	customer *Cliente,
) CreditAnalysisInput {

	return CreditAnalysisInput{
		Nome:              customer.Nome,
		ValorSolicitado:   request.Amount,
		PrazoSolicitado:   request.Term,
		EmprestimosAtivos: customer.EmprestimosAtivos,
		SaldoDevedor:      customer.SaldoDevedor,
		ParcelasTotais:    customer.ParcelasTotais,
		ParcelasPagas:     customer.ParcelasPagas,
		ParcelasEmAtraso:  customer.ParcelasEmAtraso,
		MaiorAtrasoDias:   customer.MaiorAtrasoDias,
	}
}

func BuildCreditAnalysisPrompt(input CreditAnalysisInput) string {

	return fmt.Sprintf(`
Você é um assistente de análise de crédito.

Sua função é auxiliar um analista humano de crédito utilizando EXCLUSIVAMENTE
os dados internos fornecidos nesta solicitação.

Não utilize informações externas, pesquisas de mercado ou informações que não
estejam presentes nos dados fornecidos.

Não tome a decisão final de aprovação ou reprovação do crédito.
A decisão final deve ser sempre realizada por um analista humano.

Analise:

Cliente:
- Nome: %s
- Valor solicitado: R$ %.2f
- Prazo solicitado: %d meses

Histórico interno:
- Empréstimos ativos: %d
- Saldo devedor atual: R$ %.2f
- Total de parcelas: %d
- Parcelas pagas: %d
- Parcelas em atraso: %d
- Maior atraso registrado: %d dias

Com base exclusivamente nesses dados:

1. Faça um resumo objetivo do histórico do cliente.
2. Identifique pontos positivos.
3. Identifique pontos de atenção.
4. Apresente uma avaliação preliminar de risco.
5. Apresente uma recomendação para auxiliar o analista humano.

A recomendação NÃO deve ser apresentada como decisão definitiva.
`,
		input.Nome,
		input.ValorSolicitado,
		input.PrazoSolicitado,
		input.EmprestimosAtivos,
		input.SaldoDevedor,
		input.ParcelasTotais,
		input.ParcelasPagas,
		input.ParcelasEmAtraso,
		input.MaiorAtrasoDias,
	)
}