#!/bin/bash

# module load golang/1.19

sizes=(10000 100000)
types=(1 2)
nThreads=(2 4 6 8 12)
nRuns=(1 2 3)
SLURM_JOB_ID=99999999999999
theta=0.5
nIter=1

cd ../nbody
go build -o exec-$SLURM_JOB_ID
mv exec-$SLURM_JOB_ID ../benchmark/exec-$SLURM_JOB_ID
cd ../benchmark
trap '{ rm -f -- "exec-$SLURM_JOB_ID" "particles.dat"; }' EXIT
rm slurm/out/$SLURM_JOB_ID.slurm1.stdout
for size in ${sizes[@]}
do
    echo 0 $size 1 >> slurm/out/$SLURM_JOB_ID.slurm1.stdout
    for nRun in ${nRuns[@]}
    do
        timeout 5m ./exec-$SLURM_JOB_ID $size $nIter $theta 1 0 1 1 >> slurm/out/$SLURM_JOB_ID.slurm1.stdout
    done
done

for type in ${types[@]}
do
    for size in ${sizes[@]}
    do
        for nThread in ${nThreads[@]}
        do
            echo $type $size $nThread >> slurm/out/$SLURM_JOB_ID.slurm1.stdout
            for nRun in ${nRuns[@]}
            do
                timeout 5m ./exec-$SLURM_JOB_ID $size $nIter $theta 1 $type $nThread 1 >> slurm/out/$SLURM_JOB_ID.slurm1.stdout
            done
        done
    done
done

python bench.py $SLURM_JOB_ID ${#nRuns[@]} ${#sizes[@]} ${sizes[@]}