package main

import (
	"fmt"
	"log"
	"time"

	// 環境変数などのOSの機能を使うためのライブラリ
	"os"
	// discordgo:DiscordのAPIにアクセスするためのライブラリ
	"github.com/bwmarrin/discordgo"
	// godotenv: .envファイルから環境変数を読み込むためのライブラリ
	"github.com/joho/godotenv"
)

// ユーザーの入室時間を管理するためのマップ
var userJoinTimes = make(map[string]time.Time)


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


    // ユーザーがボイスチャンネルに参加した場合、時間を記録
	// vsu.ChannelID != "": ユーザーがボイスチャンネルに参加した場合に真になる
	// vsu.BeforeUpdate == nil: ボイスチャンネルに参加する前の状態がないか確認、ユーザーが参加したか
    if vsu.ChannelID != "" && vsu.BeforeUpdate == nil {
        userJoinTimes[vsu.UserID] = time.Now() // ユーザーの参加時間を記録
        log.Printf("User %s has joined the voice channel at %v", vsu.UserID, userJoinTimes[vsu.UserID])
        return
	}


    // ユーザーの退出を確認
	// vsu は discordgo.VoiceStateUpdate 型のポインタ、ボイスチャンネルの状態が更新される情報を待つ
	// vsu.BeforeUpdate != nil:ユーザーが以前にボイスチャンネルに参加していたことを意味する
	// vsu.ChannelID == "":（""）である場合、ユーザーはボイスチャンネルを退出したことを意味する
    if vsu.BeforeUpdate != nil && vsu.ChannelID == "" {
        userID := vsu.UserID
        channelID := "1278707009549631552" // メッセージを送信するチャンネルのID
        // 滞在時間を計算
		// userJoinTimesマップから、ユーザーの参加時刻を取得。ユーザーIDをキーに、参加時刻（time.Time 型）を追跡
        joinTime, ok := userJoinTimes[userID]


		// メッセージの送信
		// if ok: ok変数は、上の行、joinTime, ok := userJoinTimes[userID] で取得した値が、マップuserJoinTimesに存在するか確認
		// trueの場合、指定したuserIDの参加時刻がマップに存在するため実行
		// falseの場合、このブロックは実行されない
        if ok {
			// joinTime:参加時刻を表す
			// time.Since(joinTime) で現在時刻との差を計算
			// よって、参加時刻から現在時刻までの経過時間を計算
            duration := time.Since(joinTime) 

			
            // 滞在時間をメッセージとして送信
            durationMessage := fmt.Sprintf("<@%s> Good job!! You stayed for %v.", userID, duration)
            _, err := s.ChannelMessageSend(channelID, durationMessage) // メッセージを送信
            if err != nil {
                log.Printf("Error sending message: %v", err)
            }
            // 記録を削除
            delete(userJoinTimes, userID) // 参加時間の記録を削除
        } else {
            log.Printf("No join time found for user %s", userID) // 参加時刻が見つからなかった場合
        }
    }
}