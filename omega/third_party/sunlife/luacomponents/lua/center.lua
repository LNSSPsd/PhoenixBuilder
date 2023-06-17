local communicator = require("communicator")

-- 通过 communicator 向 Go 发送消息
communicator.send_message("Hello, Go!")

-- 从 Go 获取消息
local message = communicator.get_message()
print("Received a message from Go: " .. message)