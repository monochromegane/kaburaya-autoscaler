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
_, _, rho, DT, delay, odelay, idelay = params
fig.suptitle(r'Simulation of AutoScaling [$DT={}$, $\rho={:.2f}$, $Delay(In,Out)={}({},{})$]'.format(DT, rho, delay, idelay ,odelay))

ideals = simulations[:, 0]
axes[0,0].plot(ideals, label='Ideal', color='C0')
servers = simulations[:, 1]
axes[0,0].plot(servers, label='Instruction', color='C1', linestyle='dashed')
delayed_servers = simulations[:, 2]
axes[0,0].plot(delayed_servers, label='Delayed [Diff from ideal={:+.1f}]'.format(np.sum(delayed_servers)-np.sum(ideals)), color='C1')
axes[0,0].legend()
axes[0,0].grid()
axes[0,0].set_ylim(0.0)
axes[0,0].set_ylabel('Servers')

waitings = simulations[:, 3]
axes[1,0].plot(waitings, color='C3', label=r'Waiting [$Sum={}$]'.format(np.sum(waitings)))
axes[1,0].legend()
axes[1,0].grid()
axes[1,0].set_ylabel('Waiting')

responseTimes = simulations[:, 4]
axes[2,0].plot(responseTimes, color='C2')
axes[2,0].grid()
axes[2,0].set_ylabel('Response time/Unit time')

lambdas = simulations[:, 5]
axes[0,1].plot(lambdas, color='C4')
axes[0,1].set_ylim(0.0)
axes[0,1].grid()
axes[0,1].set_ylabel(r'$\lambda$ (Observed)')

mus = simulations[:, 6]
axes[1,1].plot(mus, color='C5')
axes[1,1].set_ylim(0.0)
axes[1,1].grid()
axes[1,1].set_ylabel(r'$\mu$ (Observed)')

maximum = np.maximum(mus, 1.0/responseTimes)
avg = [np.average(maximum[0:i+1][~np.isnan(maximum[0:i+1])]) for i, m in enumerate(maximum)]
axes[2,1].plot(mus, color='C5', label=r'$\mu$', linestyle='dashed', linewidth=1.0)
axes[2,1].plot(1.0/responseTimes, color='C2', label=r'$1/T_s$', linestyle='dashed', linewidth=1.0)
axes[2,1].plot(avg, color='C6', label=r'average(max($\mu$, $1/ResponseTime$))', linewidth=3.0)
axes[2,1].legend()
axes[2,1].set_ylim(0.0)
axes[2,1].grid()
axes[2,1].set_ylabel(r'$\mu$ (Estimated)')

plt.savefig(args.path_out)
