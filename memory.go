package main

// Memory usage statistics. Total, Available and Used contain numbers of bytes
// for human consumption.
//
// The other fields in this struct contain kernel specific values.
type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total string `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available string `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used string `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free string `json:"free"`

	// OS X / BSD specific numbers:
	// http://www.macyourself.com/2010/02/17/what-is-free-wired-active-and-inactive-system-memory-ram/
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
	Wired    string `json:"wired"`

	// FreeBSD specific numbers:
	// https://reviews.freebsd.org/D8467
	Laundry string `json:"laundry"`

	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	// https://www.kernel.org/doc/Documentation/vm/transhuge.txt
	Buffers        string `json:"buffers"`
	Cached         string `json:"cached"`
	WriteBack      string `json:"writeBack"`
	Dirty          string `json:"dirty"`
	WriteBackTmp   string `json:"writeBackTmp"`
	Shared         string `json:"shared"`
	Slab           string `json:"slab"`
	Sreclaimable   string `json:"sreclaimable"`
	Sunreclaim     string `json:"sunreclaim"`
	PageTables     string `json:"pageTables"`
	SwapCached     string `json:"swapCached"`
	CommitLimit    string `json:"commitLimit"`
	CommittedAS    string `json:"committedAS"`
	HighTotal      string `json:"highTotal"`
	HighFree       string `json:"highFree"`
	LowTotal       string `json:"lowTotal"`
	LowFree        string `json:"lowFree"`
	SwapTotal      string `json:"swapTotal"`
	SwapFree       string `json:"swapFree"`
	Mapped         string `json:"mapped"`
	VmallocTotal   string `json:"vmallocTotal"`
	VmallocUsed    string `json:"vmallocUsed"`
	VmallocChunk   string `json:"vmallocChunk"`
	HugePagesTotal string `json:"hugePagesTotal"`
	HugePagesFree  string `json:"hugePagesFree"`
	HugePagesRsvd  string `json:"hugePagesRsvd"`
	HugePagesSurp  string `json:"hugePagesSurp"`
	HugePageSize   string `json:"hugePageSize"`
	AnonHugePages  string `json:"anonHugePages"`
}

type SwapMemoryStat struct {
	Total       string  `json:"total"`
	Used        string  `json:"used"`
	Free        string  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
	Sin         string  `json:"sin"`
	Sout        string  `json:"sout"`
	PgIn        string  `json:"pgIn"`
	PgOut       string  `json:"pgOut"`
	PgFault     string  `json:"pgFault"`

	// Linux specific numbers
	// https://www.kernel.org/doc/Documentation/cgroup-v2.txt
	PgMajFault string `json:"pgMajFault"`
}
