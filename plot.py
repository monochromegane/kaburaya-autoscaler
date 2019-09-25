import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('--path-params',     default='out/params.csv')
parser.add_argument('--path-simulation', default='out/simulation.csv')
parser.add_argument('--path-out',        default='out/plot.png')
args = parser.parse_args()

params      = np.loadtxt(args.path_params,     delimiter=",", skiprows=1)
simulations = np.loadtxt(args.path_simulation, delimiter=",", skiprows=1)

fig, axes = plt.subplots(3, 2, sharex=True, figsize=(16,9))
_, _, rho, DT, delay, cdelay = params
fig.suptitle(r'Simulation of AutoScaling [$DT={}$, $\rho={:.2f}$, $Delay(Predicted)={}({})$]'.format(DT, rho, delay, cdelay))

servers = simulations[:, 0]
axes[0,0].plot(servers, label='Servers(Instruction)', color='C1', linestyle='dashed')
delayed_servers = simulations[:, 1]
axes[0,0].plot(delayed_servers, label='Servers(Delayed)', color='C1')
axes[0,0].legend()
axes[0,0].grid()
axes[0,0].set_ylim(0.0)
axes[0,0].set_ylabel('Servers')

waitings = simulations[:, 2]
axes[1,0].plot(waitings, color='C3')
axes[1,0].grid()
axes[1,0].set_ylabel('Waiting')

responseTimes = simulations[:, 3]
axes[2,0].plot(responseTimes, color='C2')
axes[2,0].grid()
axes[2,0].set_ylabel('Response time')

lambdas = simulations[:, 4]
axes[0,1].plot(lambdas, color='C4')
axes[0,1].set_ylim(0.0)
axes[0,1].grid()
axes[0,1].set_ylabel(r'$\lambda$')

mus = simulations[:, 5]
axes[1,1].plot(mus, color='C5')
axes[1,1].set_ylim(0.0)
axes[1,1].grid()
axes[1,1].set_ylabel(r'$\mu$')

axes[2,1].axis('off')

plt.savefig(args.path_out)
