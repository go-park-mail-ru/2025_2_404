// Package main initializes and starts the Advertisement gRPC server.
package main

import (
	adhandler "2025_2_404/internal/delivery/grpc/ad"
	"2025_2_404/internal/delivery/grpc/interceptor"
	"2025_2_404/internal/service/ad/config"
	db "2025_2_404/internal/service/ad/connections"
	repoAd "2025_2_404/internal/service/ad/repository/postgres/ad"
	repoBudget "2025_2_404/internal/service/ad/repository/postgres/budget"
	usecase "2025_2_404/internal/service/ad/usecase/ad"
	budget "2025_2_404/internal/service/ad/usecase/budget"
	"2025_2_404/pkg/logger"
	adpb "2025_2_404/protos/gen/go/ad"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("starting Advertisement gRPC server")

	config := config.GetConfig()
	connCfg, err := db.New(config)
	if err != nil {
		log.Fatal("failed to connect to DB", zap.Error(err))
	}
	defer connCfg.CloseAll()

	go func() {
		log.Info("starting metrics server", zap.String("addr", ":9090"))
		http.Handle("/api/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error("metrics server failed", zap.Error(err))
		}
	}()

	// --- Auth connection
	authServiceAddr := "auth:8077"
	authConn, err := grpc.NewClient(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to Auth Service", zap.Error(err))
	}
	defer func() {
		_ = authConn.Close()
	}()

	// --- Инициализация слоёв с логгером
	repoAd := repoAd.New(connCfg.PostgresSQL, log.Named("repo.ad"))
	repoBudget := repoBudget.New(connCfg.PostgresSQL, log.Named("repo.budget"))

	adUC := usecase.New(repoAd, log.Named("usecase.ad"))
	budgetUC := budget.New(repoBudget, log.Named("usecase.budget"))

	authInterceptor, _ := interceptor.InitAuthInterceptor()
	adHandler := adhandler.New(adUC, budgetUC, log.Named("handler.ad"))

	// --- Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":"+config.AppConfig.PortAD)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
	adpb.RegisterAdServServer(grpcServer, adHandler)

	log.Info("gRPC server started",
		zap.String("host", config.AppConfig.Host),
		zap.String("port", config.AppConfig.PortAD),
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("gRPC server failed", zap.Error(err))
	}
}
