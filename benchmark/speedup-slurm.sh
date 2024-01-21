#!/bin/bash
#
#SBATCH --mail-user=superdpsingh123@gmail.com
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/dikshant/progs/MPCS-Courses/Parallel_Programming/project-3-dikshant293/proj3/benchmark
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=59:00

module load golang/1.19

sizes=(10000 100000)
types=(1 2)
nThreads=(2 4 6 8 12)
nRuns=(1 2 3)
theta=0.5
nIter=10

cd ../nbody
go build -o exec-$SLURM_JOB_ID
mv exec-$SLURM_JOB_ID ../benchmark/exec-$SLURM_JOB_ID
cd ../benchmark
trap '{ rm -f -- "exec-$SLURM_JOB_ID" "particles.dat"; }' EXIT
# rm slurm/out/$SLURM_JOB_ID.slurm1.stdout
for size in ${sizes[@]}
do
    echo 0 $size 1
    for nRun in ${nRuns[@]}
    do
        timeout 5m ./exec-$SLURM_JOB_ID $size $nIter $theta 1 0 1 1
    done
done

for type in ${types[@]}
do
    for size in ${sizes[@]}
    do
        for nThread in ${nThreads[@]}
        do
            echo $type $size $nThread
            for nRun in ${nRuns[@]}
            do
                timeout 5m ./exec-$SLURM_JOB_ID $size $nIter $theta 1 $type $nThread 1
            done
        done
    done
done

python bench.py $SLURM_JOB_ID ${#nRuns[@]} ${#sizes[@]} ${sizes[@]}