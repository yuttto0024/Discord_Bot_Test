package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv" // .envファイルを読み込むために追加
)

func main() {
    // .envファイルを読み込み
    err := godotenv.Load()
    if err != nil {
        log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
    }

    // 環境変数からDiscord Botのトークン取得
    token := os.Getenv("DISCORDTOKEN")
    if token == "" {
        log.Fatal("Discordトークンが設定されていません。環境変数DISCORDTOKENを設定してください。")
    }

    // 新しいDiscordセッションを作成
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Fatalf("Error creating Discord session: %v", err)
    }

    // Botを起動し、Discordのサーバーに接続
    err = dg.Open()
    if err != nil {
        log.Fatalf("Error opening connection: %v", err)
    }
    fmt.Println("Bot is now running. Press CTRL+C to exit")

    // 定期的にメッセージを送信するためのタイマー
    ticker := time.NewTicker(1 * time.Minute) // 1分ごとにメッセージを送信
    defer ticker.Stop()

    // メッセージを送信するテキストチャンネルのID
    channelID := "1278707009549631552"

    // 1分ごとに「Good job!!」を送信
    for {
        <-ticker.C
        _, err := dg.ChannelMessageSend(channelID, "Good job!!")
        if err != nil {
            log.Printf("Error sending message: %v", err)
        }
    }
}
