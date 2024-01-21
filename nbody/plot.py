import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np
import matplotlib.animation as animation
import os

def main():
	data = []
	with open("particles.dat","r") as f:
		data = f.readlines()
	data=[i.replace('\t', '').replace('\n','').replace('  ',' ').split(" ") for i in data]
	nParticles = int(data[0][0])
	nt = int(data[0][1])+1
	theta = float(data[0][2])
	d = {1:'random',2:'lattice',3:'circle'}
	init_type = d[int(data[0][3])]
	nthreads = int(data[0][4])
	t = {'0':'seq','1':'parallel','2':'work-steal','3':'test'}
	typ = t[data[0][5]]
	data=data[1:]
	data = [[float(x) for x in i] for i in data]
	timedata=[]
	for i in range(nt):	
		timedata.append(data[i*nParticles:(i+1)*nParticles])
	
	median_x = np.median([i[0] for i in data])
	range_x = (np.percentile([i[0] for i in data],95)-np.percentile([i[0] for i in data],5))/2.0
	median_y = np.median([i[1] for i in data])
	range_y = (np.percentile([i[1] for i in data],95)-np.percentile([i[1] for i in data],5))/2.0
	semi_side = min([range_x,range_y])
	
	fig = plt.figure(figsize=(6,6), dpi=200)
	ax = plt.gca()
	
	def animate(i):
		x = [i[0] for i in timedata[i]]
		y = [i[1] for i in timedata[i]]
		vx = [i[2] for i in timedata[i]]
		vy = [i[3] for i in timedata[i]]
		mass = [i[4] for i in timedata[i]]
		# plot in real time
		plt.cla()
		plt.grid()
		plt.title("nParticles = %d theta = %.2f nthreads = %d\n on timestep %d out of %d"%(nParticles,theta,nthreads,i,nt-1))
		plt.scatter(x,y,color='r',s=mass)
		plt.quiver(x,y,vx,vy,scale_units='xy',scale=7)
		ax.set(xlim=(median_x-semi_side,median_x+semi_side), ylim=(median_y-semi_side,median_y+semi_side))

	# Save figure
	try:
		os.mkdir('result_gifs')
	except:
		pass

	print("Saving simulation GIF...please wait...")
	ani = animation.FuncAnimation(fig, animate, repeat=False,
                                    frames=nt, interval=100)
	writer = animation.PillowWriter(fps=10,
	                                metadata=dict(artist='dikshant'),
	                                bitrate=1200)
	# ani.save("result_gifs/BH-simulation-%d-%d-%.2f-%s.gif"%(nParticles,nt-1,theta,init_type), writer=writer)
	ani.save("result_gifs/"+typ+"-BH-simulation-.gif", writer=writer)
	print("Done")
	# plt.show()
	    
	return 0

if __name__== "__main__":
	main()
