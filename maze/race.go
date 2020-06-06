package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type position struct {
	x, y, z uint32
}

func (p position) moveTowards(dst position) position {
	step := 100.0
	dx := (float64(dst.x) - float64(p.x)) / step
	dz := (float64(dst.z) - float64(p.z)) / step
	l := math.Sqrt(dx*dx + dz*dz)
	if l < 1 {
		return dst
	}
	ux, uz := dx/l, dz/l
	// fmt.Println(dx, dz, l, ux*step, uz*step)
	return position{
		x: uint32((float64(p.x)/step + ux) * step),
		y: 0,
		z: uint32((float64(p.z)/step + uz) * step),
	}
}

func (p position) equals(other position) bool {
	return p.x == other.x && p.z == other.z
}

type playerCommand struct {
	*command
	position
	original []byte
	ts       uint32
}

func newPlayerCommand(cmdBytes []byte) *playerCommand {
	cmd := parseClientMsg(cmdBytes)
	pos := position{
		x: binary.LittleEndian.Uint32(cmdBytes[17:21]),
		y: binary.LittleEndian.Uint32(cmdBytes[21:25]),
		z: binary.LittleEndian.Uint32(cmdBytes[25:29]),
	}
	ts := binary.LittleEndian.Uint32(cmdBytes[9:13])
	return &playerCommand{cmd, pos, cmdBytes, ts}
}

// extracted by printing out memory values using injected dll
var checkpoints = []position{
	{2038800, 0, 1939300},
	{1806600, 0, 1790100},
	{1730300, 0, 2084200},
	{1879600, 0, 2333000},
	{1652300, 0, 2322100},
	{1506600, 0, 1862200},
	{1800000, 0, 1622000},
	{1652000, 0, 1185000},
	{1213000, 0, 966000},
	{1199980, 0, 1260100},
	{1120000, 0, 1940000},
	{756000, 0, 2087000},
	{604800, 0, 2089600},
}

var currentCheckpoint = 0

func teleportMsg(pos position) []byte {
	decoded := make([]byte, 14)
	decoded[0] = 0x54
	decoded[1] = 0x01
	binary.LittleEndian.PutUint32(decoded[2:6], pos.x)
	binary.LittleEndian.PutUint32(decoded[6:10], pos.y)
	binary.LittleEndian.PutUint32(decoded[10:14], pos.z)
	return decoded
}

func parseTeleportMsg(decoded []byte) position {
	pos := position{}
	pos.x = binary.LittleEndian.Uint32(decoded[2:6])
	pos.y = binary.LittleEndian.Uint32(decoded[6:10])
	pos.z = binary.LittleEndian.Uint32(decoded[10:14])
	return pos
}

func loginMsg() []byte {
	msg := make([]byte, 42)
	msg[0] = 0x4c
	copy(msg[1:9], secret)
	copy(msg[9:14], []byte{0x04, 0x63, 0x69, 0x6d, 0x69})
	fmt.Printf("%x\n", msg)
	return msg
}

func heartbeatMsg(ts uint32) []byte {
	msg := make([]byte, 18)
	msg[0] = 0x3c
	msg[1] = 0x33
	copy(msg[2:10], secret)
	binary.LittleEndian.PutUint32(msg[10:14], ts)
	binary.LittleEndian.PutUint32(msg[14:18], 0)
	return msg
}

func parseHeartbeatMsg(decoded []byte) (ts uint32) {
	return binary.LittleEndian.Uint32(decoded[2:6])
}

func movementMsg(ts uint32, pos position) []byte {
	msg := make([]byte, 46)
	msg[0] = 0x50
	copy(msg[1:9], secret)
	binary.LittleEndian.PutUint32(msg[9:13], ts)
	binary.LittleEndian.PutUint32(msg[17:21], pos.x)
	binary.LittleEndian.PutUint32(msg[21:25], pos.y)
	binary.LittleEndian.PutUint32(msg[25:29], pos.z)
	copy(msg[len(msg)-4:len(msg)-2], []byte{0x90, 0x01})
	return msg
}

type raceBot struct {
	serverConn *net.UDPConn
	serverAddr *net.UDPAddr
	outgoing   chan []byte
	incoming   chan []byte
	debug      chan []byte

	serverTs uint32
	clientTs uint32
	pos      position
	mux      sync.Mutex
}

var raceCommands = parseRaceDump("notes/race-full-dump-3.txt")
var raceTeleport = position{2965974, 0, 2321971}
var startTs = 750805

func winRace() {
	port := 1337 + rand.Intn(20)
	serverConn, err := net.DialUDP("udp", nil, addr(serverIP, port))
	if err != nil {
		panic(err)
	}
	bot := &raceBot{
		serverConn: serverConn,
		serverAddr: addr(serverIP, port),
		incoming:   make(chan []byte),
		outgoing:   make(chan []byte),
		debug:      make(chan []byte),
	}
	log.Printf("Bot connected to %s:%d\n", serverIP, port)
	rsp := bot.login()
	for {
		if rsp[0] == 0x4c {
			break
		}
		time.Sleep(5 * time.Second)
		rsp = bot.login()
	}
	log.Printf("Bot logged in successfully! %x", rsp)
	done := make(chan bool, 1)
	go bot.raceLoop()
	<-done
}

func (bot *raceBot) ensureServerConn() {
	if bot.serverConn == nil {
		serverConn, err := net.DialUDP("udp", nil, bot.serverAddr)
		if err != nil {
			fmt.Println(err)
		}
		bot.serverConn = serverConn
	}
}

func (bot *raceBot) sendOne(msg []byte) {
	log.Printf("Send: %s\n", parseClientMsg(msg))
	encoded := encode(msg, 0xfe, 0xed)
	bot.serverConn.WriteMsgUDP(encoded, nil, nil)
}

