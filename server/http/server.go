package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/ivan1993spb/snake-server/client"
	"github.com/ivan1993spb/snake-server/config"
	"github.com/ivan1993spb/snake-server/connections"
	"github.com/ivan1993spb/snake-server/server/http/handlers"
	"github.com/ivan1993spb/snake-server/server/http/middlewares"
)

const ServerName = "Snake-Server"

const logName = "api"

type Server struct {
	Addr    string
	Handler http.Handler
	Logger  *logrus.Logger

	GroupManager *connections.ConnectionGroupManager
}

func NewServer(cfg config.Config, groupManager *connections.ConnectionGroupManager, logger *logrus.Logger,
	author, license, version, build string) *Server {
	srv := &Server{
		Addr:         cfg.Server.Address,
		Logger:       logger,
		GroupManager: groupManager,
	}

	// TODO: Refactor this function.

	srv.InitRoutes(
		cfg.Server.EnableWeb,
		cfg.Server.EnableBroadcast,
		cfg.Server.ForbidCORS,
		author,
		license,
		version,
		build,
	)

	return srv
}

func (srv *Server) InitRoutes(enableWeb, enableBroadcast, forbidCORS bool, author, license, version, build string) {
	// TODO: Refactor this function.

	rootRouter := mux.NewRouter().StrictSlash(true)
	rootRouter.Path("/metrics").Handler(promhttp.Handler())
	rootRouter.Path(handlers.URLRouteOpenAPI).Handler(handlers.NewOpenAPIHandler())
	if enableWeb {
		rootRouter.Path(client.URLRouteServerEndpoint).Handler(http.RedirectHandler(client.URLRouteClient, http.StatusFound))
		rootRouter.PathPrefix(client.URLRouteClient).Handler(negroni.New(gzip.Gzip(gzip.DefaultCompression), negroni.Wrap(client.NewHandler())))
	} else {
		rootRouter.Path(handlers.URLRouteWelcome).Methods(handlers.MethodWelcome).Handler(handlers.NewWelcomeHandler(srv.Logger))
	}
	rootRouter.NotFoundHandler = handlers.NewNotFoundHandler(srv.Logger)

	// Web-Socket routes
	wsRouter := rootRouter.PathPrefix("/ws").Subrouter()
	wsRouter.Path(handlers.URLRouteGameWebSocketByID).Methods(handlers.MethodGame).Handler(handlers.NewGameWebSocketHandler(srv.Logger, srv.GroupManager))

	// API routes
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiRouter.Path(handlers.URLRouteGetInfo).Methods(handlers.MethodGetInfo).Handler(handlers.NewGetInfoHandler(srv.Logger, author, license, version, build))
	apiRouter.Path(handlers.URLRouteGetCapacity).Methods(handlers.MethodGetCapacity).Handler(handlers.NewGetCapacityHandler(srv.Logger, srv.GroupManager))
	apiRouter.Path(handlers.URLRouteCreateGame).Methods(handlers.MethodCreateGame).Handler(handlers.NewCreateGameHandler(srv.Logger, srv.GroupManager))
	apiRouter.Path(handlers.URLRouteGetGameByID).Methods(handlers.MethodGetGame).Handler(handlers.NewGetGameHandler(srv.Logger, srv.GroupManager))
	apiRouter.Path(handlers.URLRouteDeleteGameByID).Methods(handlers.MethodDeleteGame).Handler(handlers.NewDeleteGameHandler(srv.Logger, srv.GroupManager))
	apiRouter.Path(handlers.URLRouteGetGames).Methods(handlers.MethodGetGames).Handler(handlers.NewGetGamesHandler(srv.Logger, srv.GroupManager))
	if enableBroadcast {
		apiRouter.Path(handlers.URLRouteBroadcast).Methods(handlers.MethodBroadcast).Handler(handlers.NewBroadcastHandler(srv.Logger, srv.GroupManager))
	}
	apiRouter.Path(handlers.URLRouteGetObjects).Methods(handlers.MethodGetObjects).Handler(handlers.NewGetObjectsHandler(srv.Logger, srv.GroupManager))
	apiRouter.Path(handlers.URLRoutePing).Methods(handlers.MethodPing).Handler(handlers.NewPingHandler(srv.Logger))

	n := negroni.New(
		middlewares.NewRecovery(srv.Logger),
		middlewares.NewServerInfo(ServerName, version, build),
		middlewares.NewLogger(srv.Logger, logName),
	)

	if !forbidCORS {
		n.Use(middlewares.NewCORS())
	}

	n.UseHandler(rootRouter)

	srv.Handler = n
}

func (srv *Server) ListenAndServe() error {
	return http.ListenAndServe(srv.Addr, srv.Handler)
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return http.ListenAndServeTLS(srv.Addr, certFile, keyFile, srv.Handler)
}
