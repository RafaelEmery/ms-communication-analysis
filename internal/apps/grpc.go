package apps

import (
	"context"
	"log"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

// TODO: improve .proto and grpc files structure
// TODO: change gRPC file test function name
// TODO: define error message as struct on .proto and grpc files

type GRPCServer struct {
	// s   p.UnimplementedProductHandlerServer
	c   Creator
	rg  ReportGenerator
	dpg ProductByDiscountGetter
}

func NewGRPCServer(c Creator, rg ReportGenerator, dpg ProductByDiscountGetter) GRPCServer {
	return GRPCServer{c: c, rg: rg, dpg: dpg}
}

func (s GRPCServer) Create(ctx context.Context, request *CreateProductRequest) (*CreateProductResponse, error) {
	log.Default().Println("create product called with request: ", request)

	// TODO: validate use case correct apply here
	i := domain.Product{
		Name:              request.Name,
		SKU:               request.Sku,
		SellerName:        request.SellerName,
		Price:             float64(request.Price),
		AvailableDiscount: float64(request.AvailableDiscount),
		AvailableQuantity: int(request.AvailableQuantity),
		SalesQuantity:     int(request.SalesQuantity),
		Active:            request.Active,
	}

	o, err := s.c.Create(ctx, i)
	if err != nil {
		log.Default().Println("error on creating use case", err.Error())
		return &CreateProductResponse{}, err
	}

	return &CreateProductResponse{
		Id:                o.ID,
		Name:              o.Name,
		Sku:               o.SKU,
		SellerName:        o.SellerName,
		Price:             float32(o.Price),
		AvailableDiscount: float32(o.AvailableDiscount),
		AvailableQuantity: int32(o.AvailableQuantity),
		SalesQuantity:     int32(o.SalesQuantity),
		Active:            o.Active,
		DiscountApplied:   o.DiscountApplied,
		CreatedAt:         o.CreatedAt.String(),
		UpdatedAt:         o.UpdatedAt.String(),
	}, nil
}

func (s GRPCServer) Report(ctx context.Context, in *EmptyRequest) (*GenerateReportResponse, error) {
	return &GenerateReportResponse{}, nil
}

func (s GRPCServer) GetByDiscount(ctx context.Context, in *EmptyRequest) (*GetByDiscountResponse, error) {
	return &GetByDiscountResponse{}, nil
}

// TODO: validate function definition and it's utility
func (s GRPCServer) mustEmbedUnimplementedProductHandlerServer() {}
