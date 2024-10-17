package errors

// Nhóm lỗi E000: Các lỗi chung

// Nhóm lỗi E001: Các lỗi liên quan tới payment record
const PAYMENT_E001_001 string = "yêu cầu thanh toán không tồn tại"
const PAYMENT_E001_002 string = "yêu cầu thanh toán không ở trạng thái khởi tạo"
const PAYMENT_E001_003 string = "hình thức thanh toán của yêu cầu này không phải chuyển khoản"
const PAYMENT_E001_004 string = "yêu cầu thanh toán không ở trạng thái khách hàng xác nhận đã chuyển khoản"
const PAYMENT_E001_005 string = "phương thức thanh toán chưa được định nghĩa"
const PAYMENT_E001_AUTH_SERVICE_NAME = "auth service"

// Nhóm lỗi E002: Các lỗi liên quan đến kết nối service to service
const PAYMENT_E002_CATEGORY_SERVICE_NAME = "category service"
const PAYMENT_E002_001 string = "có lỗi khi kết nối tới %s. StatusCode %s"
const PAYMENT_E002_002 string = "phương thức thanh toán đã ngưng sử dụng"
const PAYMENT_E002_003 string = "gói thanh toán đã ngưng sử dụng"
const PAYMENT_E002_004 string = "thời gian kích hoạt của gói thanh toán không hợp lệ"
const PAYMENT_E002_005 string = "giá của gói kích hoạt không hợp lệ"
const PAYMENT_E001_006 string = "cần điền đầy đủ thông tin chuyển khoản"

// Nhóm lỗi E003: Các lỗi liên quan đến kết nối VNPAY
const PAYMENT_E003_001 string = "có lỗi khi tạo yêu cầu thanh toán. StatusCode %s"
const PAYMENT_E003_002 string = "có lỗi khi lấy thông tin thanh toán từ VNPay."
const PAYMENT_E003_003 string = "lỗi khởi tạo giao dịch từ VNPay."
const PAYMENT_E003_004 string = "dữ liệu IPN không hợp lệ"
const PAYMENT_E003_005 string = "dữ liệu IPN thiếu thông tin bắt buộc"
const PAYMENT_E003_006 string = "checksum IPN không hợp lệ"
const PAYMENT_E003_007 string = "giao dịch không tồn tại"
const PAYMENT_E003_008 string = "lỗi trong quá trình xác minh giao dịch"
const PAYMENT_E003_009 string = "số tiền thanh toán không khớp"
const PAYMENT_E003_010 string = "giao dịch đã được xác nhận"

// Nhóm lỗi E004: Các lỗi theo kich ban SIT cua VNPAY https://sandbox.vnpayment.vn/vnpaygw-sit-testing/order/instruction
const PAYMENT_E004_001 string = "order not found"
const PAYMENT_E004_002 string = "order already confirm"
const PAYMENT_E004_003 string = "invalid amount"
const PAYMENT_E004_004 string = "invalid ip "

// Nhóm lỗi E099: Các lỗi liên quan đến hệ thống
const PAYMENT_E099_001 string = "lỗi xảy ra trong quá trình xử lý #01" // Lỗi đọc/ghi/thực hiện transaction vào mongodb
const PAYMENT_E099_002 string = "chữ ký dữ liệu không hợp lệ #02"      // Lỗi kiểm tra checksum
