# Discord Bot

Một Discord bot đơn giản được viết bằng Golang với các chức năng cơ bản như phản hồi tin nhắn, tìm kiếm Google, tích hợp AI Gemini, và quản lý hệ thống thông báo.

## Dành cho Developers

### Bước 1: Chuẩn bị
- Clone repository
- Đăng ký tài khoản [Discord Developer](https://discord.com/developers/applications) và tạo một bot mới
- Cài đặt [Golang](https://golang.org/dl/) trên máy tính của bạn (nếu cần)
- Cài đặt Docker
- Lấy token của botChat Discord
- Tạo tài khoản và lấy API key từ [SerpAPI](https://serpapi.com/) cho chức năng tìm kiếm Google
- Tạo tài khoản và lấy API key từ [Google AI Studio](https://aistudio.google.com/) cho chức năng AI Gemini

### Bước 2: Cài đặt thư viện (Nếu Cần thiết)
Các thư viện mặc định đã có sãn nếu chưa có hãy cài đặt các thư viện cần thiết bằng lệnh:
```bash
go get github.com/bwmarrin/discordgo
go get github.com/google/generative-ai-go
go get github.com/tidwall/gjson
go get github.com/go-redis/redis/v8
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/joho/godotenv
```

### Bước 3: Đăng ký bot với Discord Server
- Truy cập [Discord Developer Portal](https://discord.com/developers/applications)
- Chọn ứng dụng bot của bạn
- Vào mục "Bot" để lấy token (thay thế "Token" trong `main.go`)
- Vào mục "OAuth2" > "URL Generator":
  - Chọn scopes: bot
  - Chọn bot permissions: Send Messages, Read Message History, Manage Messages, Use Slash Commands
  - Sử dụng URL được tạo để thêm bot vào server Discord

### Bước 4: Cấu hình bot
- Tạo và cấu hình lại file `.env`, dựa vào file `.env.example`.

### Bước 5: Chạy Docker
Sử dụng Terminal để chạy docker và khởi tạo PostgreSQL và Redis:
```bash
docker-compose up -d
```

## Dành cho Users

### Cài đặt
1. Yêu cầu quản trị viên Discord thêm bot vào server của bạn
2. Đảm bảo bot có quyền đọc và gửi tin nhắn trong kênh bạn muốn sử dụng

### Các lệnh cơ bản
| Lệnh | Mô tả |
|---------------------|-----------------------------------------------------|
| `!ping` | Bot sẽ phản hồi "Pong!" |
| `!hello` | Bot sẽ chào bạn với tên người dùng của bạn |
| `!search [từ khóa]` | Bot sẽ tìm kiếm Google và trả về kết quả tìm kiếm |

### Lệnh tích hợp AI
| Lệnh | Mô tả |
|---------------------------|-------------------------------------------------------------|
| `!ask [câu hỏi]` | Gửi câu hỏi cho AI Gemini và nhận câu trả lời |
| `!ask [câu hỏi]` + hình | Gửi câu hỏi và hình ảnh để AI Gemini phân tích và trả lời |

### Lệnh đăng ký và thông báo
| Lệnh | Mô tả |
|----------------------------|--------------------------------------------------------|
| `!subscribe` | Đăng ký nhận thông báo từ bot |
| `!unsubscribe` | Hủy đăng ký nhận thông báo từ bot |
| `!notify [nội dung]` | Gửi thông báo đến tất cả người đã đăng ký (chỉ Admin) |
| `!dm [nội dung]` | Gửi tin nhắn trực tiếp đến tất cả người đăng ký (chỉ Admin) |

### Lệnh quản trị viên
| Lệnh                | Mô tả |
|---------------------|---------------------------------------------------------------|
| `!setadmin @[User]` | Cấp quyền Admin cho người dùng |
| `!listsubscribers`  | Hiển thị danh sách người dùng đã đăng ký nhận thông báo |

### Quyền Admin
Người dùng có quyền Admin có thể:
- Gửi thông báo đến tất cả người đăng ký bằng lệnh `!notify`
- Gửi tin nhắn trực tiếp đến tất cả người đăng ký bằng lệnh `!dm`
- Cấp quyền Admin cho người dùng khác bằng lệnh `!setadmin`
- Xem danh sách người đăng ký bằng lệnh `!listsubscribers`

### Thiết lập Admin đầu tiên
- Khi mới cài đặt bot, chưa có người dùng nào được cấp quyền Admin
- **Người đầu tiên** sử dụng lệnh `!setadmin` sẽ tự động trở thành Admin đầu tiên của hệ thống
- Sau đó, chỉ người có quyền Admin mới có thể cấp quyền Admin cho người khác

### Hệ thống đăng ký
- Người dùng có thể đăng ký để nhận thông báo từ bot bằng lệnh `!subscribe`
- Khi đã đăng ký, người dùng sẽ nhận các thông báo từ Admin qua kênh chung hoặc tin nhắn riêng
- Người dùng có thể hủy đăng ký bất cứ lúc nào bằng lệnh `!unsubscribe`

### Ví dụ sử dụng
- Để nhận lời chào: `!hello`
- Để tìm kiếm thông tin về Việt Nam: `!search Việt Nam`
- Để hỏi AI Gemini về Việt Nam: `!ask Việt Nam là quốc gia như thế nào?`
- Để đăng ký nhận thông báo: `!subscribe`
- Để hủy đăng ký: `!unsubscribe`
- Để trở thành Admin đầu tiên (khi chưa có Admin nào): `!setadmin`
- Để gửi thông báo (nếu bạn là Admin): `!notify Cập nhật mới vừa được phát hành!`
- Để gửi tin nhắn riêng (nếu bạn là Admin): `!dm Chào mừng bạn đến với cộng đồng của chúng tôi!`
- Để cấp quyền Admin (nếu bạn là Admin): `!setadmin @{user} true`
- Để xóa quyền Admin (nếu bạn là Admin): `!setadmin @{user} false`

### Chức năng phân tích hình ảnh
Bot có khả năng phân tích hình ảnh kết hợp với văn bản:
1. Sử dụng lệnh `!ask` cùng với câu hỏi của bạn
2. Đính kèm một hoặc nhiều hình ảnh vào tin nhắn
3. Bot sẽ sử dụng AI Gemini để phân tích và trả lời dựa trên cả văn bản và hình ảnh

### Xử lý sự cố
- Nếu gặp vấn đề với lệnh tìm kiếm, hãy thông báo cho quản trị viên để kiểm tra cấu hình API
- Nếu bot không phản hồi, hãy kiểm tra xem bot có trực tuyến hay không
- Đảm bảo bot có đủ quyền để đọc và gửi tin nhắn trong kênh bạn đang sử dụng

