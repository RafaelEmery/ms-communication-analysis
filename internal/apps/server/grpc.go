package apps

import (
	"context"
	"log"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	// TODO: understand the HandlerServer below
	// s   p.UnimplementedProductHandlerServer
	c   Creator
	rg  ReportGenerator
	dpg ProductByDiscountGetter
}

func NewGRPCServer(c Creator, rg ReportGenerator, dpg ProductByDiscountGetter) GRPCServer {
	return GRPCServer{c: c, rg: rg, dpg: dpg}
}

func (s GRPCServer) Create(ctx context.Context, request *CreateProductRequest) (*CreateProductResponse, error) {
	log.Default().Println("called grpc method Create")
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
		return &CreateProductResponse{}, status.Error(codes.Internal, err.Error())
	}

	logMemStats()

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
	log.Default().Println("called grpc method Report")
	o, err := s.rg.GenerateReport(ctx)
	if err != nil {
		return &GenerateReportResponse{}, status.Error(codes.Internal, err.Error())
	}

	logMemStats()

	return &GenerateReportResponse{FileName: o}, nil
}

func (s GRPCServer) GetByDiscount(ctx context.Context, in *EmptyRequest) (*GetByDiscountResponse, error) {
	log.Default().Println("called grpc method GetByDiscount")
	o, err := s.dpg.GetByDiscount(ctx)
	if err != nil {
		return &GetByDiscountResponse{}, status.Error(codes.Internal, err.Error())
	}

	ps := make([]*Product, 0)
	for _, v := range o {
		p := &Product{
			Id:                v.ID,
			Name:              v.Name,
			Sku:               v.SKU,
			SellerName:        v.SellerName,
			Price:             float32(v.Price),
			AvailableDiscount: float32(v.AvailableDiscount),
			AvailableQuantity: int32(v.AvailableQuantity),
			SalesQuantity:     int32(v.SalesQuantity),
			Active:            v.Active,
			DiscountApplied:   v.DiscountApplied,
			CreatedAt:         v.CreatedAt.String(),
			UpdatedAt:         v.UpdatedAt.String(),
		}

		ps = append(ps, p)
	}

	logMemStats()

	return &GetByDiscountResponse{Products: ps}, nil
}

// TODO: validate function definition and it's utility
func (s GRPCServer) mustEmbedUnimplementedProductHandlerServer() {}
