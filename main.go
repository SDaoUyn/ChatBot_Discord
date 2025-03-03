package main

import (
	"ChatBotDiscord/common"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Token Bot Discord của bạn
var Token string = "YOUR_TOKEN"

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

	// Thêm chức năng tìm kiếm Google
	if strings.HasPrefix(m.Content, "!search ") {
		// Lấy từ khóa tìm kiếm từ lệnh
		query := strings.TrimPrefix(m.Content, "!search ")

		if query != "" {
			// Gửi tin nhắn thông báo đang tìm kiếm
			s.ChannelMessageSend(m.ChannelID, "Đang tìm kiếm thông tin về: "+query)

			// Thực hiện tìm kiếm và trả về kết quả
			results, err := common.SearchGoogle(query)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Có lỗi khi tìm kiếm: "+err.Error())
				return
			}

			// Gửi kết quả tìm kiếm
			s.ChannelMessageSend(m.ChannelID, results)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Vui lòng nhập từ khóa để tìm kiếm. Ví dụ: !search Việt Nam")
		}
	}

	// Chức năng AI chat
	if strings.HasPrefix(m.Content, "!ask ") {
		// Lấy câu hỏi từ lệnh
		question := strings.TrimPrefix(m.Content, "!ask ")

		if question != "" {
			// Hiển thị bot đang suy nghĩ
			s.ChannelMessageSend(m.ChannelID, "Đang phản hồi...")

			var answer string
			var err error

			answer, err = common.AskGemini(question)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Có lỗi khi truy vấn AI: %s", err.Error()))
				return
			}

			// Gửi câu trả lời từ AI
			s.ChannelMessageSend(m.ChannelID, answer)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Vui lòng nhập câu hỏi. Ví dụ: !ask Thủ đô của Việt Nam là gì?")
		}
	}
}

// Hàm gọi API Google Gemini
