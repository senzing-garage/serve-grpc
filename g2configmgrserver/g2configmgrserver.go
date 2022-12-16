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
	var err error = nil
	response := pb.AddConfigResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	var err error = nil
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfig(ctx context.Context, request *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	var err error = nil
	response := pb.GetConfigResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetConfigList(ctx context.Context, request *pb.GetConfigListRequest) (*pb.GetConfigListResponse, error) {
	var err error = nil
	response := pb.GetConfigListResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) GetDefaultConfigID(ctx context.Context, request *pb.GetDefaultConfigIDRequest) (*pb.GetDefaultConfigIDResponse, error) {
	var err error = nil
	response := pb.GetDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	var err error = nil
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) ReplaceDefaultConfigID(ctx context.Context, request *pb.ReplaceDefaultConfigIDRequest) (*pb.ReplaceDefaultConfigIDResponse, error) {
	var err error = nil
	response := pb.ReplaceDefaultConfigIDResponse{}
	return &response, err
}

func (server *G2ConfigmgrServer) SetDefaultConfigID(ctx context.Context, request *pb.SetDefaultConfigIDRequest) (*pb.SetDefaultConfigIDResponse, error) {
	var err error = nil
	response := pb.SetDefaultConfigIDResponse{}
	return &response, err
}
