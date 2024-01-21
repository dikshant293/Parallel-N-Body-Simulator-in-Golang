import matplotlib
matplotlib.use('Agg')

import os,sys,copy
import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from pprint import pprint
argc = len(sys.argv)
jobNumber = sys.argv[1]
nRuns = int(sys.argv[2])
nSizes = int(sys.argv[3])
# size_list = [int(i) for i in sys.argv[4:]]
size_list = sys.argv[4:]
# sizes = [int(i) for i in sys.argv[4:]]
# print(sys.argv,sizes)

absolute_path = os.getcwd()
full_path = os.path.join(absolute_path, "slurm/out/%s.slurm1.stdout"%jobNumber)
f = open(full_path,'r')
lines = [i.strip() for i in f.readlines()]
f.close()
n=len(lines)

nThread_list = [2,4,6,8,12]
dummy_dict = {}
for i in nThread_list:
    dummy_dict[i]=-1.0
    
data = {"s":{},"p":{},"ws":{}}
speedup = {"p":{},"ws":{}}
for i in size_list:
    data["s"][i]=0.0
    speedup['p'][i] = copy.deepcopy(dummy_dict)
    speedup['ws'][i] = copy.deepcopy(dummy_dict)
    
for j in ["p","ws"]:
    for i in size_list:
        data[j][i] = copy.deepcopy(dummy_dict)
# data = {
#             "s":{
#                 "xsmall":0.0,
#                 "small":0.0,
#                 "medium":0.0,
#                 "large":0.0,
#                 "xlarge":0.0
#             },
#             "p":{
#                 "xsmall":{},
#                 "small":{},
#                 "medium":{},
#                 "large":{},
#                 "xlarge":{}
#             },
#             "ws":{
#                 "xsmall":{},
#                 "small":{},
#                 "medium":{},
#                 "large":{},
#                 "xlarge":{}
#             }
#         }
# speedup ={
#             "xsmall":{},
#             "small":{},
#             "medium":{},
#             "large":{},
#             "xlarge":{}
#         }

# for i in size_list:
#     data['p'][i]=copy.deepcopy(dummy_dict)
#     speedup[i]=copy.deepcopy(dummy_dict)
# print(lines,lines[0].split(' ')[0])
d = {'0':"s",'1':"p","2":"ws"}
while len(lines):
    rType = d[lines[0].split(' ')[0]]
    size = lines[0].split(' ')[1]
    nThread = -1
    if rType!='s':
        nThread = int(lines[0].split(' ')[2])
    lines.pop(0)
    vals = []
    while True:
        try:
            vals.append(float(lines[0]))
        except:
            break
        lines.pop(0)
    mean = np.mean(vals).round(4)
    if rType=='s':
        data[rType][size]=mean
    else:
        data[rType][size][nThread]=mean


for size in size_list:
    for n in nThread_list:
        for par in ['p','ws']:
            speedup[par][size][n] = (data['s'][size]/data[par][size][n]).round(4)

# pprint(data)
# pprint(speedup)
try:
    os.mkdir('speedup-graphs')
except:
    pass

plt.figure(figsize=[i*1.25 for i in [6.4,4.8]],dpi=500)
for par in ['p','ws']:
    for k in speedup[par].keys():
        if par=='p':
            plt.plot(nThread_list,[speedup[par][k][i] for i in nThread_list],label="Parallel, nParticles = "+k)
        else:
            plt.plot(nThread_list,[speedup[par][k][i] for i in nThread_list],label="WorkSteal, nParticles = "+k)

        [(plt.text(i,speedup[par][k][i],speedup[par][k][i],fontdict={'size':8})) for i in nThread_list]
plt.legend() 
plt.grid()
plt.title("Speedup Graph")
plt.xlabel("number of threads")
plt.ylabel("speedup")
plt.savefig("speedup-graphs/"+str(jobNumber)+"-graph"+'.png')
plt.close()
