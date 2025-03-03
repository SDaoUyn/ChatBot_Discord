# Discord Bot

Một Discord bot đơn giản được viết bằng Golang với các chức năng cơ bản như phản hồi tin nhắn, tìm kiếm Google và quản lý trạng thái hoạt động.

## Dành cho Developers

### Bước 1: Chuẩn bị
- Đăng ký tài khoản [Discord Developer](https://discord.com/developers/applications) và tạo một bot mới
- Lấy token của bot
- Cài đặt [Golang](https://golang.org/dl/) trên máy tính của bạn

### Bước 2: Cài đặt thư viện nếu chưa có
Cài đặt các thư viện cần thiết bằng các lệnh:
```bash
go get github.com/bwmarrin/discordgo
go get github.com/google/generative-ai-go
go get github.com/tidwall/gjson
```

### Bước 3: Đăng ký bot với server Discord
- Vào trang [Discord Developer Portal](https://discord.com/developers/applications)
- Chọn ứng dụng bot của bạn
- Vào mục "Bot" để lấy token (thay thế "TOKEN_BOT_CỦA_BẠN" trong code)
- Vào mục "OAuth2" > "URL Generator":
    - Chọn scopes: bot
    - Chọn bot permissions: Send Messages, Read Message History
    - Sử dụng URL được tạo để thêm bot vào server Discord

### Bước 4: Cấu hình bot
- Mở file `main.go` và thay thế `"Token"` bằng token bot Discord của bạn
- Mở file `common.SearchGoogle` và thay thế `"SerpAPIKey"` bằng API key của bạn từ [SerpAPI](https://serpapi.com/)
- Mở file `common.AskGemini` và thay thế `"GeminiKey"` bằng API key của bạn từ [GeminiAPI](https://aistudio.google.com/)

### Bước 5: Chạy bot
Chạy mã nguồn bằng lệnh:
```bash
go run main.go
```

### Chỉnh sửa và mở rộng
- Tất cả logic xử lý lệnh nằm trong hàm `messageCreate`
- Thêm lệnh mới bằng cách thêm các khối `if strings.HasPrefix()` mới
- Khi sửa đổi, hãy đảm bảo tài nguyên (như kết nối http) được đóng đúng cách

## Dành cho Users

### Cài đặt
1. Yêu cầu quản trị viên Discord thêm bot vào server của bạn
2. Đảm bảo bot có quyền đọc và gửi tin nhắn trong kênh bạn muốn sử dụng

### Các lệnh có sẵn

| Lệnh                | Mô tả                                               |
|---------------------|-----------------------------------------------------|
| `!ping`             | Bot sẽ phản hồi "Pong!"                             |
| `!hello`            | Bot sẽ chào bạn với tên người dùng của bạn          |
| `!search [từ khóa]` | Bot sẽ tìm kiếm Google và trả về 3 kết quả đầu tiên |
| `!ask [từ khóa]`    | Bot sẽ gọi AI gemini để trả lời                     |

### Ví dụ sử dụng
- Để nhận lời chào: `!hello`
- Để tìm kiếm thông tin về Việt Nam: `!search Việt Nam`
- Để hỏi AI gemini về Việt Nam: `!ask Việt Nam`

### Xử lý sự cố
- Nếu gặp vấn đề với lệnh tìm kiếm, hãy thông báo cho quản trị viên để kiểm tra cấu hình API