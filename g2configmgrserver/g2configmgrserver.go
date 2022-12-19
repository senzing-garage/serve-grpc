package g2configmgrserver

import (
	"context"
	"sync"

	g2sdk "github.com/senzing/g2-sdk-go/g2configmgr"
	pb "github.com/senzing/go-servegrpc/protobuf/g2configmgr"
)

var (
	g2configmgrSingleton *g2sdk.G2configmgrImpl
	g2configmgrSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2configmgr.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2configmgr() *g2sdk.G2configmgrImpl {
	g2configmgrSyncOnce.Do(func() {
		g2configmgrSingleton = &g2sdk.G2configmgrImpl{}
	})
	return g2configmgrSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ConfigmgrServer) AddConfig(ctx context.Context, request *pb.AddConfigRequest) (*pb.AddConfigResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.AddConfig(ctx, request.GetConfigStr(), request.GetConfigComments())
	response := pb.AddConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.Destroy(ctx)
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfig(ctx context.Context, request *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfig(ctx, request.GetConfigID())
	response := pb.GetConfigResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfigList(ctx context.Context, request *pb.GetConfigListRequest) (*pb.GetConfigListResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetConfigList(ctx)
	response := pb.GetConfigListResponse{
		Result: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) GetDefaultConfigID(ctx context.Context, request *pb.GetDefaultConfigIDRequest) (*pb.GetDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	result, err := g2configmgr.GetDefaultConfigID(ctx)
	response := pb.GetDefaultConfigIDResponse{
		ConfigID: result,
	}
	return &response, err
}

func (server *G2ConfigmgrServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.Init(ctx, request.GetModuleName(), request.GetIniParams(), int(request.GetVerboseLogging()))
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) ReplaceDefaultConfigID(ctx context.Context, request *pb.ReplaceDefaultConfigIDRequest) (*pb.ReplaceDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.ReplaceDefaultConfigID(ctx, request.GetOldConfigID(), request.GetNewConfigID())
	response := pb.ReplaceDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) SetDefaultConfigID(ctx context.Context, request *pb.SetDefaultConfigIDRequest) (*pb.SetDefaultConfigIDResponse, error) {
	g2configmgr := getG2configmgr()
	err := g2configmgr.SetDefaultConfigID(ctx, request.GetConfigID())
	response := pb.SetDefaultConfigIDResponse{}
	return &response, err
}
