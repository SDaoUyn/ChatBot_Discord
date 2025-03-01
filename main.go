package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string = "Your-token"

func main() {
	// Tạo một session Discord mới
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Lỗi khi tạo session Discord:", err)
		return
	}

	// Đăng ký hàm xử lý sự kiện messageCreate
	dg.AddHandler(messageCreate)

	// Đăng ký intent để nhận tin nhắn
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Mở kết nối tới Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Lỗi khi kết nối:", err)
		return
	}

	// Thông báo bot đã hoạt động
	fmt.Println("Bot đang chạy. Nhấn CTRL-C để thoát.")

	// Chờ tín hiệu CTRL-C hoặc tương tự để tắt bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// Hàm xử lý khi có tin nhắn mới
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bỏ qua tin nhắn từ chính bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!ping") {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.HasPrefix(m.Content, "!hello") {
		s.ChannelMessageSend(m.ChannelID, "Xin chào, "+m.Author.Username+"!")
	}

}
