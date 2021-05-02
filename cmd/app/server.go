package app

import (
	"context"
	"currencyParser/pkg/currencySVC"
	"currencyParser/pkg/parser"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	currencySVC     *currencySVC.Service
	parserSVC       *parser.Parser
	mux             chi.Router
}

func NewServer(saveSVC *currencySVC.Service, mux chi.Router, parserSVC *parser.Parser) *Server {
	return &Server{currencySVC: saveSVC,  mux: mux, parserSVC: parserSVC}
}
func (s *Server) StartSrapping() error{
	c := cron.New()
	err := c.AddFunc("@every 2s", s.update())
	if err != nil{
		log.Println(err)
		return err
	}
	c.Start()
	return nil
}
func (s *Server) Init() error {
	s.mux.Use(middleware.Logger)
	s.mux.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Get("/all", s.all)
		r.Post("/getSingle", s.single)
		r.Post("/getK", s.getK)
	})

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) all(writer http.ResponseWriter, request *http.Request) {
	data, err := s.currencySVC.GetAll(request.Context())
	if err != nil {
		log.Printf("can't read data from DB: %v", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("can't marshall data: %v", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(respBody)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) single(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Print("can't parse form")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	currency := request.PostForm.Get("currency")

	currencyRecod, err := s.currencySVC.GetSingle(request.Context(), currency)
	if err != nil {
		log.Printf("can't get %s from DB: %v", currency, err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(currencyRecod)
	if err != nil{
		log.Printf("fail to marshal: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}
func (s *Server) getK(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Print("can't parse form")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	k_str := request.PostForm.Get("k")
	k, err :=  strconv.ParseInt(k_str, 0, 64)
	if err != nil{
		log.Printf("cannot convert k = %s to int: %v", k_str, err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	currencies, err := s.currencySVC.GetK(request.Context(), k)
	if err != nil {
		log.Printf("can't get %d records from DB: %v", k, err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(currencies)
	if err != nil{
		log.Printf("fail to marshal: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}
func (s*Server) update(){
		data, err := s.parserSVC.GetUpdate()
		if err != nil{
			log.Println(err)
			return
		}
		for _, item := range data{
			err = s.currencySVC.UpdateInfo(context.Background(), item.Symbol, item.Price)
			if err != nil{
				log.Println(err)
				return
			}
		}
}
