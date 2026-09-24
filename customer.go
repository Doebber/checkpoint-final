package function

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/customers.json
var customersData []byte

type Cliente struct {
	ClienteID		 string `json:"cliente_id"`
	CPF              string  `json:"cpf"`
	Nome             string  `json:"nome"`
	EmprestimosAtivos int     `json:"emprestimos_ativos"`
	SaldoDevedor      float64 `json:"saldo_devedor"`
	ParcelasTotais    int     `json:"parcelas_totais"`
	ParcelasPagas     int     `json:"parcelas_pagas"`
	ParcelasEmAtraso  int     `json:"parcelas_em_atraso"`
	MaiorAtrasoDias  int     `json:"maior_atraso_dias"`
}

type BaseClientes struct {
	Clientes []Cliente `json:"clientes"`
}

// Interface que permite trocar JSON por banco de dados no futuro.
type CustomerRepository interface {
	GetByCPF(ctx context.Context, cpf string) (*Cliente, error)
}

// Implementação atual utilizando customers.json.
type JSONCustomerRepository struct{}

func (r JSONCustomerRepository) GetByCPF(ctx context.Context, cpf string) (*Cliente, error) {

	var base BaseClientes

	if err := json.Unmarshal(customersData, &base); err != nil {
		return nil, fmt.Errorf("erro ao carregar base de clientes: %w", err)
	}

	for i := range base.Clientes {
		if base.Clientes[i].CPF == cpf {
			return &base.Clientes[i], nil
		}
	}

	return nil, fmt.Errorf("cliente não encontrado")
}