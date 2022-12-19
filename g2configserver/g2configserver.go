package g2configserver

import (
	"context"
	"sync"

	g2sdk "github.com/senzing/g2-sdk-go/g2config"
	pb "github.com/senzing/go-servegrpc/protobuf/g2config"
)

var (
	g2configSingleton *g2sdk.G2configImpl
	g2configSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2config() *g2sdk.G2configImpl {
	g2configSyncOnce.Do(func() {
		g2configSingleton = &g2sdk.G2configImpl{}
	})
	return g2configSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ConfigServer) AddDataSource(ctx context.Context, request *pb.AddDataSourceRequest) (*pb.AddDataSourceResponse, error) {
	g2config := getG2config()
	result, err := g2config.AddDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := pb.AddDataSourceResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Close(ctx context.Context, request *pb.CloseRequest) (*pb.CloseResponse, error) {
	g2config := getG2config()
	err := g2config.Close(ctx, uintptr(request.GetConfigHandle()))
	response := pb.CloseResponse{}
	return &response, err
}

func (server *G2ConfigServer) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	g2config := getG2config()
	result, err := g2config.Create(ctx)
	response := pb.CreateResponse{
		Result: int64(result),
	}
	return &response, err
}

func (server *G2ConfigServer) DeleteDataSource(ctx context.Context, request *pb.DeleteDataSourceRequest) (*pb.DeleteDataSourceResponse, error) {
	g2config := getG2config()
	err := g2config.DeleteDataSource(ctx, uintptr(request.GetConfigHandle()), request.GetInputJson())
	response := pb.DeleteDataSourceResponse{}
	return &response, err
}

func (server *G2ConfigServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2config := getG2config()
	err := g2config.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2config := getG2config()
	err := g2config.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigServer) ListDataSources(ctx context.Context, request *pb.ListDataSourcesRequest) (*pb.ListDataSourcesResponse, error) {
	g2config := getG2config()
	result, err := g2config.ListDataSources(ctx, uintptr(request.GetConfigHandle()))
	response := pb.ListDataSourcesResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigServer) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	g2config := getG2config()
	err := g2config.Load(ctx, uintptr(request.GetConfigHandle()), (request.GetJsonConfig()))
	response := pb.LoadResponse{}
	return &response, err
}

func (server *G2ConfigServer) Save(ctx context.Context, request *pb.SaveRequest) (*pb.SaveResponse, error) {
	g2config := getG2config()
	result, err := g2config.Save(ctx, uintptr(request.GetConfigHandle()))
	response := pb.SaveResponse{
		Result: result,
	}
	return &response, err
}
