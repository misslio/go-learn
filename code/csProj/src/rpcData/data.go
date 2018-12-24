package rpcData

import (
	"bytes"
	"encoding/gob"
	"time"
)

/* 定义数据包的type类型 */
type RpcDataType	int
/** 
 * Start前缀表示只在启动时发送一次的信息
 * Run前缀标识以一定时间间隔发送的运行信息
 **/
const (
	StartMonitorTime = iota	/* 监控时间，即客户端启动时间 */
	StartDeviceInfo	/* 设备配置信息 */
	RunUptime		/* 系统启动时间 */
	RunSysRunInfo	/* 系统运行信息 */
	RunProcInfo		/* 进程信息 */
	RunMonitorDevice	/* 监控设备运行概况 */
	RunMonitorProgStart	/* 程序的启动与停止情况 */
)

type Tinfo struct { 
	Tid	string		/* 线程id */
	Tname	string	/* 线程名字 */
	Ts	string		/* 时间 */
	Virt	string	/* 虚拟内存 */
	Res		string	/* 物理内存 */
	Shr		string	/* 共享内存 */
	Mem		string	/* 内存使用率 */
	Cpu		string	/* CPU使用率 */
	Run_time	string	/* 运行时间 */
}

type Pinfo struct {
	Pname string	/* 进程名字 */
	Pid	string		/* 进程id	*/
	Tinfo []Tinfo	/* 线程信息 */
}

type MonitorDeviceT struct {
	Ip	string
	BaseStatus string
	AppStatus string
	RestoreStatus string
}

type DeviceInfoT struct {
	CpuMode string		/* CPU型号 */
	CpuPhyCount string	/* CPU物理个数 */
	CpuPerCores	string	/* 每个CPU的核数 */
	CpuLogCount	string	/* CPU总的逻辑核数 */
	Mem	string
	Disk string
}

type SysInfoT struct {
	Info string
}

/* 进程的启动停止信息 */
type ProcStartStopInfoT struct {
	Datetime string	/* 进程启动时间 */
	Pid	string	/* 进程号 */
	Isfatal	string	/* 是否有fatal日志 */
	Isrun	string	/* 是否正在运行 */
}

type InfoBucket struct {
	InfoType int
	MonitorTime	time.Time
	Uptime string
	ProcInfo Pinfo
	DeviceInfo	DeviceInfoT
	SysInfo SysInfoT	
	MonitorDevice MonitorDeviceT
	ProcStartStopInfo map[string] map[string] ProcStartStopInfoT
}

func Encode(data InfoBucket) ([]byte, error) {  
	buf := bytes.NewBuffer(nil)  
	enc := gob.NewEncoder(buf)  
	err := enc.Encode(data)  
	if err != nil {  
		return nil, err  
	}
	
	return buf.Bytes(), nil  
}  
	
func Decode(data []byte, to *InfoBucket) error {  
    buf := bytes.NewBuffer(data)  
    dec := gob.NewDecoder(buf)  
    return dec.Decode(to)  
}  	