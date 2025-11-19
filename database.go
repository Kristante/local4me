package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"time"
)

type NoteBD struct {
	ReqID int64
	Date  time.Time
}

func CreateDatabasePool() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Unable to connect to database with error: %v\n", err)
		os.Exit(1)
	}
	return pool
}

func SelectRequest(pool *pgxpool.Pool, reqID int64) int64 {
	var result int64
	query := "SELECT reqid FROM public.requests WHERE reqid = $1"

	err := pool.QueryRow(context.Background(), query, reqID).Scan(&result)
	if err != nil {
		fmt.Println("scan err [selectReques]", err)
		return 0
	}
	return result
}

func AddRequest(pool *pgxpool.Pool, reqID int64) error {
	_, err := pool.Exec(context.Background(), "INSERT into public.requests (reqid) VALUES ($1)", reqID)
	if err != nil {
		return errors.New("ошибка при добавлении заявки в базу данных")
	}
	return nil
}

func CheckNotes(pool *pgxpool.Pool, reqID int64, updateTime string) bool {
	query := "SELECT reqid, date FROM public.notes WHERE reqid = $1 ORDER BY date ASC LIMIT 1"

	var note NoteBD
	err := pool.QueryRow(context.Background(), query, reqID).Scan(&note.ReqID, &note.Date)
	if err != nil {
		fmt.Println("scan err: [CheckNotes]", err)
		return false
	}

	// Парсим updateTime
	t, err := time.Parse(time.RFC3339, updateTime)
	if err != nil {
		fmt.Println("parse time err:", err)
		return false
	}

	// Приводим обе даты к YYYY-MM-DD
	noteDate := note.Date.Format("2006-01-02")
	updateDate := t.Format("2006-01-02")

	return noteDate == updateDate
}

func AddNote(pool *pgxpool.Pool, reqID int64, noteTime string) error {
	t, err := time.Parse(time.RFC3339, noteTime)
	if err != nil {
		return fmt.Errorf("не удалось распарсить время: %v", err)
	}

	// SQL с плейсхолдерами $1, $2
	_, err = pool.Exec(
		context.Background(),
		"INSERT INTO public.notes (reqid, date) VALUES ($1, $2)",
		reqID,
		t.Format("2006-01-02"), // сохраняем только дату
	)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении заметки в базу данных: %v", err)
	}
	return nil
}
