package main

import (
	"fmt"
	"math"
	"os"
	// "sort"
)

var data struct {
	R int // rows
	C int // columns
	F int // fleet
	N int // rides
	B int // bonus
	T int // steps
}

type ride struct {
	a     int // row start
	b     int // col start
	x     int // row end
	y     int // col end
	s     int // start
	f     int // finish
	taken bool
}

func (r ride) Start() pos {
	return pos{r.a, r.b}
}

func (r ride) End() pos {
	return pos{r.x, r.y}
}

var rides []ride

type startSort []ride

func (r startSort) Len() int {
	return len(r)
}

func (r startSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (p pos) Dist(q pos) int {
	return int(math.Abs(float64(p.x-q.x)) + math.Abs(float64(p.y-q.y)))
}

func (r startSort) Priority(i int) float64 {
	return float64(float64(r[i].f-r[i].s) / (math.Abs(float64(r[i].y-r[i].b)) + math.Abs(float64(r[i].x-r[i].a))))
}

func (r startSort) Less(i, j int) bool {

	return r.Priority(i) < r.Priority(j)
}

type pos struct {
	x int
	y int
}

type car struct {
	p      pos
	rId    int
	busyTS int
	rides  []int
	long   bool
}

var cars []car

func load() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	fmt.Fscanf(file, "%d %d %d %d %d %d", &data.R, &data.C, &data.F, &data.N, &data.B, &data.T)
	for i := 0; i < data.N; i++ {
		var r ride
		fmt.Fscanf(file, "%d %d %d %d %d %d", &r.a, &r.b, &r.x, &r.y, &r.s, &r.f)
		rides = append(rides, r)
	}
}

func save() {
	file, err := os.Create(os.Args[2])
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	tot := 0
	for _, c := range cars {
		tot += len(c.rides)
		fmt.Fprint(file, len(c.rides))
		for _, r := range c.rides {
			fmt.Fprintf(file, " %d", r)
		}
		fmt.Fprint(file, "\n")
	}
	fmt.Println(tot)
}

func (c car) EndTime(t, id int) int {

	cur, start, end := c.p, rides[id].Start(), rides[id].End()

	tts := cur.Dist(start) // travel time start
	tte := start.Dist(end) // travel time end
	wait := rides[id].s - (t + tts)

	if wait < 0 {
		wait = 0
	}

	return t + tts + wait + tte
}

func (c car) StartTime(id int) int {
	return c.p.Dist(rides[id].Start())
}

func (r ride) TimeLeft(t int) int {
	start, end := r.Start(), r.End()
	return r.f - (t + start.Dist(end))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("main [in] [out]")
		os.Exit(1)
	}

	load()

	cars = make([]car, data.F)

	for i := int(float64(len(cars))*0.275) + 1; i < len(cars); i++ {
		cars[i].long = true
	}

	// sort.Sort(startSort(rides))
	// fmt.Println(rides)

	for i := 0; i < data.T; i++ {
		for ic, c := range cars {
			if i > c.busyTS {
				cid, bt, bet := -1, 0, 0
				for id, r := range rides {
					if r.taken {
						continue
					}

					if c.long {
						t := r.TimeLeft(i)
						et := c.EndTime(i, id)

						if c.StartTime(id) > t {
							continue
						}
						if cid == -1 || t < bt || (t == bt && et < bet) {
							cid, bt, bet = id, t, et
						}
					} else {
						t := c.EndTime(i, id)
						if t > rides[id].f {
							continue
						}
						if cid == -1 || t < bt {
							cid, bt = id, t
						}
					}

				}

				if cid != -1 {
					rides[cid].taken = true
					cars[ic].busyTS = c.EndTime(i, cid)
					cars[ic].p = rides[cid].End()
					cars[ic].rides = append(cars[ic].rides, cid)
				}
			}
		}
	}
	save()
}
