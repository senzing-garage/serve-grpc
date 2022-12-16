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
	var err error = nil
	response := pb.AddDataSourceResponse{}
	return &response, err
}

func (server *G2ConfigServer) Close(ctx context.Context, request *pb.CloseRequest) (*pb.CloseResponse, error) {
	var err error = nil
	response := pb.CloseResponse{}
	return &response, err
}

func (server *G2ConfigServer) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	var err error = nil
	response := pb.CreateResponse{}
	return &response, err
}

func (server *G2ConfigServer) DeleteDataSource(ctx context.Context, request *pb.DeleteDataSourceRequest) (*pb.DeleteDataSourceResponse, error) {
	var err error = nil
	response := pb.DeleteDataSourceResponse{}
	return &response, err
}

func (server *G2ConfigServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	var err error = nil
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	var err error = nil
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigServer) ListDataSources(ctx context.Context, request *pb.ListDataSourcesRequest) (*pb.ListDataSourcesResponse, error) {
	var err error = nil
	response := pb.ListDataSourcesResponse{}
	return &response, err
}

func (server *G2ConfigServer) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	var err error = nil
	response := pb.LoadResponse{}
	return &response, err
}

func (server *G2ConfigServer) Save(ctx context.Context, request *pb.SaveRequest) (*pb.SaveResponse, error) {
	var err error = nil
	response := pb.SaveResponse{}
	return &response, err
}
