// main.go
package galaction

import (
	"context"
	"fmt"
	"hsLineOpc/api"
	"log"
	"time"
)

func EntryStartHs(ctx context.Context, hsClient *api.OpcClient) {
	// Инициализация модулей
	hs := NewHS(hsClient)
	packs := NewPackS(hsClient)
	procs := NewProcS(hsClient)
	ss := NewSS(hsClient)

	// Главный управляющий цикл
	if err := controlLoop(ctx, hs, procs, packs, ss); err != nil {
		log.Fatal("Ошибка управляющего процесса:  ", err)
	}
}

func controlLoop(ctx context.Context, hs *HS, procs *ProcS, packs *PackS, ss *SS) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Полный рабочий цикл
			if err := workflow(ctx, hs, procs, packs, ss); err != nil {
				return err
			}
			log.Println("Завершён рабочий цикл, ожидание следующего запуска...")
			time.Sleep(5 * time.Second)
		}
	}
}

func workflow(ctx context.Context, hs *HS, procs *ProcS, packs *PackS, ss *SS) error {
	hs.GrMoveToStart(ctx)
	// Шаг 1: Размещение объекта на карусели
	if err := hs.GrMovePuckToCarousel(ctx); err != nil {
		return fmt.Errorf("Ошибка размещения объекта: %w", err)
	}

	// Шаг 2: Запуск процесса обработки
	if err := procs.Start(ctx); err != nil {
		return fmt.Errorf("Ошибка процесса обработки:  %w", err)
	}

	// Шаг 3: Перемещение к упаковке
	if err := hs.GrMovePuckToPack(ctx); err != nil {
		return fmt.Errorf("Ошибка перемещения к упаковке: %w", err)
	}

	// Шаг 4: Запуск упаковки
	if err := packs.Start(ctx); err != nil {
		return fmt.Errorf("Ошибка упаковки:  %w", err)
	}

	// Шаг 5: Перемещение к сортировке
	if err := hs.GrMovePuckToConveyor(ctx); err != nil {
		return fmt.Errorf("Ошибка перемещения к сортировке: %v", err)
	}

	// Шаг 6: Запуск сортировки
	if err := ss.Start(ctx); err != nil {
		return fmt.Errorf("Ошибка сортировки: %w", err)
	}

	// Шаг 7: Возврат в исходное положение
	return hs.GrMoveToStart(ctx)
}
