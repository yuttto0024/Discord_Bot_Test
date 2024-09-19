package main

import (
	"fmt"
	"log"

	// 環境変数などのOSの機能を使うためのライブラリ
	"os"
	// discordgo:DiscordのAPIにアクセスするためのライブラリ
	"github.com/bwmarrin/discordgo"
	// godotenv: .envファイルから環境変数を読み込むためのライブラリ
	"github.com/joho/godotenv"
)

func main() {
	// godotenv.Load()は.envファイルを読み込む関数
	// ファイル内の環境変数を読み込み、プログラムに利用できるようにする
	// 成功するとnilを返し、失敗するとエラー→エラーの詳細がerrに入る
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

	// 環境変数からDiscord Botのトークン取得
	// 上記と似た仕組み
	token := os.Getenv("DISCORDTOKEN")
	if token == "" {
		log.Fatal("Discordトークンが設定されていません。環境変数DISCORDTOKENを設定してください。")
	}

	// 新しいDiscordセッションを作成
	// discordgo.New関数でDiscordAPIに接続するためのセクション作成
	// dgに作成したセクションの結果を格納
	// このdg(セクション)を通じ、Botでメッセージを送ったり、イベントに反応できる
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Botを起動し、Discordのサーバーに接続
	// WebSocketを使用し、BotがDiscordのAPIに接続し、サーバー上のイベントをリアルタイムで受け取れるようにしている
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit")

	// ボイスチャンネルの入退出をリッスンするためのイベントハンドラを登録
	dg.AddHandler(voiceStateUpdate)

	// プログラムが終了するまで無限待機、外部イベントずっと監視状態
	// select{}がないと、main()が終了し、プログラムも終了する
	select {}
}

    // ボイスチャンネルの状態が更新されたときに呼ばれるイベントハンドラ
	// 関数内でポインタ s を使って discord.Session の値にアクセス
	// 関数内でポインタ vsu を使って discordgo.VoiceStateUpdate の値にアクセス
	// それぞれセッション、入退室が保存されたメモリのアドレスで、関数内でデータを操作できる
	// == nil は「ポインタが有効なデータを持っているかどうか」を確認するために使う
	// vsu が nil の時、入退室イベントが発生していないと判断、関数を終了する
func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu == nil {
		log.Println("VoiceStateUpdate event is nil")
		return
	}

	// ユーザーがボイスチャンネルから退出したときの処理
	// vsu.BeforeUpdateは、discordgo.VoiceStateUpdateを参照している
	// vsu.UserIDは、vsuの中のUserIDを参照している
	// vsu.BeforeUpdate != nilは、入退出の状態が存在 = ユーザーの退出を表す
	// vsu.ChannelID == "" は、ユーザーがチャンネルから抜けたことを表す
	// sを使って、メッセージ送信、ユーザー情報の取得、イベント監視などを行う
	if vsu.BeforeUpdate != nil && vsu.ChannelID == "" {
		userID := vsu.UserID
		channelID := "1278707009549631552" // メッセージを送信するチャンネルのID

		// 退出したユーザーの情報をログに出力
		log.Printf("User %s has left the voice channel", userID)

		// メンションを飛ばすためのメッセージ作成
		// @sは文字列のフォーマット指定、ポインタは関係ない
		mention := fmt.Sprintf("<@%s> Good job!!", userID)

		// メッセージを送信
		// ChannelMessageSend(channelID, mention):discordgoライブラリに用意されたメソッド
		// 戻り値を持つ関数を２つ(送信したメッセージデータとエラー情報)を使っているが、前者は_で無視
		_, err := s.ChannelMessageSend(channelID, mention)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}



