package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	// キャンセル可能な「context.Context」のオブジェクトを作る
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでテスト対象の「run」関数を実行してHTTPサーバーを起動する
	eg.Go(func() error {
		return run(ctx)
	})

	// エンドポイントに対してGETリクエストを送信する
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTPサーバーの戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信する
	cancel()

	// run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
