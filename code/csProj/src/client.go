package main

import (
	"net"
	"os"
	"flag"
	"time"
//	"flag"
	"fmt"
	"os/exec"
//	"strconv"
//	"regexp"
	"strings"
	"rpcData"
)

//服务器IP和端口
var serverConn string

func getUptime() (string, error) {
	cmdStr := "uptime"
	cmd := exec.Command("sh", "-c", cmdStr)
	
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
    return strings.TrimSpace(string(out)), err
}

func getProcInfo(Pid string) (string, error){
	cmdStr := "/usr/bin/top -n 1 -H -b -p " + Pid + " | grep root"
	cmd := exec.Command("sh", "-c", cmdStr)
	
	out, err := cmd.Output()
	if err != nil {
	//	fmt.Println(err)
		return "", err
	}
	//fmt.Println(string(out))
	
    return string(out), err
}

func getPid(name string) (string, error) {
	cmd := exec.Command("/usr/sbin/pidof", name)
	
	out, err := cmd.Output()
	if err != nil {
//		fmt.Println(err)
		return "", err
	}
	/* 去除最后的回车 */
	pidStr := string(out)

	return strings.TrimSpace(pidStr), err
}

func handleProcInfo(info string) ([]rpcData.Tinfo) {
	var tinfo []rpcData.Tinfo
	var tinfoElement rpcData.Tinfo
//	var min_tinfo rpcData.Tinfo
	var line string
	
	result := strings.Split(info, "\n")
	
	for i := 0; i < len(result); i++{
		line = strings.TrimSpace(result[i])
		splitResult := strings.Split(line, " ")
		length := len(splitResult)
		var newResult []string
		for i:= 0; i < length; i++ {
//			fmt.Printf("==>%s<===\n", splitResult[i]);
			if len(splitResult[i]) != 0 {
				newResult = append(newResult, splitResult[i])
			}
			
		}
	
		for index, v := range newResult {
			switch index {
				case 0:	/* PID/TID */
					tinfoElement.Tid = v
				case 1:	/* USER */
				case 2:	/* PR */
				case 3:	/* NI */
				case 4:	/* VIRT */
					tinfoElement.Virt = v
				case 5:	/* RES */
					tinfoElement.Res = v
				case 6:	/* SHR */
					tinfoElement.Shr = v
				case 7:	/* S */
				case 8:	/* %CPU */
					tinfoElement.Cpu = v
				case 9:	/* %MEM */
					tinfoElement.Mem = v
				case 10:	/* TIME+ */
					tinfoElement.Run_time = v
				case 11:	/* COMMAND */
					tinfoElement.Tname = v
				
			}
		}
		tinfo = append(tinfo, tinfoElement)
		
	}

//	fmt.Printf("****%s****\n",  tinfo)
//	fmt.Println()
	
	return tinfo
}

/**
 * @brief - 发送系统启动时间
 * @params - count : 以秒为单位的发送时间间隔
 * @author - lc
 * @date - 2018-10-22
 **/
func sendUptime(count time.Duration) () {
	var uptime string
	var infoBucket 	rpcData.InfoBucket
	var sendBytes	[]byte
	
	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	for {
		/* 发送系统启动时间 */
		uptime, err = getUptime()
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		infoBucket.InfoType = rpcData.RunUptime
		infoBucket.Uptime = uptime
		sendBytes, err = rpcData.Encode(infoBucket)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		client.Write(sendBytes)
		time.Sleep(count * time.Second)
	}
	client.Close()
}

/**
 * @brief - 发送进程信息，top命令执行结果
 * @params - procName : 要监控的进程名称
 * @author - lc
 * @date - 2018-10-22
 **/
func sendProcInfo(procName string) () {
	var pidStr, procInfo string
	var infoBucket 	rpcData.InfoBucket
	var sendBytes	[]byte
	var pinfo rpcData.Pinfo
	var tinfo []rpcData.Tinfo
	
	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	for {
		pinfo.Pname = procName
		pidStr, err = getPid(pinfo.Pname)
		if err != nil{
			continue
		}
		pinfo.Pid = pidStr

		procInfo, err = getProcInfo(pidStr)
		if err != nil{
			continue
		}
	
		tinfo = handleProcInfo(procInfo)
		pinfo.Tinfo = tinfo
		
		infoBucket.InfoType = rpcData.RunProcInfo
		infoBucket.ProcInfo = pinfo
		fmt.Println("Send proc info, infoType is ", infoBucket.InfoType)
		sendBytes, err = rpcData.Encode(infoBucket)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		client.Write(sendBytes)
		time.Sleep(30 * time.Second)
	}
	client.Close()
}


/**
 * @brief - 发送监控设备的基本情况
 * @params - 
 * @author - lc
 * @date - 2018-10-25
 **/
