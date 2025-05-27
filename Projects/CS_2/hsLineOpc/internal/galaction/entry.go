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
	// 初始化各模块
	hs := NewHS(hsClient)
	packs := NewPackS(hsClient)
	procs := NewProcS(hsClient)
	ss := NewSS(hsClient)

	// 主控制流程
	if err := controlLoop(ctx, hs, procs, packs, ss); err != nil {
		log.Fatal("控制流程错误: ", err)
	}
}

func controlLoop(ctx context.Context, hs *HS, procs *ProcS, packs *PackS, ss *SS) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// 完整工作流程
			if err := workflow(ctx, hs, procs, packs, ss); err != nil {
				return err
			}
			log.Println("完成一个工作周期，等待下一次触发...")
			time.Sleep(5 * time.Second)
		}
	}
}

func workflow(ctx context.Context, hs *HS, procs *ProcS, packs *PackS, ss *SS) error {
	hs.GrMoveToStart(ctx)
	// 步骤1: 放置物体到转盘
	if err := hs.GrMovePuckToCarousel(ctx); err != nil {
		return fmt.Errorf("放置物体失败: %w", err)
	}

	// 步骤2: 执行加工流程
	if err := procs.Start(ctx); err != nil {
		return fmt.Errorf("加工流程失败: %w", err)
	}

	// 步骤3: 移动到包装位置
	if err := hs.GrMovePuckToPack(ctx); err != nil {
		return fmt.Errorf("移动至包装失败: %w", err)
	}

	// 步骤4: 执行包装
	if err := packs.Start(ctx); err != nil {
		return fmt.Errorf("包装失败: %w", err)
	}

	// 步骤5: 移动到分拣位置
	if err := hs.GrMovePuckToConveyor(ctx); err != nil {
		return fmt.Errorf("移动至分拣失败: %v", err)
	}

	// 步骤6: 执行分拣
	if err := ss.Start(ctx); err != nil {
		return fmt.Errorf("分拣失败: %w", err)
	}

	// 步骤7: 返回起始位置
	return hs.GrMoveToStart(ctx)
}
