package main

import ( 
	"log"
	"net"
	"os"
	"fmt"
	"rpcData"
	"strings"
	"database/sql"	
	_ "github.com/mattn/go-sqlite3"
)

/* 获取数据库连接句柄 */
func getDbConn() *sql.DB{
	db, err:= sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return db
}

func handleMonitorDevices(infoBucket rpcData.InfoBucket, clientAddr string){
	fmt.Println("MonitorDevice:")
	fmt.Println("	MonitorDevice:", infoBucket.MonitorDevice)
	fmt.Println("--------------------------")
	sqlCreate := "create table if not exists monitor_devices(id integer primary key autoincrement, ip varchar(255), status varchar(255), base_status varchar(255), app_status varchar(255), restore_status varchar(255));"
	db.Exec(sqlCreate)
	sqlUpdate := fmt.Sprintf("update monitor_devices set status=1, base_status=%s, app_status=%s, restore_status=%s where ip=\"%s\"", 
		infoBucket.MonitorDevice.BaseStatus,
		infoBucket.MonitorDevice.AppStatus,
		infoBucket.MonitorDevice.RestoreStatus,
		infoBucket.MonitorDevice.Ip)
	fmt.Println(sqlUpdate)
	db.Exec(sqlUpdate)
}

func handleMonitorTime(infoBucket rpcData.InfoBucket, clientAddr string){
	fmt.Println("MonitorTime:")
	fmt.Println("	Time:", infoBucket.MonitorTime)
	fmt.Println("--------------------------")
	
	tableName := "monitor_time_" + clientAddr
	fmt.Println("--------------------------", tableName)
	sqlCreate := fmt.Sprintf("create table if not exists %s (id INTEGER PRIMARY KEY AUTOINCREMENT, monitor_time  VARCHAR(255))", tableName)
	fmt.Println(sqlCreate)
	db.Exec(sqlCreate)
	sqlInsert := fmt.Sprintf("insert into %s (monitor_time) values (\"%s\")", tableName, infoBucket.MonitorTime)
	fmt.Println(sqlInsert)
	db.Exec(sqlInsert)
}

func handleDeviceInfo(infoBucket rpcData.InfoBucket, clientAddr string){
	fmt.Println("DeviceInfo:")
	fmt.Println("	CpuMode:", infoBucket.DeviceInfo.CpuMode)
	fmt.Println("	CpuPhyCount:", infoBucket.DeviceInfo.CpuPhyCount)
	fmt.Println("	CpuPerCores:", infoBucket.DeviceInfo.CpuPerCores)
	fmt.Println("	CpuLogCount:", infoBucket.DeviceInfo.CpuLogCount)
	fmt.Println("	Mem:", infoBucket.DeviceInfo.Mem)
	fmt.Println("	Disk:", infoBucket.DeviceInfo.Disk)
	fmt.Println("--------------------------")
	
	tableName := "device_info_" + clientAddr
	fmt.Println("--------------------------", tableName)
	sqlCreate := fmt.Sprintf("create table if not exists %s (id INTEGER PRIMARY KEY AUTOINCREMENT, cpu_mode VARCHAR(255),cpu_phycount VARCHAR(255),cpu_percores VARCHAR(255),cpu_logcount VARCHAR(255),mem VARCHAR(255))", tableName)
	fmt.Println(sqlCreate)
	db.Exec(sqlCreate)
	sqlInsert := fmt.Sprintf("insert into %s (cpu_mode,cpu_phycount,cpu_percores,cpu_logcount,mem) values (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")", tableName, 
		infoBucket.DeviceInfo.CpuMode,
		infoBucket.DeviceInfo.CpuPhyCount,
		infoBucket.DeviceInfo.CpuPerCores,
		infoBucket.DeviceInfo.CpuLogCount,
		infoBucket.DeviceInfo.Mem)
	fmt.Println(sqlInsert)
	db.Exec(sqlInsert)
}

func handleUptime(infoBucket rpcData.InfoBucket, clientAddr string){
	fmt.Println("Uptime:")
	fmt.Println("	Time:", infoBucket.Uptime)
	fmt.Println("--------------------------")
	tableName := "uptime_" + clientAddr
	fmt.Println("--------------------------", tableName)
	sqlCreate := fmt.Sprintf("create table if not exists %s (id INTEGER PRIMARY KEY AUTOINCREMENT, uptime VARCHAR(255))", tableName)
	fmt.Println(sqlCreate)
	db.Exec(sqlCreate)
	sqlInsert := fmt.Sprintf("insert into %s (uptime) values (\"%s\")", tableName, 
		infoBucket.Uptime)
	fmt.Println(sqlInsert)
	db.Exec(sqlInsert)
}

