package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/rs/cors"
)

type Request struct {
	ValorInicial float64 `json:"valor_inicial"`
	TaxaJuros    float64 `json:"taxa_juros"`
	Periodo      int     `json:"periodo"`
}

type Response struct {
	ValorFinal float64 `json:"valor_final"`
}

func calculateCompoundInterest(valorInicial, taxaJuros float64, periodo int) float64 {
	// Fórmula para o cálculo dos juros compostos
	return valorInicial * math.Pow((1+taxaJuros), float64(periodo))
}

func interestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	var request Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	valorFinal := calculateCompoundInterest(request.ValorInicial, request.TaxaJuros, request.Periodo)

	response := Response{
		ValorFinal: valorFinal,
	}

	// Garantir que a resposta seja bem formatada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Erro ao codificar a resposta", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/calculate", interestHandler)

	// Configuração do CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                            // Permite qualquer origem (modifique conforme necessário)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // Métodos permitidos
		AllowedHeaders: []string{"Content-Type"},
	})

	// Inicializando o servidor com CORS habilitado
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
