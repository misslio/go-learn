package main

import (
	"net"
	"os"
	"flag"
	"time"
	"sort"
	"fmt"
	"os/exec"
//	"strconv"
	"regexp"
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

		fmt.Println("Send sendUptime info, infoType is ", infoBucket.InfoType)
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
func getCmdResut(cmdStr string) (string, error){
	cmd := exec.Command("sh", "-c", cmdStr)
	
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(out)), err
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
	deviceInfo.CpuMode, _ = getCmdResut("cat /proc/cpuinfo | grep name | cut -f2 -d: | uniq -c")
	deviceInfo.CpuPhyCount, _  = getCmdResut("cat /proc/cpuinfo| grep \"physical id\"| sort| uniq| wc -l")
	deviceInfo.CpuPerCores, _  = getCmdResut("cat /proc/cpuinfo| grep \"cpu cores\"| uniq | cut -f2 -d:")
	deviceInfo.CpuLogCount, _  = getCmdResut("cat /proc/cpuinfo| grep \"processor\"| wc -l")
	
	
	/**
	 * 查看总的内存大小
	 *  cat /proc/meminfo | grep MemTotal| cut -f2 -d:
	 **/
	deviceInfo.Mem, _  = getCmdResut("cat /proc/meminfo | grep MemTotal| cut -f2 -d:") 
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

	fmt.Println("Send SystemInfo info, infoType is ", infoBucket.InfoType)
	
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
	fmt.Println("Send MonitorTime info, infoType is ", infoBucket.InfoType)

	sendBytes, err = rpcData.Encode(infoBucket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client.Write(sendBytes)
	
	client.Close()
}


/*** 下面代码通过日志监控进程的情况 ***/

/**
 * @desc - 解析启动日志监控
 * @paras - progName : 要监控的进程名称
 * 			dirName : 日志所在目录
 * @ret
 * @author - lc
 * @date - 2018-10-29 
 **/
 func monitorStartLog(progName string, dirName string) (map[string] rpcData.ProcStartStopInfoT) {
	var procStartStopInfo rpcData.ProcStartStopInfoT

	var err error
	var cmd, cmdResult string
	var result []string
	var startMap map[string] rpcData.ProcStartStopInfoT
	var fatalMap map[string] int
	startMap = make(map[string] rpcData.ProcStartStopInfoT)
	fatalMap = make(map[string] int)

	//查看是否有fatal日志
	cmd = "ls " + dirName + "/"+ progName + 
	".apm_analysis_fatal*"
	// fmt.Println(cmd)
	cmdResult, err = getCmdResut(cmd)

	if err == nil {
		result = strings.Split(cmdResult, "\n")

		for i:= 0; i < len(result); i++ {
			// fmt.Println(result[i])

			//获取进程id
			re, _ := regexp.Compile(`[\.]([0-9]*)[\.log]`)
			matchInfo := re.FindStringSubmatch(result[i])

			fatalMap[matchInfo[len(matchInfo) - 1]] = 1
		}
	}
	
	cmd = "ls " + dirName + "/"+ progName + 
		".apm_analysis_start*"
//	fmt.Println(cmd)
	cmdResult, err = getCmdResut(cmd)
	
	//查看是否有进程正在运行
	pidStr, _ := getPid(progName)

	if err == nil {
		result = strings.Split(cmdResult, "\n")

		for i:= 0; i < len(result); i++ {
	//		fmt.Println(result[i])	
			//获取日期
			re, _ := regexp.Compile(`(\d{8}-\d{6})`)
			matchInfo := re.FindStringSubmatch(result[i])
			//获取进程id
			re2, _ := regexp.Compile(`[\.]([0-9]*)[\.log]`)
			matchInfo2 := re2.FindStringSubmatch(result[i])

			procStartStopInfo.Datetime = matchInfo[len(matchInfo) - 1]
			procStartStopInfo.Pid = matchInfo2[len(matchInfo2) - 1]
			
			_, ok := fatalMap[matchInfo2[len(matchInfo2) - 1]]
			if ok {
				procStartStopInfo.Isfatal = "1"
			} else {
				procStartStopInfo.Isfatal = "0"
			}
			if pidStr == matchInfo2[len(matchInfo2) - 1] {
				procStartStopInfo.Isrun = "1"
			} else {
				procStartStopInfo.Isrun = "0"	
			}
			
			startMap[matchInfo2[len(matchInfo2) - 1]] = procStartStopInfo
		}
	}
	


	// fmt.Println(startMap)
	return startMap
}

func monitorProgStart() (map[string] map[string] rpcData.ProcStartStopInfoT){
	var progs = []string{"base_analysis", "app_analysis", "deep_analysis"}

	var progsStart map[string] map[string] rpcData.ProcStartStopInfoT
	progsStart = make(map[string] map[string] rpcData.ProcStartStopInfoT)

	for i := 0 ; i < len(progs); i++ {
		result := monitorStartLog(progs[i], "/var/log/iridium")
		progsStart[progs[i]] = result
	}

	// for k, v := range progsStart {
	// 	fmt.Println(k)
	// 	for _, v2 := range v {
	// 		fmt.Println("	", v2)
	// 	}
	// }

	return progsStart
}

func sendMonitorProgStart() {
	var sendBytes	[]byte
	var infoBucket rpcData.InfoBucket
	var progsStart map[string] map[string] rpcData.ProcStartStopInfoT

	progsStart = monitorProgStart()

	client, err := net.Dial("tcp", serverConn)
	if err != nil {
		fmt.Println("Client is dailing failed")
		os.Exit(1)
	}
	
	for {
		infoBucket.InfoType = rpcData.RunMonitorProgStart
		infoBucket.ProcStartStopInfo = progsStart

		fmt.Println("Send MonitorProgStart info, infoType is ", infoBucket.InfoType)
		sendBytes, err = rpcData.Encode(infoBucket)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		client.Write(sendBytes)
		time.Sleep(600 * time.Second)
	}
	client.Close()
}

/***  以下为读取解析运行日志文件，获取解析统计内容  ***/
/* 根据进程名称获取所有run目录下名称，日期大小形成数组 */
func getALLRunLogDate(progName string, dirName string) ([]int) {
	var cmd string
	var cmdResult, pidStr string
	var err error 
	var result  []string
	var logDate, newLogDate string
	var newResult []int

	/* 获取pid */
	pidStr,  err = getPid(progName)

	cmd = "ls " + dirName + " | grep " + pidStr

	cmdResult, err = getCmdResut(cmd)
	if err != nil {
		return newResult
	}

	result = strings.Split(cmdResult, "\n")
	fmt.Println(result)
	for i := 0; i < len(result); i++ {
		//获取日期
		re, _ := regexp.Compile(`(\d{8}-\d{6})`)
		matchInfo := re.FindStringSubmatch(result[i])

		logDate = matchInfo[len(matchInfo) - 1]

		newLogDate = ""
		for j := 0; j < len(logDate); j++ {
			if j == 4 {
				newLogDate += "-"
			} else if j == 6 {
				newLogDate += "-"
			}  else if j == 8 {
				newLogDate += " "
				continue
			} else if j == 11 {
				newLogDate += ":"
			} else if j == 13 {
				newLogDate += ":"
			}
			newLogDate += string(logDate[j])
		}
		fmt.Println("-->", newLogDate, "<--")
		loc, _ := time.LoadLocation("Local")
		tm, _ := time.ParseInLocation("2006-01-02 15:04:05", newLogDate, loc)
		fmt.Println("-->", tm.Unix(), "<--")
	//	newResult[i] = int(tm.Unix())
		newResult = append(newResult, int(tm.Unix()))
	}
	sort.Ints(newResult)
	fmt.Println(newResult)
	return newResult
}


/* 通过入参，获取解析线程的个数 */
func main() {
	serverConnInput := flag.String("s", "127.0.0.1:9700", "Server IP and Port")
	flag.Parse()

	getALLRunLogDate("base_analysis", "/var/log/iridium/run")

	serverConn = *serverConnInput
	fmt.Println("Server Address is : ", serverConn)
	sendMonitorTime()
	sendSystemInfo()
	
	go sendUptime(15)

	go sendProcInfo("base_analysis")
	
	go sendMonitorDevice()
	
	go sendMonitorProgStart()

	sendSystemRunInfo() /* 系统运行信息占用主线程 */
}
