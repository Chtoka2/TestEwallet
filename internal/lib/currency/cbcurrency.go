//cd internal/lib/currency

package currency

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const cbrURL = "https://www.cbr.ru/scripts/XML_daily.asp"

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}

type Valute struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Value    string `xml:"Value"` // e.g. "90,1234"
}

// GetCBRRate возвращает обменный курс: сколько charCodeTo за 1 charCodeFrom.
// Пример: GetCBRRate(ctx, log, "USD", "EUR") → ~0.93
func GetCBRRate(ctx context.Context, log *slog.Logger, charCodeFrom, charCodeTo string) (float64, error) {
	// Случай одинаковых валют
	if charCodeFrom == charCodeTo {
		return 1.0, nil
	}

	// Обработка RUB как базовой валюты
	var rateFrom, rateTo float64
	var err error

	rateFrom, err = getRateToRUB(ctx, log, charCodeFrom)
	if err != nil {
		return 0, fmt.Errorf("get rate for %s: %w", charCodeFrom, err)
	}

	rateTo, err = getRateToRUB(ctx, log, charCodeTo)
	if err != nil {
		return 0, fmt.Errorf("get rate for %s: %w", charCodeTo, err)
	}

	// Конвертация через RUB: from → RUB → to
	// 1 from = rateFrom RUB
	// 1 to = rateTo RUB ⇒ 1 RUB = 1/rateTo to
	// ⇒ 1 from = rateFrom / rateTo to
	conversionRate := rateFrom / rateTo

	log.Debug("Calculated conversion rate",
		slog.String("from", charCodeFrom),
		slog.String("to", charCodeTo),
		slog.Float64("rate", conversionRate))

	return conversionRate, nil
}

// getRateToRUB возвращает, сколько рублей стоит 1 единица валюты charCode.
func getRateToRUB(ctx context.Context, log *slog.Logger, charCode string) (float64, error) {
	if charCode == "RUB" {
		return 1.0, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", cbrURL, nil)
	if err != nil {
		log.Error("Failed to create CBR request", slog.Any("error", err))
		return 0, fmt.Errorf("create CBR request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to fetch from CBR", slog.Any("error", err))
		return 0, fmt.Errorf("fetch CBR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("CBR returned non-200 status", slog.Int("status", resp.StatusCode))
		return 0, fmt.Errorf("CBR returned status %d", resp.StatusCode)
	}

	// Поддержка кодировки windows-1251 (официальный XML ЦБ в ней)
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		log.Error("Failed to decode CBR XML", slog.Any("error", err))
		return 0, fmt.Errorf("decode CBR XML: %w", err)
	}

	for _, v := range valCurs.Valute {
		if v.CharCode == charCode {
			cleanValue := strings.ReplaceAll(v.Value, ",", ".")
			rate, err := strconv.ParseFloat(cleanValue, 64)
			if err != nil {
				log.Error("Failed to parse currency value", slog.String("value", v.Value), slog.Any("error", err))
				return 0, fmt.Errorf("parse CBR value '%s': %w", v.Value, err)
			}
			result := rate / float64(v.Nominal)
			log.Debug("CBR rate to RUB found", slog.String("currency", charCode), slog.Float64("rate", result))
			return result, nil
		}
	}

	log.Warn("Currency not found in CBR response", slog.String("currency", charCode))
	return 0, fmt.Errorf("currency %s not found in CBR response", charCode)
}