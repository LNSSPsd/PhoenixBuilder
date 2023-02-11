# 插件 开

try:
    from omega_side.python3_omega_sync.bootstrap import install_lib
    install_lib(lib_name="rich",lib_install_name="rich")
    install_lib(lib_name="toml",lib_install_name="toml")
    import os,sys,time,base64 # 万恶的重名
    from omega_python_plugins import Tools_Logging
    import os,sys,time,toml
    from threading import Thread
    from omega_side.python3_omega_sync import API
    from omega_side.python3_omega_sync import frame as omega
    from omega_side.python3_omega_sync.protocol import *
    import socket
    import platform,psutil,datetime
    # 初始化logging
    Tools_Logging.__init__(True)
    lprint = Tools_Logging.logprint
    lprint(" Logging-INFO ","库加载完毕!","Lu")
except Exception as err_msg:
    lprint = Tools_Logging.logprint
    lprint(" IMPORT-ERROR ",str(err_msg),"Hong")
def plugin_main(api:API):
    config = toml.load(os.getcwd ()+"./omega_storage/side/omega_python_plugins/Tools_config.toml")
    def Server ():
        IP_PORT=(config["Main"]["IP"],config["Main"]["PORT"]) 
        Server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        Server.bind((IP_PORT[0],IP_PORT[1]))
        Server.listen(5)
        lprint("Server-INFO","Server waiting for connect","Lu")
        lprint("Server-INFO","IP:"+IP_PORT[0]+" PORT:"+str(IP_PORT[1]),"Lu")
        lprint("Server-INFO","开始进入循环!","Lu")
        while True:
            client_executor, addr = Server.accept()
            lprint("Server-MSG-INFO","来自IP地址:"+addr[0]+":"+str(addr[1])+"的连接","Lu")
            Thread(target=Server_MSG,args=(client_executor,addr)).start()
    def Server_MSG (client_executor,addr):
        MSG = client_executor.recv(1024)
        if not MSG:
            lprint("Server-MSG-INFO","收到来自:"+addr[0]+":"+str(addr[1])+"的信息为空!","Hong")
        else:
            lprint("Server-MSG-INFO","收到来自:"+addr[0]+":"+str(addr[1])+"的信息->"+str(MSG.decode('gbk')),"Lu")
            if str(MSG.decode('gbk')) == "连接成功!":
                time.sleep(0.5)
                lprint("Server-Client-Con","客户端已连接成功!","Lu")
                client_executor.send(bytes(str(True),'utf-8')) 
            if "指令执行" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                Command = str(MSG.decode('gbk'))[7:]
                lprint("Server-Command","执行指令-"+Command,"Lu")
                api.do_send_ws_cmd(Command,cb=None)
            if "命令执行" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                Command = str(MSG.decode('gbk'))[7:]
                lprint("Server-Command","执行命令-"+Command,"Lu")
                Command_return=str(os.popen(Command).read())
                lprint("Server-Command-return","命令返回-> "+Command_return,"Lu")
                client_executor.send(bytes(str(Command_return),'utf-8')) 
            if "主机信息" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                lprint("Server-GET-MSG","客户端请求数据:主机信息","Lu")
                主机信息 = "设备名称:"+os.environ['COMPUTERNAME']+"\n核心型号:"+platform.processor()+"\n服务器端口:"+str(config["Main"]["PORT"])+"\n内存占用:"+str(psutil.virtual_memory().percent)+"%"
                lprint("Server-GET-MSG","返回数据:\n"+主机信息,"Lu")
                client_executor.send(bytes(str(主机信息),'gbk')) 
            if "CPU利用率" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                lprint("Server-GET-MSG","客户端请求数据:CPU利用率","Lu")
                CPU利用率 = psutil.cpu_percent(percpu=True)
                利用率 = 0
                for i in CPU利用率:
                    利用率 = 利用率+i
                else:
                    CPU利用率 = 利用率/4
                lprint("Server-GET-MSG","返回数据:"+str(round(CPU利用率)),"Lu")
                client_executor.send(bytes(str(round(CPU利用率)),'gbk')) 
            if "INFO" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                lprint("Server-GET-MSG","客户端请求数据:日志","Lu")
                log=open("./XingChenTools_data/logs/"+datetime.datetime.now().strftime('%Y-%m-%d')+".log",'r',)
                client_executor.send(bytes(str(log.read()),'gbk')) 
            if "关闭" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                lprint("Server-EXIT","关闭","Lu")
                os.system( 'taskkill /pid ' + str(os.getppid()) + ' /f')
            if "玩家列表" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                playerlist = []
                response=api.do_get_players_list(cb=None)
                for player in response:
                    playerlist.append(player["name"])
                for i in playerlist:
                    time.sleep(1)
                    client_executor.send(bytes(str(i),'gbk')) 
            if "kick" in str(MSG.decode('gbk')):
                time.sleep(0.5)
                lprint("Server-Command-kick","踢出 - "+MSG.decode('gbk')[6:-1],"Lu")
                api.do_send_ws_cmd("/kick "+MSG.decode('gbk')[6:-1],cb=None)
    Thread(target=Server).start()
omega.add_plugin(plugin=plugin_main)
# omega.run(addr=None)