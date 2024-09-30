package main

import (
	"fmt"
	"log"
	"time"

	// 環境変数など、OSの機能を使うためのライブラリ
	"os"
	// discordgo:DiscordのAPIにアクセスするためのライブラリ
	"github.com/bwmarrin/discordgo"
	// godotenv:.envファイルから環境変数を読み込むためのライブラリ
	"github.com/joho/godotenv"
)

// ユーザーの入室時間を管理するためのマップ
// userJoinTimesという名前の変数を宣言
// これはユーザーID（文字列型）をキー、
// ユーザーがボイスチャンネルに参加した時刻（time.Time型）を値とするマップを作成
// このマップを、userJoinTimesに割り当てた
var userJoinTimes = make(map[string]time.Time)

func main() {
	// godotenv.Load():.envファイルを読み込む関数
	// 環境変数をプログラムに利用できるようにする
	// 成功するとnilを返し、失敗するとエラーが返る
	// err変数にエラーの値を格納
	err := godotenv.Load()
	
	// エラーが発生した場合の処理
	// if err != nil: エラーがnilでない場合に実行される
	// log.Fatalf()を使ってエラーメッセージを出力し、プログラムを終了する
	if err != nil {
		log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
	}

	// 環境変数からDiscordBotのトークンを取得
	token := os.Getenv("DISCORDTOKEN")
	if token == "" {
		log.Fatal("Discordトークンが設定されていません。環境変数DISCORDTOKENを設定してください。")
	}

	// discordgo.New():DiscordAPIに接続するためのセクションを作成する
	// dgに作成したセクションの結果を格納
	// このdg(セクション)を通じ、Botでメッセージを送ったり、イベントに反応できる
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Botを起動し、Discordのサーバーに接続
	// WebSocketを使用し、サーバー上のイベントをリアルタイムで受け取る
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit")

	// ボイスチャンネルの入退出を監視するのイベントハンドラの登録
	dg.AddHandler(voiceStateUpdate)

	// プログラムが終了まで無限待機、外部イベントずっと監視状態
	// select{}がないと、main()が終了し、プログラムも終了する
	select {}
}

    // ボイスチャンネルの状態が更新されたときに呼ばれるイベントハンドラ
    // 関数内でポインタsを使い、discord.Sessionの値にアクセス
    // 関数内でポインタvsuを使い、discordgo.VoiceStateUpdateの値にアクセス
    // セッション、入退室が保存されたメモリのアドレスで、関数内のデータを操作できる
    // == nil:「ポインタが有効なデータを持っているか」確認するために使う
    // vsuがnilの時、入退室イベントが発生してないと判断し、関数を終了する
func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
    if vsu == nil {
        log.Println("VoiceStateUpdate event is nil")
        return
    }

    // ユーザーがボイスチャンネルに参加した場合、時間を記録
	// vsu.ChannelID != "": ユーザーがボイスチャンネルに参加した場合に真となる
	// vsu.BeforeUpdate == nil: ボイスチャンネル参加する前の状態が存在しない = ユーザーが新しく参加したことを意味する
    if vsu.ChannelID != "" && vsu.BeforeUpdate == nil {
        userJoinTimes[vsu.UserID] = time.Now()
        log.Printf("User %s has joined the voice channel at %v", vsu.UserID, userJoinTimes[vsu.UserID])
        return
	}

    // ユーザーの退出を確認
	// vsu.BeforeUpdate != nil: ユーザーが以前にボイスチャンネルに参加していたことを意味する
	// vsu.ChannelID == "": ("")である場合、ユーザーはボイスチャンネルを退出したことを意味する
    if vsu.BeforeUpdate != nil && vsu.ChannelID == "" {
        userID := vsu.UserID
        channelID := "1278707009549631552" // メッセージを送信するチャンネルのID

        // 滞在時間を計算
		// userJoinTimesマップから、ユーザーの参加時刻を取得。
        // ユーザーIDをキーに、参加時刻(time.Time 型)を追跡
        joinTime, ok := userJoinTimes[userID]

		// if ok: ok変数は、joinTime, ok := userJoinTimes[userID] で取得した値が、マップuserJoinTimesに存在するか確認
		// trueの場合、指定したuserIDの参加時刻がマップに存在するため実行
		// falseの場合、このブロックは実行されない
        if ok {

			// joinTime: 参加時刻を表す
			// time.Since(joinTime): 現在時刻との差を計算 = 参加時刻から現在時刻までの経過時間を計算
            duration := time.Since(joinTime) 

            // メッセージをフォーマットで作成
            durationMessage := fmt.Sprintf("<@%s> Good job!! You stayed for %v.", userID, duration)

			// 作成したメッセージをDiscordの特定のチャンネルに送信
            _, err := s.ChannelMessageSend(channelID, durationMessage)

			// メッセージ送信が失敗した場合のエラーハンドリング
            if err != nil {
                log.Printf("Error sending message: %v", err)
            }

            // ユーザーの入室時刻の記録を削除
            delete(userJoinTimes, userID)
        } else {
            log.Printf("No join time found for user %s", userID)
        }
    }
}