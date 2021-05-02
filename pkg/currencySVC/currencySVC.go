package currencySVC

import (
	"context"
	"github.com/DaniilOr/currencyParser/cmd/app/dtos"
	"errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}
var ErrCurrencyNotFound = errors.New("currency not found")

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) GetSingle(ctx context.Context, currency string) (interface{}, error) {
	currencyDetails := &dtos.Currency{}
	err := s.pool.QueryRow(ctx, `
		SELECT currency, price FROM currencies WHERE currency = $1
	`, currency).Scan(&currencyDetails.Symbol, &currencyDetails.Price)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, ErrCurrencyNotFound
		}
		return nil, err
	}

	return currencyDetails, nil
}

func (s*Service) GetK(ctx context.Context, k int64) ([]dtos.Currency, error){
	rows, err := s.pool.Query(ctx, `
	SELECT * FROM currencies
	LIMIT $1
`,
		k)
	if err != nil{
		return nil, err
	}
	currencies := []dtos.Currency{}
	for rows.Next() {
		var currency dtos.Currency
		rows.Scan(
			&currency.Symbol,
			&currency.Price,
		)
		currencies = append(currencies, currency)
	}
	if rows.Err() != nil{
		log.Println(rows.Err())
		return nil, rows.Err()
	}
	return currencies, nil
}
// Лучше не использовать этот запрос. Это затратно

func (s*Service) GetAll(ctx context.Context) ([]dtos.Currency, error){
	rows, err := s.pool.Query(ctx, `
	SELECT * FROM currencies
`)
	log.Println("Data has been extracted")
	if err != nil{
		return nil, err
	}
	currencies := []dtos.Currency{}
	for rows.Next() {
		var currency dtos.Currency
		rows.Scan(
			&currency.Symbol,
			&currency.Price,
		)
		currencies = append(currencies, currency)
	}
	if rows.Err() != nil{
		log.Println(rows.Err())
		return nil, rows.Err()
	}
	return currencies, nil
}
func (s*Service) UpdateInfo(ctx context.Context, currency string, price string) (error){
	_, err := s.pool.Exec(ctx, `INSERT INTO currencies (currency, price)  VALUES ($1, $2)
                        ON CONFLICT (currency) DO UPDATE SET price = $2`, currency, price)
	if err != nil{
		log.Println(err)
		return  err
	}
	return nil
}
