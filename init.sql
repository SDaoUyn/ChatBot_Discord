-- Tạo bảng subscribers nếu chưa tồn tại
CREATE TABLE IF NOT EXISTS subscribers (
   id SERIAL PRIMARY KEY,
   discord_user_id BIGINT NOT NULL UNIQUE,
   subscribed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   is_active BOOLEAN DEFAULT TRUE,
   is_Admin BOOLEAN DEFAULT FALSE
);

-- Tạo function để tự động cập nhật timestamp
CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Tạo trigger cho sự kiện update
DROP TRIGGER IF EXISTS set_timestamp ON subscribers;
CREATE TRIGGER set_timestamp
  BEFORE UPDATE ON subscribers
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp_column();

-- Tạo trigger cho sự kiện insert
DROP TRIGGER IF EXISTS set_timestamp_before_insert ON subscribers;
CREATE TRIGGER set_timestamp_before_insert
  BEFORE INSERT ON subscribers
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp_column();