package main

import (
	"context"
	"fmt"
	"log"
	"time"

	orderpb "github.com/Mellanie-Marques/microservices-proto/golang/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func criarPedido(client orderpb.OrderClient, idCliente int32, itens []*orderpb.OrderItem, precoTotal float32, nomeTeste string) {
	fmt.Printf("\n=== %s ===\n", nomeTeste)

	requisicao := &orderpb.CreateOrderRequest{
		CostumerId: idCliente,
		OrderItems: itens,
		TotalPrice: precoTotal,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resposta, err := client.Create(ctx, requisicao)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			fmt.Printf("Código do Erro: %s (%d)\n", st.Code().String(), st.Code())
			fmt.Printf("Mensagem de Erro: %s\n", st.Message())
		} else {
			fmt.Printf("Erro desconhecido: %v\n", err)
		}
		return
	}

	fmt.Printf("SUCESSO - ID do Pedido: %d\n", resposta.OrderId)
}

func main() {
	// Conectar ao servidor gRPC
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar: %v", err)
	}
	defer conn.Close()

	cliente := orderpb.NewOrderClient(conn)

	fmt.Println("=== TESTES DE TRATAMENTO DE ERROS ===")

	// Teste 1: Pedido válido (deve ter sucesso)
	itensValidos := []*orderpb.OrderItem{
		{ProductCode: "prod1", UnitPrice: 10.0, Quantity: 2},
		{ProductCode: "prod2", UnitPrice: 5.0, Quantity: 1},
	}
	criarPedido(cliente, 123, itensValidos, 25.0, "Teste 1: Pedido Válido (Total: R$25,00, Qtd: 3)")

	// Teste 2: Pedido com quantidade total > 50 (deve falhar com INVALID_ARGUMENT)
	itensQuantidadeGrande := []*orderpb.OrderItem{
		{ProductCode: "prod1", UnitPrice: 1.0, Quantity: 30},
		{ProductCode: "prod2", UnitPrice: 1.0, Quantity: 25},
	}
	criarPedido(cliente, 124, itensQuantidadeGrande, 55.0, "Teste 2: Pedido com Quantidade Grande (Qtd Total: 55)")

	// Teste 3: Pedido com PreçoTotal > 1000 (deve falhar com INVALID_ARGUMENT)
	itensCaros := []*orderpb.OrderItem{
		{ProductCode: "prod1", UnitPrice: 500.0, Quantity: 1},
		{ProductCode: "prod2", UnitPrice: 600.0, Quantity: 1},
	}
	criarPedido(cliente, 125, itensCaros, 1100.0, "Teste 3: Pedido Caro (Total: R$1100,00)")

	// Teste 4: Pedido válido após erros (deve ter sucesso)
	itensValidos2 := []*orderpb.OrderItem{
		{ProductCode: "prod3", UnitPrice: 15.0, Quantity: 1},
		{ProductCode: "prod4", UnitPrice: 20.0, Quantity: 2},
	}
	criarPedido(cliente, 126, itensValidos2, 55.0, "Teste 4: Pedido Válido Após Erros (Total: R$55,00, Qtd: 3)")

	fmt.Println("\n=== TESTES CONCLUÍDOS ===")
	fmt.Println("Resultados esperados:")
	fmt.Println("- Teste 1: SUCESSO")
	fmt.Println("- Teste 2: INVALID_ARGUMENT (quantidade > 50)")
	fmt.Println("- Teste 3: INVALID_ARGUMENT (preço > 1000)")
	fmt.Println("- Teste 4: SUCESSO")
}
