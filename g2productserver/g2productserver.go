package g2productserver

import (
	"context"
	"sync"

	g2sdk "github.com/senzing/g2-sdk-go/g2product"
	pb "github.com/senzing/go-servegrpc/protobuf/g2product"
)

var (
	g2productSingleton *g2sdk.G2productImpl
	g2productSyncOnce  sync.Once
)

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func getG2product() *g2sdk.G2productImpl {
	g2productSyncOnce.Do(func() {
		g2productSingleton = &g2sdk.G2productImpl{}
	})
	return g2productSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (server *G2ProductServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	var err error = nil
	response := pb.DestroyResponse{}
	return &response, err
}

func (server *G2ProductServer) Init(ctx context.Context, request *pb.InitRequest) (*pb.InitResponse, error) {
	var err error = nil
	response := pb.InitResponse{}
	return &response, err
}

func (server *G2ProductServer) License(ctx context.Context, request *pb.LicenseRequest) (*pb.LicenseResponse, error) {
	var err error = nil
	response := pb.LicenseResponse{}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseFile(ctx context.Context, request *pb.ValidateLicenseFileRequest) (*pb.ValidateLicenseFileResponse, error) {
	var err error = nil
	response := pb.ValidateLicenseFileResponse{}
	return &response, err
}

func (server *G2ProductServer) ValidateLicenseStringBase64(ctx context.Context, request *pb.ValidateLicenseStringBase64Request) (*pb.ValidateLicenseStringBase64Response, error) {
	var err error = nil
	response := pb.ValidateLicenseStringBase64Response{}
	return &response, err
}

func (server *G2ProductServer) Version(ctx context.Context, request *pb.VersionRequest) (*pb.VersionResponse, error) {
	var err error = nil
	response := pb.VersionResponse{}
	return &response, err
}
