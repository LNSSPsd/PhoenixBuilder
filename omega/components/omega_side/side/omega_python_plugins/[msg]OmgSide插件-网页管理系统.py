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
    def on_text_msg(packet):
        if not os.path.exists(data+"msg-"+datetime.datetime.now().strftime('%Y-%m-%d')+".log"):   
            file = open(data+"msg-"+datetime.datetime.now().strftime('%Y-%m-%d')+".log",'w')
            file.close()
        with open(data+"msg-"+datetime.datetime.now().strftime('%Y-%m-%d')+".log", 'a') as f:   
            msg = packet.SourceName+":"+packet.Message
            f.write(str(msg)+"\n") 
    def plugin_login():
        api.listen_mc_packet(pkt_type="IDText",cb=None,on_new_packet_cb=on_text_msg)
    plugin_login()
omega.add_plugin(plugin=plugin_main)