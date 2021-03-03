// @Author LiuYong
// @Created at 2021-01-28
// @Modified at 2021-01-28
package orderRPCClient

import (
	"common/rpc_manage"
	"common/tools/logging"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"rpc/pay/proto/orderRPCpb"
	"sync"
	"time"
)

type OrderRPCClient struct {
	pbClient *orderRPCpb.OrderRPCServiceClient
	// 支付成功后进行处理的函数
	dealPayOK func(payOKInfo *PayOKOrderInfo)
}

type PayOKOrderInfo struct {
	UserID    uint
	OutsideID string
	AffairID  string
	Money     int
	State     int
}

// 唯一
var (
	client           *OrderRPCClient
	kafkaDefaultHost = "kafka:9092"
	targetService    = "nginx:18082"
	// 在重设kafka host是停止上一个线程
	runFlag = true
	wg      sync.WaitGroup
	once    sync.Once
	lock    sync.Mutex
)

const (
	targetServiceName_localhost = ":8082"
	topic                       = "PayOK"
	partition                   = 0
)

// NewClient 创建一个 OrderRPCClient
//
// Deprecated
func NewClient() (*OrderRPCClient, error) {
	once.Do(func() {
		if client != nil {
			return
		}
		client = &OrderRPCClient{pbClient: nil, dealPayOK: nil}
		conn, err := rpc_manage.NewGRPCClientConn(targetService) // 如果出错可能运行在本地, 尝试使用本地连接
		if err != nil {
			logging.Error(err)
			client = nil
			return
		}
		tclient := orderRPCpb.NewOrderRPCServiceClient(conn)
		client.pbClient = &tclient

		// 连接kafka并接受消息
		go runRecvKafka(context.Background(), client)
	})
	return client, nil
}

func tryDirectConnent() *OrderRPCClient {
	if client != nil {
		lock.Lock()
		if client == nil {
			return nil
		}
		logging.Info("order RPC reConnent to localhost:8082")
		conn, err := rpc_manage.NewGRPCClientConn(targetServiceName_localhost) // 如果出错可能运行在本地, 尝试使用本地连接
		if err != nil {
			logging.Error(err)
			client = nil
			return nil
		}
		tclient := orderRPCpb.NewOrderRPCServiceClient(conn)
		client.pbClient = &tclient
		lock.Unlock()
	}
	return client
}

// NewClientWithMQHost 创建一个 OrderRPCClient 并设置消息队列的服务地址
//
// 在docker中默认设置未kafka:9092, 如果是本地运行可以使用localhost:9092
// Deprecated
func NewClientWithMQHost(MQHost string) (*OrderRPCClient, error) {
	kafkaDefaultHost = MQHost
	return NewClient()
}

func NewClientWithTargetAndMQHost(target, MQHost string) (*OrderRPCClient, error) {
	targetService = target
	return NewClientWithMQHost(MQHost)
}

// runRecvKafka 运行接受协程, 只能在协程中运行
func runRecvKafka(ctx context.Context, cli *OrderRPCClient) {
	wg.Add(1)
	defer wg.Done()
	// 创建kafka消费者
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaDefaultHost},
		Topic:     topic,
		Partition: partition,
		//StartOffset:            kafka.LastOffset,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		// 接受消息的等待时间, 如果不设置默认会等待9秒后才向kafka服务器请求接受消息, 坑!
		MaxWait: 200 * time.Millisecond,
	})
	// 设置消息偏移, 避免收到之前的消息
	err := kafkaReader.SetOffset(kafka.LastOffset)
	if err != nil {
		logging.Error(err)
	}

	for runFlag {
		select {
		case <-ctx.Done():
			break
		default:
			msg, err := kafkaReader.ReadMessage(context.Background())
			if err != nil {
				logging.Error("order get kafka message error: ", err)
				break
			}
			if cli.dealPayOK != nil {
				var payOK PayOKOrderInfo
				err := json.Unmarshal(msg.Value, &payOK)
				if err != nil {
					logging.Error("order get kafka message error: ", err)
					break
				}
				cli.dealPayOK(&payOK)
			}
		}
	}

	if err := kafkaReader.Close(); err != nil {
		logging.Error("failed to close connection:", err)
	}
}

// SetMQHost 重设队列的地址, 该方法可能会出现等待, 最好在NewClientWithMQHost时设置好
func (c *OrderRPCClient) SetMQHost(MQHost string) {
	runFlag = false
	wg.Wait()
	runFlag = true
	kafkaDefaultHost = MQHost
	go runRecvKafka(context.Background(), c)
}

// SetDealPayOK 设置支付成功后的出处理函数
func (c *OrderRPCClient) SetDealPayOK(deal func(payOK *PayOKOrderInfo)) {
	c.dealPayOK = deal
}

func (c *OrderRPCClient) Create(info *orderRPCpb.CreateRequest) (*orderRPCpb.CreateRespond, error) {
	tclient := *c.pbClient
	resp, err := tclient.Create(context.Background(), info)
	if err != nil {
		// logging.Error(err)
		// tryDirectConnent()
		// tclient := *c.pbClient
		// resp, err := tclient.Create(context.Background(), info)
		// if err != nil {
		// 	return nil, err
		// }
		return resp, err
	}
	return resp, nil
}

func (c *OrderRPCClient) Read(info *orderRPCpb.SearchCondition) (*orderRPCpb.ReadRespond, error) {
	tclient := *c.pbClient
	resp, err := tclient.Read(context.Background(), info)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateState(info *orderRPCpb.UpdateStateRequest) (*orderRPCpb.Respond, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateState(context.Background(), info)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *OrderRPCClient) UpdateStateWithRelativeOrder(info *orderRPCpb.UpdateStateWithRRequest) (*orderRPCpb.Respond, error) {
	tclient := *c.pbClient
	resp, err := tclient.UpdateStateWithRelativeOrder(context.Background(), info)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *OrderRPCClient) ExistNoFinishOrder(info *orderRPCpb.SearchCondition) (*orderRPCpb.ExistNoFinishOrderRespond, error) {
	tclient := *c.pbClient
	resp, err := tclient.ExistNoFinishOrder(context.Background(), info)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *OrderRPCClient) Refund(request *orderRPCpb.RefundRequest) (*orderRPCpb.Respond, error) {
	tclient := *c.pbClient
	resp, err := tclient.Refund(context.Background(), request)
	if err != nil {
		return resp, err
	}
	return resp, err
}
