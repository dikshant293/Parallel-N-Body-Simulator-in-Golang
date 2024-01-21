package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"proj3-redesigned/barrier"
	"proj3-redesigned/bdequeue"
	"proj3-redesigned/qTree"
	"proj3-redesigned/worksteal"
	"strconv"
	"time"
)

func rand_init(pArray []qTree.Particle, nParticles int, r *rand.Rand) {
	for i := 0; i < nParticles; i++ {
		pArray[i].Mass = r.Float64() * 100.0
		pArray[i].X = r.Float64()*20.0 - 10.0
		pArray[i].Y = r.Float64()*20.0 - 10.0
		pArray[i].Vx = r.Float64()*20.0 - 10.0
		pArray[i].Vy = r.Float64()*20.0 - 10.0
	}
}

func circular_init(pArray []qTree.Particle, nParticles int) {
	rad_per_particle := 2.0 * math.Pi / float64(nParticles)
	radius := 5.0 / math.Sin(rad_per_particle)
	for i := 0; i < nParticles; i++ {
		pArray[i].X = radius * math.Cos(float64(i)*rad_per_particle)
		pArray[i].Y = radius * math.Sin(float64(i)*rad_per_particle)
		pArray[i].Mass = 100.0
		v := math.Sqrt(float64(nParticles)*pArray[i].Mass/radius) / 2.0
		pArray[i].Vx = -1.0 * v * math.Sin(float64(i)*rad_per_particle)
		pArray[i].Vy = v * math.Cos(float64(i)*rad_per_particle)
	}
}

func orbit_init(pArray []qTree.Particle, nParticles int) {
	pArray[nParticles-1].X = 0.0
	pArray[nParticles-1].Y = 0.0
	pArray[nParticles-1].Vx = 0.0
	pArray[nParticles-1].Vy = 0.0
	pArray[nParticles-1].Mass = 1000.0
	nParticles--
	rad_per_particle := 2.0 * math.Pi / float64(nParticles)
	radius := 5.0 / math.Sin(rad_per_particle)
	for i := 0; i < nParticles; i++ {
		pArray[i].X = radius * math.Cos(float64(i)*rad_per_particle)
		pArray[i].Y = radius * math.Sin(float64(i)*rad_per_particle)
		pArray[i].Mass = 100.0
		v := math.Sqrt(float64(nParticles)*pArray[i].Mass/radius) / 2.0
		pArray[i].Vx = -1.0 * v * math.Sin(float64(i)*rad_per_particle)
		pArray[i].Vy = v * math.Cos(float64(i)*rad_per_particle)
	}
}

func calc_force_on_particle(p *qTree.Particle, root *qTree.QTNode, F []float64, theta float64) {
	if root == nil {
		return
	}
	dx := root.P.X - p.X
	dy := root.P.Y - p.Y
	distSqr := dx*dx + dy*dy + 1e-9
	if (root.Rb-root.Lb)/math.Sqrt(distSqr) < theta || root.Size == 1 {
		invDist := 1.0 / math.Sqrt(distSqr)
		invDist3 := invDist * invDist * invDist
		F[0] += p.Mass * root.Total_mass * dx * invDist3
		F[1] += p.Mass * root.Total_mass * dy * invDist3
	} else {
		for i := 0; i < 4; i++ {
			calc_force_on_particle(p, root.Child[i], F, theta)
		}
	}
}

func divideArray(totalSize int, n int, isEqual bool) []int {
	subarraySizes := make([]int, n)
	startIndices := make([]int, n+1)
	subarraySizes[n-1] = totalSize
	startIndices[n] = totalSize
	for i := 0; i < n-1; i++ {
		subarraySizes[i] = rand.Intn(totalSize)
		totalSize -= subarraySizes[i]
	}
	startIndex := 0
	for i := 0; i < n; i++ {
		startIndices[i] = startIndex
		if isEqual {
			startIndex += subarraySizes[n-1] / n
		} else {
			startIndex += subarraySizes[i]
		}
	}
	return startIndices
}

func get_start_end(i, per_thread, total int) (int, int) {
	start := per_thread * i
	if start > total {
		start = total
	}
	end := start + per_thread
	if end > total {
		end = total
	}
	return start, end
}

