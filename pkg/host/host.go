package host

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yl2chen/cidranger"
)

const googleCloudIPURL = "https://www.gstatic.com/ipranges/cloud.json"

type googleCloudIPList struct {
    Prefixes []struct {
        IPv4Prefix string `json:"ipv4Prefix"`
    } `json:"prefixes"`
}



var ranger cidranger.Ranger

func init() {
    ranger = cidranger.NewPCTrieRanger()
    resp, err := http.Get(googleCloudIPURL)
    if err != nil {
        log.Fatalf("Failed to get Google Cloud IP list: %v", err)
    }
    defer resp.Body.Close()

    var list googleCloudIPList
    if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
        log.Fatalf("Failed to decode Google Cloud IP list: %v", err)
    }
    for _, prefix := range list.Prefixes {
        if prefix.IPv4Prefix == "" {
            continue
        }
        _, network, err := net.ParseCIDR(prefix.IPv4Prefix)
        if err != nil {
            log.Fatalf("Failed to parse CIDR: %v", err)
        }
        ranger.Insert(cidranger.NewBasicRangerEntry(*network))
    }
}

func isGoogleCloudIP(ip string) bool {
	ipNet := net.ParseIP(ip)
	contains, err := ranger.Contains(ipNet)
	if err != nil {
			log.Printf("Failed to check IP: %v", err)
            return false
	}
	return contains
}

func ValidateXForwardedFor(c *gin.Context) {
    lastIP := c.ClientIP()	

    // Only accept requests coming from either localhost (testing) or Google Cloud (container environment)
	if lastIP != "127.0.0.1" && lastIP != "::1" && !isGoogleCloudIP(lastIP) {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid X-Forwarded-For header"})
        return
    }

    c.Next()
}