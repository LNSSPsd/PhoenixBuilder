import os,sys,datetime,socket
def __init__(self) -> None:
    if not os.path.exists("./XingChenTools_data/logs"):
        os.makedirs("./XingChenTools_data/logs")
    if not os.path.exists("./XingChenTools_data/tmps"):
        os.makedirs("./XingChenTools_data/tmps")
    if not os.path.exists("./XingChenTools_data/logs/"+datetime.datetime.now().strftime('%Y-%m-%d')+".log"):   
            file = open("./XingChenTools_data/logs/"+datetime.datetime.now().strftime('%Y-%m-%d')+".log",'w')
            file.close()
def logprint(level:str,msg,color):
    try:
        color_dict = {"Hong":"\033[31m","Lu":"\033[32m","Huang":"\033[33m","Lan":"\033[34m","Mr":"\033[38m"}
        print("["+datetime.datetime.now().strftime("%Y-%m-%d-%H:%M:%S")+"]"+color_dict[color]+" ["+level+"] "+msg+"\033[0m")
        with open("./XingChenTools_data/logs/"+datetime.datetime.now().strftime('%Y-%m-%d')+".log",'a') as log:   
            log.write("["+datetime.datetime.now().strftime("%Y-%m-%d-%H:%M:%S")+"] ["+level+"] "+msg+"\n") 
            log.close()
        exelog("["+datetime.datetime.now().strftime("%Y-%m-%d-%H:%M:%S")+"] ["+level+"] "+msg+"\n")
    except:
        print("["+datetime.datetime.now().strftime("%Y-%m-%d-%H:%M:%S")+"]"+color_dict["Hong"]+" ["+"Logging-ERROR"+"] "+"哇,好像不能输出日志啦!"+"\033[0m")
def exelog(msg):
    try:
        IP=("127.0.0.1",19730)
        exelog=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        exelog.connect(IP)
        exelog.send(bytes(str(msg),'gbk')) 
        exelog.close()
    except:
        pass