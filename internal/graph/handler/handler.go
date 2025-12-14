package handler

import (
	"context"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/graph/generated"
	"github.com/xinliangnote/go-gin-api/internal/graph/resolvers"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"
)

var _ Gql = (*gql)(nil)

type Gql interface {
	i()
	Playground() core.HandlerFunc
	Query() core.HandlerFunc
}

type gql struct {
	logger *zap.Logger
	db     mysql.Repo
	cache  redis.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Gql {
	return &gql{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}

func (g *gql) i() {}

// Query 一个GraphQL查询处理器工厂函数，负责创建和配置GraphQL服务器的核心处理逻辑
func (g *gql) Query() core.HandlerFunc {

	// 定义扩展字段
	extensions := make(map[string]interface{})

	// 创建GraphQL可执行schema
	// 注入日志、数据库、缓存等依赖
	// 生成GraphQL查询处理器
	h := handler.New(generated.NewExecutableSchema(
		resolvers.NewRootResolvers(g.logger, g.db, g.cache)),
	)

	// 添加Websocket传输支持
	// 启用KeepAlivePingInterval，每10秒发送一次心跳包，保持连接活跃
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	// 设置 transport
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	// 设置查询缓存，缓存1000个查询结果
	h.SetQueryCache(lru.New(1000))

	// 启用侧边栏文档
	h.Use(extension.Introspection{})

	// 启用自动持久化查询
	// 自动缓存100个查询，避免重复解析
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return func(c core.Context) {
		var responses interface{}

		defer func() {
			// 设置GraphQL响应日志；集成到项目的统一日志系统中
			c.GraphPayload(responses)
		}()

		// 设置 core trace_id
		extensions["trace_id"] = c.Trace().ID()

		// 响应拦截器：在GraphQL响应返回前，添加自定义扩展字段
		h.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			resp := next(ctx)
			resp.Extensions = extensions // 添加trace_id扩展字段
			responses = resp             // 记录响应用于日志
			return resp
		})

		// 将项目的core.Context注入到GraphQL上下文中
		// 确保GraphQL解析器可以访问项目框架功能
		coreContext := context.WithValue(c.Request().Context(), resolvers.CoreContextKey, c)
		h.ServeHTTP(c.ResponseWriter(), c.Request().WithContext(coreContext))
	}
}

func (g *gql) Playground() core.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql/query")
	return func(c core.Context) {
		h.ServeHTTP(c.ResponseWriter(), c.Request())
	}
}