func update_particles(start int, end int, pArray []qTree.Particle, root *qTree.QTNode, dt float64, theta float64, bar *barrier.Barrier) {
	var F []float64 = make([]float64, 2)
	for i := start; i < end; i++ {
		F[0] = 0.0
		F[1] = 0.0
		calc_force_on_particle(&pArray[i], root, F, theta)
		pArray[i].Vx += dt * F[0] / pArray[i].Mass
		pArray[i].Vy += dt * F[1] / pArray[i].Mass
	}
	if bar != nil {
		bar.Wait()
	}
}

func move_particles(start int, end int, pArray []qTree.Particle, dt float64, bar *barrier.Barrier) {
	for i := start; i < end; i++ {
		pArray[i].X += pArray[i].Vx * dt
		pArray[i].Y += pArray[i].Vy * dt
	}
	if bar != nil {
		bar.Wait()
	}
}

func simulate_timesteps(pArray []qTree.Particle, nParticles int, niters int, theta float64, dt float64, f *os.File) {
	if is_benchmark == 0 {
		fmt.Println("seq")
	}
	for k := 0; k < niters; k++ {
		var root *qTree.QTNode = qTree.Create_node(nil, -1, pArray, nParticles)
		for i := 0; i < nParticles; i++ {
			qTree.QTree_insert(&pArray[i], root)
		}
		root = qTree.Remove_empty_nodes(root)
		update_particles(0, nParticles, pArray, root, dt, theta, nil)
		move_particles(0, nParticles, pArray, dt, nil)

		if is_benchmark==0 && file_output==1{
			for i := 0; i < nParticles; i++ {
				f.WriteString(fmt.Sprintf("%g %g %g %g %g\n", pArray[i].X, pArray[i].Y, pArray[i].Vx, pArray[i].Vy, pArray[i].Mass))
			}
		}
		if is_benchmark == 0 {
			fmt.Printf("\riteration: %d", k+1)
		}
	}
	if is_benchmark == 0 {
		fmt.Printf("\n")
	}
}

func par_helper(start int, end int, pArray []qTree.Particle, root *qTree.QTNode, dt float64, theta float64, b1, b2 *barrier.Barrier) {
	update_particles(start, end, pArray, root, dt, theta, b1)
	b1.Wait()
	move_particles(start, end, pArray, dt, b2)
	b2.Wait()
}

func par_simulate_timesteps(pArray []qTree.Particle, nParticles int, niters int, theta float64, dt float64, f *os.File, nthreads int) {
	if is_benchmark == 0 {
		fmt.Println("parallel : nthreads = ", nthreads)
	}
	for k := 0; k < niters; k++ {
		var root *qTree.QTNode = qTree.Create_node(nil, -1, pArray, nParticles)
		for i := 0; i < nParticles; i++ {
			qTree.QTree_insert(&pArray[i], root)
		}
		root = qTree.Remove_empty_nodes(root)

		par_per_thread := (nParticles + nthreads - 1) / nthreads

		b1 := barrier.NewBarrier(int32(nthreads))
		b2 := barrier.NewBarrier(int32(nthreads) + 1)

		for i := 0; i < nthreads; i++ {
			chunkStart, chunkEnd := get_start_end(i, par_per_thread, nParticles)
			go par_helper(chunkStart, chunkEnd, pArray, root, dt, theta, b1, b2)
		}
		b2.Wait()

		if is_benchmark==0 && file_output==1{
			for i := 0; i < nParticles; i++ {
				f.WriteString(fmt.Sprintf("%g %g %g %g %g\n", pArray[i].X, pArray[i].Y, pArray[i].Vx, pArray[i].Vy, pArray[i].Mass))
			}
		}
		if is_benchmark == 0 {
			fmt.Printf("\riteration: %d", k+1)
		}
	}
	if is_benchmark == 0 {
		fmt.Printf("\n")
	}
}

func populateQueue(i int, queueArray []*bdequeue.Queue, chunkIdxs []int, idxArr []int64, nParticles int, ntreads int, bar *barrier.Barrier) {

	queueArray[i] = bdequeue.NewBDEQueue(nParticles)
	for j := chunkIdxs[i]; j < chunkIdxs[i+1]; j++ {
		queueArray[i].PushBottom(&idxArr[j])
	}
	bar.Wait()
}

