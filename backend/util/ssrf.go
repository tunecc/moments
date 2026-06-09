package fs_util

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ValidateExternalURL 校验外部 URL,防止 SSRF。
// 仅允许 http/https,且解析出的 IP 不能指向内网/环回/链路本地等地址。
func ValidateExternalURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("非法的 URL: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("仅支持 http/https 协议")
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL 缺少主机名")
	}

	// 解析所有 IP(域名可能解析到内网)
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("无法解析主机名: %w", err)
	}

	for _, ip := range ips {
		if isPrivateOrReserved(ip) {
			return fmt.Errorf("禁止访问内网或保留地址: %s", ip.String())
		}
	}

	return nil
}

// isPrivateOrReserved 判断 IP 是否为内网/环回/链路本地/未指定等需要拒绝的地址。
func isPrivateOrReserved(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}

	// 额外拦截 IPv4 的 100.64.0.0/10 (CGNAT) 与 IPv6 唯一本地地址已被 IsPrivate 覆盖
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 100 && ip4[1]&0xc0 == 0x40 {
			return true
		}
	}

	return false
}
