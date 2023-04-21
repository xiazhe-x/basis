package base

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xiazhe-x/basis"
	"time"
)

var client *redis.Client

func RedisClient() *redis.Client {
	Check(client)
	return client
}

//redis starter，并且设置为全局
type GoRedisStarter struct {
	basis.BaseStarter
}

func (r *GoRedisStarter) Setup(ctx basis.StarterContext) {
	conf := ctx.Props()
	client = redis.NewClient(&redis.Options{
		Network:  "tcp",                                           //网络类型，tcp or unix，默认tcp
		Addr:     conf.GetDefault("redis.addr", "127.0.0.1:6379"), //主机名+冒号+端口，默认localhost:6379
		Password: conf.GetDefault("redis.pwd", ""),                //密码
		DB:       conf.GetIntDefault("redis.db", 0),               // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     16, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 64, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		/*//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 120 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,   //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,   //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
			//可自定义连接函数
			Dialer: func() (net.Conn, error) {
				netDialer := &net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 5 * time.Minute,
				}
				return netDialer.Dial("tcp", "127.0.0.1:6379")
			},
			//钩子函数
			OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
				fmt.Printf("conn=%v\n", conn)
				return nil
			},
		*/
	})

	_, err := client.Do(context.Background(), "ping").Result()
	if err != nil {
		logrus.Panic("redis：", err)
		panic(err)
	}
}
