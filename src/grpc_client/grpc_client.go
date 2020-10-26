package grpc_client

import (
	"community-sig/constants"
	"community-sig/logging"
	"community-sig/protobuffer_def"
	"context"
	"sync"
	"time"

	"community-sig/config"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"
	"github.com/micro/go-plugins/registry/zookeeper/v2"
	rgrpc "google.golang.org/grpc"
)

var (
	statusClient  protobuffer_def.StatusServerService
	statusOnce    = &sync.Once{}
	statusCancel  context.CancelFunc
	statusCtx     context.Context
	statusService service.Service
	//信令服务器 rpc连接池
	sigClients map[string]*SigClient
)

func init() {
	sigClients = make(map[string]*SigClient)
}

func StatusClient() protobuffer_def.StatusServerService {
	statusOnce.Do(func() {
		r := zookeeper.NewRegistry(func(op *registry.Options) {
			op.Addrs = config.GetConf().RegisterCenter.Address
			op.Context = context.Background()
			op.Timeout = time.Second * 5
		})

		statusCtx, statusCancel = context.WithCancel(context.Background())
		statusService = grpc.NewService(
			service.Name(constants.GetStatusRpcClientName()),
			service.Registry(r),
			service.Context(statusCtx),
		)

		err := statusService.Client().Init(client.Retries(3), client.PoolSize(200),
			client.PoolTTL(time.Second*20), client.RequestTimeout(time.Second*5))
		if err != nil {
			logging.Logger.Info("service client init failed, err is ", err)
			return
		}

		statusClient = protobuffer_def.NewStatusServerService(config.GetConf().Base.PeerServiceName, statusService.Client())

		logging.Logger.Info("peer service name is ", config.GetConf().Base.PeerServiceName)
		logging.Logger.Infof("sig service client is %v", statusClient)
	})

	return statusClient
}

func CloseStatusServiceClient() {
	if statusClient != nil {
		statusCancel()
		statusClient = nil
	}
}

func SendRpcMsg(baseReq *protobuffer_def.BaseRequest) (*protobuffer_def.BaseResponse, error) {
	baseRsp, err := StatusClient().BaseInterface(context.Background(), baseReq)
	if err != nil {
		logging.Logger.Info("query failed, err is ", err)
		return nil, err
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("query failed, code is ", baseRsp.Code)
		return baseRsp, nil
	}

	return baseRsp, nil
}

//信令服务器rpc客户端结构
type SigClient struct {
	Conn   *rgrpc.ClientConn
	Client protobuffer_def.ComSigServerClient
}

func (sg *SigClient) PostRpcMsg(baseReq *protobuffer_def.BaseRequest) (*protobuffer_def.BaseResponse, error) {
	baseRsp, err := sg.Client.BaseInterface(context.Background(), baseReq)
	if err != nil {
		logging.Logger.Info("query failed, err is ", err)
		return nil, err
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		logging.Logger.Info("query failed, code is ", baseRsp.Code)
		return baseRsp, nil
	}

	return baseRsp, nil
}

func GetSigRPCClient(addr string) (*SigClient, error) {
	if addr == "" {
		logging.Logger.Info("error rpc addr is empty!!!")
		return nil, constants.ErrRpcAddrEmpty
	}

	val, ok := sigClients[addr]
	if ok {
		logging.Logger.Info("get rpc client from sigservice map ", val)
		return val, nil
	}

	conn, err := rgrpc.Dial(addr, rgrpc.WithInsecure())
	if err != nil {
		logging.Logger.Infof("did not connect: %v ", err)
		return nil, err
	}

	sc := protobuffer_def.NewComSigServerClient(conn)
	sigClient := &SigClient{}
	sigClient.Conn = conn
	sigClient.Client = sc

	sigClients[addr] = sigClient
	logging.Logger.Infof("add %s into client map success", addr)
	return sigClient, nil
}

func CloseSigClients() {
	if sigClients == nil {
		return
	}
	for _, val := range sigClients {
		val.Conn.Close()
	}
	sigClients = nil
}

/*
func Start() {
	r := zookeeper.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:2181"}
		op.Context = context.Background()
		op.Timeout = time.Second * 5
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create GRPC service
	service := grpc.NewService(
		service.Name("test.client2"),
		service.Registry(r),
		service.Context(ctx),
	)

	err := service.Client().Init(client.Retries(3), client.PoolSize(200), client.PoolTTL(time.Second*20), client.RequestTimeout(time.Second*5))
	if err != nil {
		fmt.Println("service client init failed, err is ", err)
		return
	}

	test := protobuffer_def.NewComSigServerService(config.GetConf().Base.PeerServiceName, service.Client())
	fmt.Println("test is ", test)
	baseReq := &protobuffer_def.BaseRequest{}
	baseReq.RequestId = "user_info_reg"
	baseReq.C = protobuffer_def.CMD_REGISTER_STATUS
	queryReq := &protobuffer_def.RegisterStatusRequest{}
	queryReq.Identity = "test01"
	queryReq.Phone = "18301152007"
	queryReq.RegisterInfo = "127.0.0.1:8092"
	queryReq.RoomList = []string{"101", "102"}
	queryReq.UserAvator = "avator.png"
	queryReq.UserName = "test-user"

	body, err := ptypes.MarshalAny(queryReq)
	if err != nil {
		fmt.Println("proto marshal failed")
		return
	}
	baseReq.Body = body

	fmt.Println("baseReq is ", baseReq)
	baseRsp, err := test.BaseInterface(context.Background(), baseReq)
	if err != nil {
		fmt.Println("query failed, err is ", err)
		return
	}

	if baseRsp.Code != protobuffer_def.ReturnCode_SUCCESS {
		fmt.Println("query failed, code is ", baseRsp.Code)
		return
	}

	queryRsp := &protobuffer_def.RegisterStatusResponse{}
	err = ptypes.UnmarshalAny(baseRsp.GetBody(), queryRsp)
	if err != nil {
		fmt.Println("proto unmarsh failed")
		return
	}
	fmt.Println("query rsp is ", queryRsp)
}
*/
