local communicator = require("communicator")

-- 通过 communicator 向 Go 发送消息
communicator.send_message("Hello, Go!")
num = 0
-- 从 Go 获取消息
while true do
  function center()
    num =num +1
  --local message = communicator.get_message()
  communicator.send_message(num)
  --print( message)
end
end