func sendMonitorDevice() () {
	var sendBytes	[]byte
	var infoBucket 	rpcData.InfoBucket
	var MonitorDevice rpcData.MonitorDeviceT
	
	/* 获取本机IP */
	var ipStr string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipStr = ipnet.IP.String()
			}
		}
	}
	
	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	for {
		_, err = getPid("base_analysis")
		if err != nil{
			MonitorDevice.BaseStatus = "2"
		}else{
			MonitorDevice.BaseStatus = "1"
		}

		_, err = getPid("app_analysis")
		if err != nil{
			MonitorDevice.AppStatus = "2"
		}else{
			MonitorDevice.AppStatus = "1"
		}
		_, err = getPid("restore_analysis")
		if err != nil{
			MonitorDevice.RestoreStatus = "2"
		}else{
			MonitorDevice.RestoreStatus = "1"
		}
		MonitorDevice.Ip = ipStr
		infoBucket.InfoType = rpcData.RunMonitorDevice
		infoBucket.MonitorDevice = MonitorDevice
		
		fmt.Println("Send MonitorDevice info, infoType is ", infoBucket.InfoType)
		sendBytes, err = rpcData.Encode(infoBucket)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		client.Write(sendBytes)
		time.Sleep(30 * time.Second)
	}
	client.Close()
}

/**
 * @brief - 获取命令行运行结果
 * @params - None
 * @author - lc
 * @date - 2018-10-22
 **/
func getCmdResut(cmdStr string) string{
	cmd := exec.Command("sh", "-c", cmdStr)
	
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	return strings.TrimSpace(string(out))
}

/**
 * @brief - 发送系统运行状况
 * @params - None
 * @author - lc
 * @date - 2018-10-22
 **/
func sendSystemRunInfo() (){
	for {
		time.Sleep(1 * time.Second)
	}
}

/**
 * @brief - 发送设备硬件信息，只在启动时发送一次
 *			CPU, 内存，硬盘
 * @params - None
 * @author - lc
 * @date - 2018-10-22
 **/
func sendSystemInfo() () {
	var deviceInfo rpcData.DeviceInfoT
	/**
     * CPU型号，物理个数，总核数，总逻辑CPU数 
	 * 总核数 = 物理CPU个数 X 每颗物理CPU的核数 
	 * 总逻辑CPU数 = 物理CPU个数 X 每颗物理CPU的核数 X 超线程数
	 *
	 * CPU型号命令
	 *	cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c
	 * CPU物理个数
	 *  cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
	 * 每个物理CPU中的core个数
	 *  cat /proc/cpuinfo| grep "cpu cores"| uniq
	 * 逻辑CPU个数
	 *  cat /proc/cpuinfo| grep "processor"| wc -l
	 **/
	deviceInfo.CpuMode = getCmdResut("cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c")
	deviceInfo.CpuPhyCount = getCmdResut("cat /proc/cpuinfo| grep \"physical id\"| sort| uniq| wc -l")
	deviceInfo.CpuPerCores = getCmdResut("cat /proc/cpuinfo| grep \"cpu cores\"| uniq | cut -f2 -d:")
	deviceInfo.CpuLogCount = getCmdResut("cat /proc/cpuinfo| grep \"processor\"| wc -l")
	
	
	/**
	 * 查看总的内存大小
	 *  cat /proc/meminfo | grep MemTotal| cut -f2 -d:
	 **/
	deviceInfo.Mem = getCmdResut("cat /proc/meminfo | grep MemTotal| cut -f2 -d:") 
	/* 磁盘暂不支持 */
	
	/* 发送消息 */
	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	var infoBucket 	rpcData.InfoBucket
	var sendBytes	[]byte

	infoBucket.InfoType = rpcData.StartDeviceInfo
	infoBucket.DeviceInfo = deviceInfo
	
	sendBytes, err = rpcData.Encode(infoBucket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client.Write(sendBytes)

	client.Close()
}

/**
 * @brief - 发送客户端启动监控时间到服务端
 * @params - None
 * @author - lc
 * @date - 2018-10-23
 **/
func sendMonitorTime() {
	/* 发送消息 */
	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	var infoBucket 	rpcData.InfoBucket
	var sendBytes	[]byte

	infoBucket.InfoType = rpcData.StartMonitorTime
	infoBucket.MonitorTime = time.Now()
	
	sendBytes, err = rpcData.Encode(infoBucket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client.Write(sendBytes)
	
	client.Close()
}


/*** 下面代码通过日志监控进程的情况 ***/

func main() {
	serverConnInput := flag.String("s", "127.0.0.1:9700", "Server IP and Port")
	flag.Parse()

	serverConn = *serverConnInput
	fmt.Println("Server Address is : ", serverConn)
	sendMonitorTime()
	sendSystemInfo()
	
	go sendUptime(15)

	go sendProcInfo("base_analysis")
	
	go sendMonitorDevice()
	
	sendSystemRunInfo() /* 系统运行信息占用主线程 */
}