func makeQueueArray(nthreads int, nParticles int, idxArr []int64, isEqual bool) []*bdequeue.Queue {
	queueArray := make([]*bdequeue.Queue, nthreads)
	chunkIdxs := divideArray(nParticles, nthreads, isEqual)
	b := barrier.NewBarrier(int32(nthreads) + 1)
	for i := 0; i < nthreads; i++ {
		go populateQueue(i, queueArray, chunkIdxs, idxArr, nParticles, nthreads, b)
	}
	b.Wait()
	return queueArray
}

func work_steal_simulate_timesteps(pArray []qTree.Particle, nParticles int, niters int, theta float64, dt float64, f *os.File, nthreads int) {
	if is_benchmark == 0 {
		fmt.Println("worksteal : nthreads = ", nthreads)
	}
	idxArr := make([]int64, nParticles)
	workArr := make([]int64, nParticles+1)
	for i := 0; i < nParticles; i++ {
		idxArr[i] = int64(i)
	}

	for k := 0; k < niters; k++ {
		var root *qTree.QTNode = qTree.Create_node(nil, -1, pArray, nParticles)
		for i := 0; i < nParticles; i++ {
			qTree.QTree_insert(&pArray[i], root)
		}
		root = qTree.Remove_empty_nodes(root)

		var a1, a2 []int64 = make([]int64, nParticles), make([]int64, nParticles)
		b1 := barrier.NewBarrier(int32(nthreads) + 1)
		b2 := barrier.NewBarrier(int32(nthreads) + 1)

		queueArray := makeQueueArray(nthreads, nParticles, idxArr, false)
		w := worksteal.NewWorkStealingThread(queueArray, 1)
		for i := 0; i < nthreads; i++ {
			go w.Run(i, worksteal.MyStruct{Idx: i, PArray: &pArray, Root: root, Dt: dt, Theta: theta, Bar: b1}, update_particles, move_particles, &workArr, &a1)
		}
		b1.Wait()

		queueArray = makeQueueArray(nthreads, nParticles, idxArr, true)
		w = worksteal.NewWorkStealingThread(queueArray, 1)
		for i := 0; i < nthreads; i++ {
			go w.Run(i, worksteal.MyStruct{Idx: i, PArray: &pArray, Root: nil, Dt: dt, Theta: theta, Bar: b2}, update_particles, move_particles, &workArr, &a2)
		}
		b2.Wait()

		if is_benchmark==0 && file_output==1{
			for i := 0; i < nParticles; i++ {
				f.WriteString(fmt.Sprintf("%g %g %g %g %g\n", pArray[i].X, pArray[i].Y, pArray[i].Vx, pArray[i].Vy, pArray[i].Mass))
			}
		}

		if is_benchmark == 0 {
			fmt.Printf("\riteration: %d", k+1)
		}
	}
	if is_benchmark == 0 {
		fmt.Printf("\n")
	}
}

func compareParticles(p1, p2 *[]qTree.Particle, nParticles int) (bool, int) {
	ans := true
	mismatch := 0
	for i := 0; i < nParticles; i++ {
		if (*p1)[i].X != (*p2)[i].X || (*p1)[i].Y != (*p2)[i].Y || (*p1)[i].Vx != (*p2)[i].Vx || (*p1)[i].Vy != (*p2)[i].Vy || (*p1)[i].Mass != (*p2)[i].Mass {
			ans = false
			mismatch++
		}
	}
	return ans, mismatch
}

