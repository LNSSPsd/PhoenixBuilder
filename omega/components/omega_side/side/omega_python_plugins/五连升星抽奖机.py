# 插件 开
# 参考:https://www.bilibili.com/video/BV1B34y1f7xN

import os,sys,time
from threading import Thread
from omega_side.python3_omega_sync.bootstrap import install_lib
install_lib(lib_name="numpy",lib_install_name="numpy")
install_lib(lib_name="rich",lib_install_name="rich")
from omega_side.python3_omega_sync import API
from omega_side.python3_omega_sync import frame as omega
from omega_side.python3_omega_sync.protocol import *
from omega_python_plugins import 抽卡
from rich import print as rprint
import random
def plugin_main(api:API):
    方式 = "检测物品" # 检测物品/检测计分板
    ID = 304 #钻石
    词典 = "minecraft:diamond" # https://www.mcmod.cn/class/1.html可以在这里找词典
    数量 = 5
    def on_menu_main(packet):
        try:
            player=packet["SourceName"]
            if 方式 in "检测物品":
                response=api.do_get_uqholder(cb=None)
                response = response.PlayersByEntityID
                for i in response:
                    data = response[i]
                    if data["Username"] == player:
                        Entity = data["Entity"]
                        Slots = Entity["Slots"]
                        if Slots == {}:
                            api.do_set_player_title(player,"§4无法检测到物品,请将物品放入物品栏!",cb=None)
                        elif len(Slots) !=0:
                            lastslot=Entity["LastPacketSlot"]
                            NewItem = Slots[str(lastslot)]["NewItem"]
                            Stack = NewItem["Stack"]
                            手持ID = Stack["NetworkID"]
                            手持数量 = Stack["Count"]
                            if 手持ID == ID and 手持数量 >= 数量-1:
                                api.do_set_player_title(player,"§a请勿修改物品栏,为了能够显示在您的前方,请勿转动视角!",cb=None)
                                time.sleep(1)
                                扣除状态 = api.do_send_ws_cmd("/clear @a[name="+player+"] "+词典+" 0 "+str(数量),cb=None)
                                扣除状态 = dict(扣除状态.result.OutputMessages[0])
                                扣除状态 = 扣除状态["Success"]
                                if 扣除状态 == True:
                                    api.do_set_player_title(player,"§a成功扣除!",cb=None)
                                    # rprint(Entity)
                                    api.do_send_ws_cmd("/tp @s "+player,cb=None)
                                    api.do_set_player_title(player,"§a即将清空前方方块,请确保无任何遮挡",cb=None)
                                    time.sleep(2)
                                    api.do_set_player_title(player,"§a所需空间 15X15X15 ",cb=None)
                                    time.sleep(1)
                                    timelist = ["10","9","8","7","6","5","4","3","2","1"]
                                    for i in timelist:
                                        api.do_set_player_title(player,"§a"+i,cb=None)
                                        time.sleep(1)
                                        api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)
                                        api.do_set_player_title(player,"§a清空中!",cb=None)
                                        if i=="1":
                                            time.sleep(1)
                                            api.do_send_ws_cmd("/playsound random.levelup "+player,cb=None)
                                    api.do_send_ws_cmd("/title @a actionbar awa",cb=None)
                                    api.do_send_ws_cmd("/fill ^0 ^0 ^0 ~+15 ~+15 ~+7 air 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/fill ^0 ^0 ^0 ~+15 ~+15 ~-7 air 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/tp "+player+" @s",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~2 ~1 ~0",cb=None)
                                    api.do_send_ws_cmd("/fill ~+2 ~+7 ~3 ~+2 ~+7 ~3 minecraft:concretepowder 7 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.anvil_land "+player,cb=None)
                                    api.do_send_ws_cmd("/fill ~+2 ~+1 ~3 ~+2 ~+1 ~3 minecraft:glass 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~+2 ~+1 ~3",cb=None)

                                    api.do_send_ws_cmd("/fill ~+4 ~+7 ~2 ~4 ~+7 ~2 minecraft:concretepowder 7 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.anvil_land "+player,cb=None)
                                    api.do_send_ws_cmd("/fill ~+4 ~+1 ~2 ~+4 ~+1 ~2 minecraft:glass 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~+4 ~+1 ~2",cb=None)
                                    
                                    api.do_send_ws_cmd("/fill ~+3 ~+7 ~0 ~3 ~+7 ~0 minecraft:concretepowder 7 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.anvil_land "+player,cb=None)
                                    api.do_send_ws_cmd("/fill ~+3 ~+1 ~0 ~+3 ~+1 ~0 minecraft:glass 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~+3 ~+1 ~0",cb=None)

                                    api.do_send_ws_cmd("/fill ~+4 ~+7 ~-2 ~4 ~+7 ~-2 minecraft:concretepowder 7 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.anvil_land "+player,cb=None)
                                    api.do_send_ws_cmd("/fill ~+4 ~+1 ~-2 ~+4 ~+1 ~-2 minecraft:glass 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~+4 ~+1 ~-2",cb=None)

                                    api.do_send_ws_cmd("/fill ~+2 ~+7 ~-3 ~2 ~+7 ~-3 minecraft:concretepowder 7 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.anvil_land "+player,cb=None)
                                    api.do_send_ws_cmd("/fill ~+2 ~+1 ~-3 ~+2 ~+1 ~-3 minecraft:glass 0 destroy",cb=None)
                                    api.do_send_ws_cmd("/particle minecraft:knockback_roar_particle ~+2 ~+1 ~-3",cb=None)
                                    结果 = 抽卡.cj();获得 = 结果[0];星级 = 结果[1]#;rprint(获得);rprint(星级)
                                    time.sleep(1)
                                    def titlesub ():
                                        api.do_send_ws_cmd("/playsound mob.enderdragon.growl "+player,cb=None)
                                        颜色列表 = ["§2","§3","§4","§5","§6","§9","§0"]
                                        # rprint(random.choice(颜色列表))
                                        api.do_set_player_title(player," "+random.choice(颜色列表)+"✡",cb=None)
                                        api.do_send_ws_cmd("/title "+player+" subtitle ✩",cb=None)
                                        time.sleep(0.4)
                                        api.do_send_ws_cmd("/title "+player+" subtitle ✩✩",cb=None)
                                        time.sleep(0.4)
                                        api.do_send_ws_cmd("/title "+player+" subtitle ✩✩✩",cb=None)
                                        time.sleep(0.4)
                                        api.do_send_ws_cmd("/title "+player+" subtitle ✩✩✩✩",cb=None)
                                        time.sleep(0.4)
                                        api.do_send_ws_cmd("/title "+player+" subtitle ✩✩✩✩✩",cb=None)
                                    titlesub()
                                    api.do_send_ws_cmd("/fill ~+1 ~+1 ~-3 ~+1 ~+1 ~-3 minecraft:frame 1 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)  
                                    titlesub()
                                    api.do_send_ws_cmd("/fill ~+3 ~+1 ~-2 ~+3 ~+1 ~-2 minecraft:frame 1 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)  
                                    titlesub()
                                    api.do_send_ws_cmd("/fill ~+2 ~+1 ~0 ~+2 ~+1 ~0 minecraft:frame 1 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)  
                                    titlesub()
                                    api.do_send_ws_cmd("/fill ~+3 ~+1 ~2 ~+3 ~+1 ~2 minecraft:frame 1 destroy",cb=None)
                                    time.sleep(0.8)
                                    api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)  
                                    titlesub()
                                    api.do_send_ws_cmd("/fill ~+1 ~+1 ~3 ~+1 ~+1 ~3 minecraft:frame 1 destroy",cb=None)
                                    time.sleep(0.8)  
                                    api.do_send_ws_cmd("/playsound random.orb "+player,cb=None)  
                                    api.do_send_ws_cmd("/say  §a玩家"+player+"抽奖获得了物品"+str(获得),cb=None)  
                                    IDDICT = 抽卡.getid()
                                    for i in 获得:
                                        api.do_send_ws_cmd("/give "+player+" "+IDDICT[i]+" 1",cb=None)  
                                elif 扣除状态 == False:
                                    api.do_set_player_title(player,"§4扣除失败!请检查物品是否在物品栏!",cb=None)
                            else:
                                api.do_set_player_title(player,"§4物品栏没有指定物品或数量小于:"+str(数量),cb=None)
        except Exception as err:
            api.do_send_ws_cmd("/say §4ERROR >> "+str(err),cb=None)
            if 扣除状态 == True:
                api.do_send_ws_cmd("/give "+player+" "+词典+" "+str(数量),cb=None)
                api.do_set_player_title(player,"§a已经退回物品",cb=None)
    def getIDAddPlayer_IDTextr_pakcet(packet):
        if packet["Message"] == "抽奖":
            Thread(target=on_menu_main,args=[packet]).start()
    api.listen_mc_packet(pkt_type="IDText",cb=None,on_new_packet_cb=getIDAddPlayer_IDTextr_pakcet)

omega.add_plugin(plugin=plugin_main)
# 17695905
# omega.run(addr=None)