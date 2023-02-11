import os,sys,time
from threading import Thread
from omega_side.python3_omega_sync.bootstrap import install_lib
install_lib(lib_name="numpy",lib_install_name="numpy")
from omega_side.python3_omega_sync import API
from omega_side.python3_omega_sync import frame as omega
from omega_side.python3_omega_sync.protocol import *
import numpy as np
import random
import toml
def cj():
    config = toml.load(os.getcwd ()+"\抽奖机_config.toml")
    五星=config["KC1"]["list_5"]
    四星=config["KC1"]["list_4"]
    三星=config["KC1"]["list_3"]
    次数 = 0
    np.random.seed(0)
    概率 = np.array(config["KC1"]["GL"]) # 实际概率=数值x100 加起来必须等于1！！！
    结果list = []
    结果星级 = []
    while 次数<5:
        结果= np.random.choice([random.choice(五星),random.choice(四星), random.choice(三星)], p=概率.ravel())
        次数=次数+1
        结果list.append(结果)
        for i in 五星:
            if i == 结果:
                结果星级.append("5")
        for i in 四星:
            if i == 结果:
                结果星级.append("4")
        for i in 三星:
            if i == 结果:
                结果星级.append("3")
    else:
        return 结果list,结果星级
def getid ():
    IDDICT = {"钻石剑":"minecraft:diamond_sword",
              "钻石块":"minecraft:diamond_block",
              "不死图腾":"minecraft:totem_of_undying",
              "钻石":"minecraft:diamond",
              "铁镐":"minecraft:iron_pickaxe",
              "木棍":"minecraft:stick",
              "树苗":"minecraft:sapling"
    }
    return IDDICT