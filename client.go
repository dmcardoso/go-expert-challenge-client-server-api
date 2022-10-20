package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"createDate"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao?currency=USD-BRL", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}

	var quote Quote
	err = json.Unmarshal(body, &quote)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)
	}

	file, err := os.Create("cotacao.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %v\n", err)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %v\n", quote.Bid))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever em arquivo: %v\n", err)
	}

	fmt.Println("Arquivo criado com sucesso")

	fileContent, err := os.ReadFile("cotacao.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler conteúdo do arquivo: %v\n", err)
	}

	fmt.Println(string(fileContent))
}