func handleSysRunInfo(infoBucket rpcData.InfoBucket, clientAddr string){
	fmt.Println("SysRunInfo:")
	fmt.Println("--------------------------")
	
}

func handleProcInfo(infoBucket rpcData.InfoBucket, clientAddr string) {
	fmt.Println("ProcInfo:")
	fmt.Println("	Info:", infoBucket.ProcInfo)
	fmt.Println("--------------------------")
	
	tableName := "proc_info" + clientAddr
	fmt.Println("--------------------------", tableName)
	sqlCreate := fmt.Sprintf("create table if not exists %s (id INTEGER PRIMARY KEY AUTOINCREMENT, pname VARCHAR(255), pid VARCHAR(255), tid VARCHAR(255), tname VARCHAR(255), ts VARCHAR(255), virt VARCHAR(255), res VARCHAR(255), shr VARCHAR(255), mem VARCHAR(255), cpu VARCHAR(255),  run_time VARCHAR(255))", tableName)
	fmt.Println(sqlCreate)
	db.Exec(sqlCreate)
	for _,v := range infoBucket.ProcInfo.Tinfo {
		sqlInsert := fmt.Sprintf("insert into %s (pname, pid, tid, tname, ts , virt, res, shr, mem, cpu, run_time) values (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")", tableName, 
		infoBucket.ProcInfo.Pname,
		infoBucket.ProcInfo.Pid,
		v.Tid,
		v.Tname,
		v.Ts,
		v.Virt,
		v.Res,
		v.Shr,
		v.Mem,
		v.Cpu,
		v.Run_time)
		fmt.Println(sqlInsert)
		db.Exec(sqlInsert)
	}
	
}

/**
 * @brief - 服务端信息处理函数
 * @author - lc
 * @date - 2018-10-23
 **/
func recvMessage(client net.Conn, clientAddr string) error {
	var message []byte
	message = make([]byte, 4096)
	var infoBucket rpcData.InfoBucket
	
	for {
		len, _ := client.Read(message)
		if len > 0 {
			rpcData.Decode(message, &infoBucket)
			fmt.Println("infoBucket.InfoType =", infoBucket.InfoType)
			switch infoBucket.InfoType {
				case rpcData.StartMonitorTime:
					handleMonitorTime(infoBucket, clientAddr)
				case rpcData.StartDeviceInfo:
					handleDeviceInfo(infoBucket, clientAddr)
				case rpcData.RunUptime:
					handleUptime(infoBucket, clientAddr)
				case rpcData.RunSysRunInfo:
					handleSysRunInfo(infoBucket, clientAddr)
				case rpcData.RunProcInfo:
					handleProcInfo(infoBucket, clientAddr)
				case rpcData.RunMonitorDevice:
					handleMonitorDevices(infoBucket, clientAddr)
			}
		}
	}

	return nil
}

func checkRemoteAddr(addr net.Addr) (string, int){
	result := strings.Split(addr.String(), ":")
	
	
	return strings.Replace(result[0], ".", "_", -1),0
}

var db *sql.DB

/**
 * @brief - 服务端主程序
 * @author - lc
 * @date - 2018-10-23
 **/
func main() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var ipStr string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipStr = ipnet.IP.String()
			}
		}
	}
	
	ipStr += ":9700"
	fmt.Println(ipStr)
	server, err := net.Listen("tcp", ipStr)
	if err != nil {
		log.Fatal("start server failed!")
		os.Exit(1)
	}

	defer server.Close()
	
	log.Println("server is running")
	
	db = getDbConn()
	
	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatal("Accept errorr\n")
			continue
		}

		/* 获取客户端IP地址，检测是否有权限 */
		remoteAddr := client.RemoteAddr()
		clientAddr, ret := checkRemoteAddr(remoteAddr)
		if 0 == ret {
			fmt.Println("Client", clientAddr, "is connected!")
			go recvMessage(client, clientAddr)
		}else{
			fmt.Println("Invalid remote addr : ", remoteAddr.String())
		}
	}
}