func run_simulation(pArray []qTree.Particle, nParticles int, nIters int, theta float64, init_type int, dt float64, nthreads int) {
	var f *os.File
	if is_benchmark==0 && file_output==1{
		f, _ = os.Create("particles.dat")
		defer f.Close()
		f.WriteString(fmt.Sprintf("%d %d %f %d %d %d\n", nParticles, nIters, theta, init_type, nthreads, run_type))
		for i := 0; i < nParticles; i++ {
			f.WriteString(fmt.Sprintf("%g %g %g %g %g\n", pArray[i].X, pArray[i].Y, pArray[i].Vx, pArray[i].Vy, pArray[i].Mass))
		}
	} else {
		f = nil
	}
	if run_type == 0 {
		simulate_timesteps(pArray, nParticles, nIters, theta, dt, f)
	} else if run_type == 1 {
		par_simulate_timesteps(pArray, nParticles, nIters, theta, dt, f, nthreads)
	} else if run_type == 2 {
		work_steal_simulate_timesteps(pArray, nParticles, nIters, theta, dt, f, nthreads)
	} else {
		fmt.Println("test")
		var pArrayWS []qTree.Particle = make([]qTree.Particle, nParticles)
		var pArrayPar []qTree.Particle = make([]qTree.Particle, nParticles)
		for j := 0; j < nParticles; j++ {
			pArrayWS[j] = pArray[j]
			pArrayPar[j] = pArray[j]
		}
		fmt.Println("\nsequential begins")
		startTime := time.Now()
		simulate_timesteps(pArray, nParticles, nIters, theta, dt, f)
		endTime := time.Since(startTime).Seconds()
		fmt.Printf("\nsequential took %f seconds\n", endTime)

		fmt.Println("\nparallel begins")
		startTime = time.Now()
		par_simulate_timesteps(pArrayPar, nParticles, nIters, theta, dt, f, nthreads)
		endTime = time.Since(startTime).Seconds()
		fmt.Printf("\nparallel took %f seconds\n", endTime)
		if bol, mis := compareParticles(&pArray, &pArrayPar, nParticles); bol {
			fmt.Println("Parallel - CORRECT")
		} else {
			fmt.Println("Parallel - INCORRECT ", bol, mis)
		}

		fmt.Println("\nwork steal begins")
		startTime = time.Now()
		work_steal_simulate_timesteps(pArrayWS, nParticles, nIters, theta, dt, f, nthreads)
		endTime = time.Since(startTime).Seconds()
		fmt.Printf("\nworkstealing took %f seconds\n", endTime)
		if bol, mis := compareParticles(&pArray, &pArrayWS, nParticles); bol {
			fmt.Println("Work steal - CORRECT")
		} else {
			fmt.Println("Work steal - INCORRECT ", bol, mis)
		}
	}
}

var run_type int = 0
var is_benchmark int = 0
var file_output int = 0

const usage = "go run nbody.go [nParticles] [nIter] [theta] [init_type] [run_type] [ntheads] [is_benchmark] [file_output]"

func main() {
	r := rand.New(rand.NewSource(69))
	var nParticles int = 1000
	var nIter int = 200
	var theta float64 = 0.5
	var dt float64 = 0.01
	var init_type = 1
	var nthreads = 0
	if len(os.Args) > 1 {
		nParticles, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		nIter, _ = strconv.Atoi(os.Args[2])
	}
	if len(os.Args) > 3 {
		theta, _ = strconv.ParseFloat(os.Args[3], 32)
	}
	if len(os.Args) > 4 {
		init_type, _ = strconv.Atoi(os.Args[4])
		if init_type < 1 || init_type > 3 {
			log.Panic("wrong init_type must be between 1 to 3: ", usage)
		}
	}
	if len(os.Args) > 5 {
		run_type, _ = strconv.Atoi(os.Args[5])
		if run_type < 0 || run_type > 3 {
			log.Panic("wrong run_type must be between 0 to 3: ", usage)
		}
	}
	if len(os.Args) > 6 {
		nthreads, _ = strconv.Atoi(os.Args[6])
	}
	if len(os.Args) > 7 {
		is_benchmark, _ = strconv.Atoi(os.Args[7])
		if is_benchmark < 0 || is_benchmark > 1 {
			log.Panic("wrong is_benchmark, must be 0 or 1")
		}
	}
	if len(os.Args) > 8 {
		file_output, _ = strconv.Atoi(os.Args[8])
		if file_output < 0 || file_output > 1 {
			log.Panic("wrong file_output, must be 0 or 1")
		}
	}

	var pArray []qTree.Particle = make([]qTree.Particle, nParticles)
	switch init_type {
	case 1:
		rand_init(pArray, nParticles, r)
	case 2:
		circular_init(pArray, nParticles)
	case 3:
		orbit_init(pArray, nParticles)
	}

	start := time.Now()
	run_simulation(pArray, nParticles, nIter, theta, init_type, dt, nthreads)
	end := time.Since(start).Seconds()
	if run_type == 3 {
		fmt.Printf("\nTotal time taken = %f\n", end)
	} else {
		fmt.Printf("%f\n", end)
	}
}
