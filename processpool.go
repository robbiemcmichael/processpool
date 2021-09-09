package processpool

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
)

type Options struct {
	Address   string
	Port      int
	Processes int
}

type Process struct {
	Pid  int `json:"pid"`
	Port int `json:"port"`
}

type Pool struct {
	Options Options
	Chan    chan Process
	Ready   int64
}

func NewPool(options Options, path string, args []string) *Pool {
	pool := Pool{
		Options: options,
		Chan:    make(chan Process, options.Processes),
		Ready:   0,
	}

	// Start the desired number of processes
	for i := 0; i < options.Processes; i++ {
		go pool.runProcess(path, args, options.Port+i)
		time.Sleep(10 * time.Millisecond)
	}

	return &pool
}

func (p *Pool) Serve() error {
	http.HandleFunc("/", p.handler)

	return http.ListenAndServe(p.Options.Address, nil)
}

func (p *Pool) handler(w http.ResponseWriter, r *http.Request) {
	process := <-p.Chan

	// Process has been received from channel, decrement the counter
	atomic.AddInt64(&p.Ready, -1)

	// Read the counter
	ready := atomic.LoadInt64(&p.Ready)

	fmt.Printf("Served process received from channel: %+v\n", process)
	fmt.Printf("Processes in pool: %d\n", ready)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(process)
}

func (p *Pool) runProcess(path string, args []string, port int) {
	for {
		cmd := exec.Command(path, args...)

		// Pass the port to the process as an environment variable
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("PROCESSPOOL_PORT=%d", port),
		)

		// Start process
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start process: %v\n", err)
			continue
		}

		// Wait for process to be ready
		for {
			time.Sleep(1 * time.Second)
			conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
			if err != nil {
				continue
			}
			if conn != nil {
				break
			}
		}

		// Process is ready and has been sent to channel, increment the counter
		atomic.AddInt64(&p.Ready, 1)

		// Read the counter
		ready := atomic.LoadInt64(&p.Ready)

		// Send process to the channel
		process := Process{Pid: cmd.Process.Pid, Port: port}
		p.Chan <- process
		fmt.Printf("New process sent to channel: %+v\n", process)
		fmt.Printf("Processes in pool: %d\n", ready)

		// Wait for process to exit
		err = cmd.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Process %d: %v\n", cmd.Process.Pid, err)
		}
	}
}
