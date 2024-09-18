package main

import (
	"fmt"
	"log"
	"os"

	// discordgo:DiscordのAPIにアクセスするためのライブラリ
	"github.com/bwmarrin/discordgo"
	// godotenv: .envファイルから環境変数を読み込むためのライブラリ
	"github.com/joho/godotenv"
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

	// ボイスチャンネルの入退出をリッスンするためのイベントハンドラを登録
	dg.AddHandler(voiceStateUpdate)

	// Botを起動し、Discordのサーバーに接続
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit")

	// プログラムが終了するまで待機
	select {}
}

// ボイスチャンネルの状態が更新されたときに呼ばれるイベントハンドラ
func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu == nil {
		log.Println("VoiceStateUpdate event is nil")
		return
	}

	// ユーザーがボイスチャンネルから退出したときの処理
	// 状態が変わった後のチャンネルIDが空である場合
	if vsu.BeforeUpdate != nil && vsu.ChannelID == "" {
		// 退出したユーザーの情報
		userID := vsu.UserID
		channelID := "1278707009549631552" // メッセージを送信するチャンネルのID

		// 退出したユーザーの情報をログに出力
		log.Printf("User %s has left the voice channel", userID)

		// メンション付きのメッセージを作成
		mention := fmt.Sprintf("<@%s> Good job!!", userID)

		// メッセージを送信
		_, err := s.ChannelMessageSend(channelID, mention)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}



