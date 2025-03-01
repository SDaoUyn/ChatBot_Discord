Dev ->
Bước 1: Chuẩn bị
Đăng ký tài khoản Discord Developer và tạo một bot mới
Lấy token của bot
Cài đặt Golang trên máy tính của bạn

Bước 2: Cài đặt thư viện
Thư viện phổ biến nhất để tạo Discord bot với Golang là DiscordGo. Bạn có thể cài đặt nó bằng lệnh:
-- go get github.com/bwmarrin/discordgo

Bước 3: Đăng ký bot với server Discord
Vào trang Discord Developer Portal
Chọn ứng dụng bot của bạn
Vào mục "Bot" để lấy token (thay thế "TOKEN_BOT_CỦA_BẠN" trong code)
Vào mục "OAuth2" > "URL Generator":
Chọn scopes: bot
Chọn bot permissions: Send Messages, Read Message History
Sử dụng URL được tạo để thêm bot vào server Discord

Bước 4: Chạy bot
Chạy mã nguồn bằng lệnh:
