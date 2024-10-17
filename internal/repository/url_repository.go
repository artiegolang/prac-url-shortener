package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"practicum-middle/pkg/database"
)

type URLRepository struct {
	db *database.DB
}

type URLPair struct {
	ShortID     string
	OriginalURL string
}

func NewURLRepository(db *database.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) SaveURL(ctx context.Context, shortID, originalURL string) (string, bool, error) {
	var returnedShortID string
	sqlQuery := `
        INSERT INTO urls (short_id, original_url)
        VALUES ($1, $2)
        ON CONFLICT (original_url) DO NOTHING
        RETURNING short_id;
    `
	err := r.db.Pool.QueryRow(ctx, sqlQuery, shortID, originalURL).Scan(&returnedShortID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Вставка не произошла, URL уже существует
			// Получаем существующий short_id
			err = r.db.Pool.QueryRow(ctx, "SELECT short_id FROM urls WHERE original_url = $1", originalURL).Scan(&returnedShortID)
			if err != nil {
				return "", false, err
			}
			return returnedShortID, true, nil // true означает, что URL уже существовал
		}
		return "", false, err // Другая ошибка
	}
	return returnedShortID, false, nil // false означает, что запись была вставлена
}

func (r *URLRepository) SaveURLsBatch(ctx context.Context, urlPairs []URLPair) (resultMap map[string]string, err error) {
	// Начинаем транзакцию
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Используем pgx.Batch для пакетной вставки
	batch := &pgx.Batch{}
	for _, pair := range urlPairs {
		batch.Queue(`
            INSERT INTO urls (short_id, original_url)
            VALUES ($1, $2)
            ON CONFLICT (original_url) DO NOTHING
        `, pair.ShortID, pair.OriginalURL)
	}

	// Отправляем батч и получаем BatchResults
	br := tx.SendBatch(ctx, batch)

	// Нужно обработать все результаты из батча
	for range urlPairs {
		_, err = br.Exec()
		if err != nil {
			br.Close()
			return nil, err
		}
	}

	// Закрываем BatchResults после использования
	err = br.Close()
	if err != nil {
		return nil, err
	}

	// После того как батч обработан и закрыт, можем выполнять другие запросы

	// Собираем все original_url
	originalURLs := make([]string, len(urlPairs))
	for i, pair := range urlPairs {
		originalURLs[i] = pair.OriginalURL
	}

	// Получаем short_id для всех original_url внутри транзакции
	rows, err := tx.Query(ctx, `
        SELECT original_url, short_id FROM urls WHERE original_url = ANY($1)
    `, originalURLs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Формируем карту соответствия original_url -> short_id
	resultMap = make(map[string]string)
	for rows.Next() {
		var originalURL, shortID string
		err = rows.Scan(&originalURL, &shortID)
		if err != nil {
			return nil, err
		}
		resultMap[originalURL] = shortID
	}

	// Проверяем наличие ошибок после итерации по rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resultMap, nil
}

func (r *URLRepository) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	var originalURL string
	err := r.db.Pool.QueryRow(ctx, "SELECT original_url FROM urls WHERE short_id=$1", shortID).Scan(&originalURL)
	return originalURL, err
}

func (r *URLRepository) Ping(ctx context.Context) error {
	conn, err := r.db.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return conn.Ping(ctx)
}
