// Playground 是 bug 复现启动脚本。
// 运行逻辑如下：
//
//  1. 启动 docker compose，运行 consul。
//  2. 启动 server。
//  3. 启动 client。
//  4. 等待服务端与客户端交互一段时间。
//  5. 停止 server 和 client。
//  6. 打印 consul 日志。
//  7. 停止 docker compose。
package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	must(cmd(ctx, "docker", "compose", "up", "-d").Run())
	defer func() { must(cmd(context.Background(), "docker", "compose", "down").Run()) }()
	sleep(ctx, 5*time.Second)

	func(ctx context.Context) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// 打包运行而非直接 go run 避免 go run 创建的子进程无法被杀死
		must(cmd(ctx, "go", "build", "-o", "./bin/", "./server", "./client").Run())

		server := cmd(ctx, "./bin/server")
		must(server.Start())

		must(sleep(ctx, 5*time.Second))

		client := cmd(ctx, "./bin/client")
		must(client.Start())

		must(sleep(ctx, 25*time.Second))

		var wg sync.WaitGroup
		wg.Add(1)
		go func() { defer wg.Done(); server.Wait() }()
		wg.Add(1)
		go func() { defer wg.Done(); client.Wait() }()

		cancel()
		wg.Wait()
	}(ctx)

	must(sleep(ctx, 2*time.Second))
	must(cmd(ctx, "docker", "compose", "logs").Run())
}

func cmd(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
