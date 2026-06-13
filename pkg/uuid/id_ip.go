package uuid

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	LogIdIp = "LogID"
)

var ip_instance string

// 获取本机 IPv4 地址（非 127.0.0.1 回环地址）
// 只获取一次，结果缓存到 ip_instance
func getLocalIP() (ip string) {
	if ip_instance != "" {
		return ip_instance
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		ip_instance = "127.0.0.1"
		return
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// 找到第一个 IPv4 地址
				ip = ipnet.IP.String()
				break
			}
		}
	}
	ip_instance = ip
	return ip
}

// 生成一个日志唯一 ID（LogID）
func GenLogIdIp() string {
	// 用当前时间微秒作为随机数种子
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	// 获取本机 IP
	ip := getLocalIP()

	now := time.Now()
	timestamp := uint32(now.Unix()) // 秒级时间戳（32位）
	timeNano := now.UnixNano()      // 纳秒级时间戳
	pid := os.Getpid()              // 当前进程号
	b := bytes.Buffer{}

	// 写入 IP（转成十六进制字符串，例如 192.168.1.10 -> c0a8010a）
	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))

	// 写入时间戳（秒，8位十六进制）
	b.WriteString(fmt.Sprintf("%x", timestamp&0xffffffff))

	// 写入时间纳秒低 16 位（保证同一秒内唯一性）
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))

	// 写入进程 PID 的低 16 位（保证不同进程不冲突）
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))

	// 写入一个 24 位随机数（6 位十六进制，增加随机性）
	b.WriteString(fmt.Sprintf("%06x", r.Int31n(1<<24)))

	// 固定后缀 "b0"，可能是内部约定的标识符
	b.WriteString("b0")

	// 返回拼接好的唯一 ID
	return b.String()
}
