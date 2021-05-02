package main

import (
	"context"
	"github.com/DaniilOr/currencyParser/cmd/app"
	"github.com/DaniilOr/currencyParser/pkg/currencySVC"
	"github.com/DaniilOr/currencyParser/pkg/parser"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort               = "9999"
	defaultHost               = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@currenciesdb:5432/db"
	defaultParseURL = "https://api.binance.com/api/v3/ticker/price"
)
func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}
	url, ok := os.LookupEnv("PARSE_URL")
	if !ok{
		url = defaultParseURL
	}
	if err := execute(net.JoinHostPort(host, port), dsn, url); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, dsn string, url string) error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}
	mux := chi.NewRouter()
	mycurrencySVC := currencySVC.NewService(pool)
	parserSVC := parser.InitService(url)
	application := app.NewServer(mycurrencySVC, mux, parserSVC)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}
	err = application.StartScrapping()
	if err != nil{
		return err
	}
	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}