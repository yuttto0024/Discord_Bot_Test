package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/bwmarrin/discordgo"
)

// ボイスチャンネルへの参加時間を記録するためのマップ
var voiceStates = make(map[string]time.Time)

// 監視するボイスチャンネルのID（実際のIDに置き換えてください）
var targetVoiceChannelID = "1278707009549631553"

// メッセージを投稿するテキストチャンネルのID（実際のIDに置き換えてください）
var targetTextChannelID = "1278707009549631552"

func main() {
    // 環境変数からDiscord Botのトークンを取得
    token := os.Getenv("DISCORD_TOKEN")

    // Discordセッションを作成
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("Error creating Discord session,", err)
        return
    }

    // ボイスチャンネルの状態が変わったときに呼ばれるハンドラを追加
    dg.AddHandler(voiceStateUpdate)

    // Discordへの接続を開始
    err = dg.Open()
    if err != nil {
        fmt.Println("Error opening connection,", err)
        return
    }

    fmt.Println("Bot is now running. Press CTRL-C to exit.")

    // プログラムが終了しないように待機
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-stop

    // プログラム終了時にDiscordセッションを閉じる
    dg.Close()
}

// ボイスチャンネルの状態が更新されたときに呼ばれる関数
func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
    userID := vsu.UserID

    // ユーザーがボイスチャンネルに参加したとき
    if vsu.BeforeUpdate == nil && vsu.ChannelID == targetVoiceChannelID {
        // 参加した時間を記録
        voiceStates[userID] = time.Now()
    } else if vsu.BeforeUpdate != nil && vsu.BeforeUpdate.ChannelID == targetVoiceChannelID && vsu.ChannelID != targetVoiceChannelID {
        // ユーザーがボイスチャンネルを退出したとき
        if joinTime, ok := voiceStates[userID]; ok {
            // 滞在時間を計算
            duration := time.Since(joinTime)
            delete(voiceStates, userID)

            // 時間、分、秒を計算
            hours := int(duration.Hours())
            minutes := int(duration.Minutes()) % 60
            seconds := int(duration.Seconds()) % 60

            // 滞在時間をフォーマットしてメッセージを作成
            message := fmt.Sprintf("<@%s> stayed in the voice channel for %02d:%02d:%02d", userID, hours, minutes, seconds)

            // テキストチャンネルにメッセージを送信
            s.ChannelMessageSend(targetTextChannelID, message)
        }
    }
}