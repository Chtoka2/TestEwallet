package currency

// cd internal/lib/currency
import (
	"context"
	"log/slog"
	"os"
	"testing"
)

type test struct{
	Name string
	CharCodeFrom, CharCodeTo string
	wantError bool
}

func TestCurrency(t *testing.T){
	tests := []test{
		{
			Name: "Succesful",
			CharCodeFrom: "USD",
			CharCodeTo:"RUB",
			wantError: false,
		},
		{
			Name: "Value of RUB",
			CharCodeFrom: "RUB",
			CharCodeTo: "USD",
			wantError: true,
		},
		{
			Name: "Inccorect value",
			CharCodeFrom: "u",
			CharCodeTo:"A",
			wantError: true,
		},
	}
	log := slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	for _, i := range tests{
		res, err := GetCBRRate(context.Background(), log, i.CharCodeFrom, i.CharCodeTo)
		if err != nil{
			if i.wantError == false{
				t.Errorf("Not wanted error %v", err)
			}
			log.Info("", slog.String("op", i.Name), slog.String("error", err.Error()))
		}
		log.Info("Course of value", slog.Float64(i.CharCodeFrom, res), slog.Float64(i.CharCodeTo, res))
	}
}