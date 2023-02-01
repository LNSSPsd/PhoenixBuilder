# 插件: 开

import os,sys
from omega_side.python3_omega_sync.bootstrap import install_lib
install_lib("flask");install_lib("requests")
from omega_side.python3_omega_sync import API
from omega_side.python3_omega_sync import frame as omega
from omega_side.python3_omega_sync.protocol import *
import json,requests,datetime
from flask import Flask,render_template

data = "./omega_python_plugins/OmgSide插件-网页管理系统DATA/"
def plugin_main(api:API):
    def on_player_login(player:PlayerInfo):
        if not os.path.exists(data+datetime.datetime.now().strftime('%Y-%m-%d')+".log"):   
            file = open(data+datetime.datetime.now().strftime('%Y-%m-%d')+".log",'w')
            file.close()
        print("["+datetime.datetime.now().strftime("%Y-%m-%d-%H:%M:%S")+"]"+"玩家进入"+">"+player.name+"")
        with open(data+datetime.datetime.now().strftime('%Y-%m-%d')+".log", 'a') as f:   
            f.write(player.name+"\n") 
    def plugin_login():
        api.listen_player_login(cb=None,on_player_login_cb=on_player_login)
    plugin_login()
omega.add_plugin(plugin=plugin_main)