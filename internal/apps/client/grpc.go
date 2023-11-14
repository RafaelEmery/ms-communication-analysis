package client

import (
	"context"
	"log"
	"strings"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	apps "github.com/RafaelEmery/performance-analysis-server/internal/apps/server"
	"google.golang.org/grpc"
)

func HandleGRPC(ctx context.Context, host string, data InteractionData) error {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := apps.NewProductHandlerClient(conn)

	for i := 0; i < data.RequestQuantity; i++ {
		start := time.Now()
		code, err := doProcedureCall(ctx, data.Resource, client)
		if err != nil {
			log.Default().Println(err)
			if data.RequestQuantity == 1 {
				return err
			}
		}

		elapsed := time.Since(start).String()
		log.Default().Printf("[%s] %s - %s", code, strings.ToUpper(string(data.Resource[0]))+data.Resource[1:], elapsed)
	}

	return nil
}

func doProcedureCall(ctx context.Context, resource string, client apps.ProductHandlerClient) (string, error) {
	var err error

	switch resource {
	case createResource:
		var product domain.Product
		fp := *product.Fake()

		req := &apps.CreateProductRequest{
			Name:              fp.Name,
			Sku:               fp.SKU,
			SellerName:        fp.SellerName,
			Price:             float32(fp.Price),
			AvailableDiscount: float32(fp.AvailableDiscount),
			AvailableQuantity: int32(fp.AvailableQuantity),
			SalesQuantity:     int32(fp.SalesQuantity),
			Active:            fp.Active,
		}

		_, err = client.Create(ctx, req)
	case reportResource:
		req := &apps.EmptyRequest{}
		_, err = client.Report(ctx, req)
	case getByDiscountResource:
		req := &apps.EmptyRequest{}
		_, err = client.GetByDiscount(ctx, req)
	}

	return "ok", err
}
