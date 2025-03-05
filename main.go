package main

import (
	"ChatBotDiscord/common"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

	// Chức năng AI chat kết hợp văn bản và hình ảnh
	if strings.HasPrefix(m.Content, "!ask ") {
		// Tách phần câu hỏi
		question := strings.TrimPrefix(m.Content, "!ask ")

		// Kiểm tra có hình ảnh không
		if len(m.Attachments) > 0 {
			// Hiển thị bot đang suy nghĩ
			s.ChannelMessageSend(m.ChannelID, "Đang xử lý hình ảnh và câu hỏi...")

			// Mảng lưu đường dẫn file tạm
			tmpImagePaths := []string{}
			defer func() {
				// Xóa các file tạm sau khi sử dụng
				for _, path := range tmpImagePaths {
					os.Remove(path)
				}
			}()

			// Tải và lưu các hình ảnh
			for _, attachment := range m.Attachments {
				// Tải hình ảnh
				resp, err := http.Get(attachment.URL)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Lỗi tải hình ảnh: "+err.Error())
					return
				}
				defer resp.Body.Close()

				// Tạo file tạm
				tmpFile, err := ioutil.TempFile("", "discord-image-*.jpg")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Lỗi tạo file tạm: "+err.Error())
					return
				}

				// Sao chép nội dung hình ảnh
				_, err = io.Copy(tmpFile, resp.Body)
				tmpFile.Close()
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Lỗi lưu hình ảnh: "+err.Error())
					return
				}

				// Lưu đường dẫn file tạm
				tmpImagePaths = append(tmpImagePaths, tmpFile.Name())
			}

			// Nếu có nhiều hình ảnh
			var answer string
			var err error
			if len(tmpImagePaths) > 1 {
				answer, err = common.AskGeminiMultipleImages(question, tmpImagePaths)
			} else {
				// Nếu chỉ có một hình ảnh
				answer, err = common.AskGeminiWithImage(question, tmpImagePaths[0])
			}

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Có lỗi khi xử lý: %s", err.Error()))
				return
			}

			// Gửi kết quả
			s.ChannelMessageSend(m.ChannelID, answer)

		} else {
			// Nếu không có hình ảnh, sử dụng chat text thông thường
			s.ChannelMessageSend(m.ChannelID, "Đang suy nghĩ...")

			answer, err := common.AskGemini(question)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Có lỗi khi truy vấn AI: %s", err.Error()))
				return
			}

			s.ChannelMessageSend(m.ChannelID, answer)
		}
	}

}