func (bot *raceBot) heartbeat() {
	bot.sendOne(heartbeatMsg(atomic.LoadUint32(&bot.clientTs)))
	atomic.AddUint32(&bot.clientTs, 10)
}

func (bot *raceBot) find(dst position) {
	if bot.pos.equals(dst) {
		return
	}
	bot.sendOne(movementMsg(atomic.LoadUint32(&bot.clientTs), bot.pos))
	bot.pos = bot.pos.moveTowards(dst)
}

func (bot *raceBot) move(dst position) {
	bot.sendOne(movementMsg(atomic.LoadUint32(&bot.clientTs), dst))
	bot.pos = dst
}

func (bot *raceBot) raceLoop() {
	reader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, 4096)
	atomic.StoreUint32(&bot.clientTs, uint32(startTs))
	bot.heartbeat()
	teleportCount := 0
	go func() {
		for {
			msg, err := bot.recvOne(buffer)
			if msg == nil || err != nil {
				bot.sendOne(heartbeatMsg(atomic.LoadUint32(&bot.clientTs)))
			}
			decoded := decode(msg)
			if decoded[0] == 0x58 || decoded[0] == 0x59 {
				panic("it ded")
			}
			cmd := parseServerMsg(decoded)
			switch cmd.name {
			case "heartbeat":
				bot.incoming <- decoded
			case "teleport":
				if teleportCount > 3 {
					panic("server is teleporting us")
				}
				bot.outgoing <- decoded
				teleportCount++
			case "player":
				// do nothing
			default:
				bot.debug <- decoded
			}
		}
	}()
	// teleport to the race
	for {
		select {
		case tMsg := <-bot.outgoing:
			// tPos := parseTeleportMsg(tMsg)
			bot.pos = parseTeleportMsg(tMsg)
			bot.sendOne(movementMsg(atomic.LoadUint32(&bot.clientTs), bot.pos))
			bot.heartbeat()
		case hMsg := <-bot.incoming:
			serverTs := parseHeartbeatMsg(hMsg)
			atomic.StoreUint32(&bot.serverTs, serverTs)
			clientTs := atomic.LoadUint32(&bot.clientTs)
			log.Printf("Client - server ts: %d", clientTs-serverTs)
			bot.find(raceTeleport)
			bot.heartbeat()
			time.Sleep(1)
		case msg := <-bot.debug:
			cmd := parseServerMsg(msg)
			if cmd.name == "race" {
				fmt.Printf("Recv: %s\nPress enter to continue", cmd)
				reader.ReadString('\n')
			}
			fmt.Println("Continuing")
		}

		if bot.pos.x != 0 && bot.pos.x < 2110000 {
			break
		}
	}

	// the index in the positions we've saved from the race right
	// pointing to right after going through the teleporter
	startIdx := 79
	curr := startIdx
	for {
		if bot.pos.x != 0 {
			select {
			case hMsg := <-bot.incoming:
				serverTs := parseHeartbeatMsg(hMsg)
				atomic.StoreUint32(&bot.serverTs, serverTs)
				clientTs := atomic.LoadUint32(&bot.clientTs)
				log.Printf("Client - server ts: %d", clientTs-serverTs)

				if bot.pos.equals(raceCommands[curr].position) {
					if curr%10 == 0 {
						fmt.Printf("Reached position %d/%d %v\n", curr, len(raceCommands), bot.pos)
					}
					curr++
				}
				atomic.StoreUint32(&bot.clientTs, raceCommands[curr].ts)
				bot.sendOne(raceCommands[curr].original)
				bot.pos = raceCommands[curr].position
				bot.heartbeat()

				time.Sleep(1)
			case <-bot.outgoing:
				curr = startIdx
				bot.heartbeat()
			case msg := <-bot.debug:
				cmd := parseServerMsg(msg)
				fmt.Printf("Recv: %s\n", cmd)
				if cmd.name != "race" {
					fmt.Println("Press enter to continue")
					reader.ReadString('\n')
					fmt.Println("Continuing")
				}
			}
		}
	}
}

func (bot *raceBot) recvOne(buffer []byte) ([]byte, error) {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		bot.ensureServerConn()
		n, _, err := bot.serverConn.ReadFromUDP(buffer)
		if err != nil {
			return nil, err
		}
		msg := make([]byte, n)
		copy(msg, buffer[:n])
		cmd := parseServerMsg(decode(msg))
		if cmd.name != "player" {
			log.Printf("Recv: %s\n", cmd)
			if cmd.name != "heartbeat" {
				fmt.Printf("Recv: %s\n", cmd)
			}
		}
		return msg, nil
	}
	// nothing received before deadline is reached
	return nil, nil
}

func (bot *raceBot) processIncoming(msg []byte) {
	decoded := decode(msg)
	log.Printf("Received %s\nOriginal %x", parseServerMsg(decoded), decoded)
}

func (bot *raceBot) login() []byte {
	rsp, err := http.Get("http://" + serverIP + "/api/login/queue")
	if err != nil {
		panic(err)
	}
	httputil.DumpResponse(rsp, true)
	buffer := make([]byte, 4096)
	for {
		bot.sendOne(loginMsg())
		msg, _ := bot.recvOne(buffer)
		if msg != nil {
			return decode(msg)
		}
	}
}

func parseRaceDump(path string) []*playerCommand {
	result := []*playerCommand{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Index(line, "CLIENT>") != 0 {
			continue
		}
		cmdBytes, err := hex.DecodeString(strings.Split(line, ">0x")[1])
		if err != nil {
			panic(err)
		}
		cmd := parseClientMsg(cmdBytes)
		if cmd.name == "player" {
			result = append(result, newPlayerCommand(cmdBytes))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}
