package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

var secret = []byte{0x9D, 0x64, 0x15, 0xA9, 0xAD, 0x96, 0x12, 0x66}

func decode(msg []byte) []byte {
	rand1, rand2 := int(msg[0]), int(msg[1])
	result := make([]byte, len(msg)-2)
	for idx, b := range msg[2:] {
		result[idx] = b ^ byte(rand1)
		rand1 = (((rand1 + rand2) & 0xff) + ((rand1 + rand2) / 0xff)) & 0xff
	}
	return result
}

func encode(msg []byte, rand1, rand2 int) []byte {
	result := make([]byte, len(msg)+2)
	result[0] = byte(rand1)
	result[1] = byte(rand2)
	for idx, b := range msg {
		result[idx+2] = b ^ byte(rand1)
		rand1 = (((rand1 + rand2) & 0xff) + ((rand1 + rand2) / 0xff)) & 0xff
	}
	return result
}

type command struct {
	name   string
	code   []byte
	secret []byte
	data   []byte
	parsed string
}

func byteStr(bytes []byte, count int) string {
	result := "|"
	groups := len(bytes) / count
	for i := 0; i < groups; i++ {
		result += fmt.Sprintf("%x|", bytes[i*count:(i+1)*count])
	}
	if count*groups < len(bytes) {
		result += fmt.Sprintf("%x|", bytes[groups*count:])
	}
	return result
}

func f64(bytes []byte) float64 {
	num := binary.LittleEndian.Uint32(bytes)
	return float64(num) / 10000
}

func istr(name string, in []byte) string {
	num := int32(binary.LittleEndian.Uint32(in))
	return fmt.Sprintf("%1s=%10d(0x%x)", name, num, in)
}

type printBytes func(data []byte) string

var parseFns = map[string]printBytes{
	"emoji": func(data []byte) string {
		return fmt.Sprintf("%2x", data)
	},
	"info": func(data []byte) string {
		return fmt.Sprintf("%9x|", data)
	},
	"player": func(data []byte) string {
		return parseUpdate(data)
	},
	"teleport": func(data []byte) string {
		nums := []int{}
		for i := 1; i < 10; i += 4 {
			nums = append(nums, int(binary.LittleEndian.Uint32(data[i:i+4])))
		}
		return fmt.Sprintf("%2x|x=%7d|y=%7d|z=%7d|", data[0], nums[0], nums[1], nums[2])
	},
	"server_player": func(data []byte) string {
		// multiple updates, 0x50 user_id(4 bytes) client_update(37 bytes) => 42 bytes per update
		updateLen := 42
		// prepend command byte to first one, make map keyed by user id?
		data = append([]byte{0x50}, data...)
		count := len(data) / updateLen
		result := "\n"
		for i := 0; i < count; i++ {
			chunk := make([]byte, updateLen)
			copy(chunk, data[i*updateLen:(i+1)*updateLen])
			uid := istr("uid", chunk[1:5])
			chunk = append([]byte{0x50}, chunk[6:]...)
			result += fmt.Sprintf("%s (%s)\n", parseUpdate(chunk), uid)
		}
		return strings.TrimRight(result, "\n")
	},
	"heartbeat": func(data []byte) string {
		// num := binary.LittleEndian.Uint64(data)
		return fmt.Sprintf("%08.4f|", f64(data))
	},
	"server_heartbeat": func(data []byte) string {
		return fmt.Sprintf("%08.4f|%08.4f|%08.4f|%08.4f|", f64(data[0:4]), f64(data[4:8]), f64(data[8:12]), f64(data[12:16]))
	},
}

func parseUpdate(data []byte) string {
	// return fmt.Sprintf("%x %x %x %x %x", data[0:8], data[8:16], data[16:24], data[24:32], data[32:])
	chunks := [][]byte{}
	for i := 0; i < 32; i += 4 {
		chunks = append(chunks, data[i:i+4])
	}
	posStr := fmt.Sprintf("%08.4f|%08.4f|%08.4f|%08.4f", f64(chunks[0]), f64(chunks[2]), f64(chunks[3]), f64(chunks[4]))
	// result := []string{
	// 	istr("ts", chunks[0]),
	// 	istr("?", chunks[1]),
	// 	istr("x", chunks[2]),
	// 	istr("y", chunks[3]),
	// 	istr("z", chunks[4]),
	// 	istr("?", chunks[5]),
	// 	istr("?c", chunks[6]),
	// 	istr("?", chunks[7]),
	// }
	rest := data[32:] // jump byte, maybe velocity, maybe some animation modifier?
	jumpByte := rest[0]
	velocity := binary.LittleEndian.Uint16(rest[1:3])
	animation := 255*rest[4] + rest[3]
	return fmt.Sprintf("%s|%x|%d(0x%x)|%d(0x%x)", posStr, jumpByte, velocity, rest[1:3], animation, rest[3:5])
}

func (c *command) String() string {
	name := c.name
	if name == "" {
		name = fmt.Sprintf("0x%x", c.code)
	}
	return fmt.Sprintf("%10s|%s", name, c.parsed)
}

var commandNames = map[string]string{
	"45":   "emoji",
	"3c33": "heartbeat",
	"49":   "info",
	"50":   "player",
	"54":   "teleport",
	"52":   "race",
	// three other server messages which are unknown: 0x52(R) when starting race and 0x46(F) when teleporting
	// teleport and send checkpoint (R) message
	// this will probably win the race on the client
	// time skew will probably trick the server
	// get the client flag today
}

func parseClientMsg(decoded []byte) *command {
	secretStart := bytes.Index(decoded, secret)
	if secretStart == -1 {
		return &command{data: decoded}
	}
	code := decoded[0:secretStart]
	name := commandNames[hex.EncodeToString(code)]
	data := decoded[secretStart+len(secret):]
	parsed := byteStr(data, 4)
	if parseFn, ok := parseFns[name]; ok {
		parsed = parseFn(data)
	}
	return &command{
		name:   name,
		code:   code,
		secret: secret,
		data:   data,
		parsed: parsed,
	}
}

func parseServerMsg(decoded []byte) *command {
	cmdStart := 1
	if bytes.Index(decoded, []byte{0x3c, 0x33}) == 0 {
		cmdStart = 2
	}
	code := decoded[0:cmdStart]
	name := commandNames[hex.EncodeToString(code)]
	data := decoded[cmdStart:]
	parsed := ""
	switch name {
	case "info", "teleport":
		parsed = parseFns[name](data)
	case "player":
	case "heartbeat":
		parsed = parseFns["server_"+name](data)
	default:
		parsed = byteStr(data, 4)
	}
	return &command{
		name:   name,
		code:   code,
		data:   data,
		parsed: parsed,
	}
}

func msgInterceptor(origin string, msg []byte) [][]byte {
	switch origin {
	case "CLIENT":
		return clientMsgInterceptor(origin, msg)
	case "SERVER":
		return serverMsgInterceptor(origin, msg)
		//return msg
	default:
		panic("Unkown origin: " + origin)
	}
}

// for map radar, plot coordinates of the bunny - done!
// leave proxy running to collect data

// for race, teleport on checkpoints
func wrap(msgs ...[]byte) [][]byte {
	batch := make([][]byte, len(msgs))
	for _, m := range msgs {
		batch = append(batch, m)
	}
	return batch
}

var teleported = false

func serverMsgInterceptor(origin string, msg []byte) [][]byte {
	decoded := decode(msg)
	cmd := parseServerMsg(decoded)

	if cmd.name == "xxteleport" {
		fmt.Printf("DROP TELE()|%s\n", cmd)
		if !teleported {
			currentCheckpoint = 0
			tMsg := teleportMsg(checkpoints[currentCheckpoint])
			rMsg := []byte{0x52, byte(currentCheckpoint)}
			fmt.Printf("%s(%3d)|%s\n", "MODSRV", len(tMsg), parseServerMsg(tMsg))
			fmt.Printf("%s(%3d)|%10s|%2d\n", "MODSRV", 2, "checkpoint", byte(currentCheckpoint))
			currentCheckpoint++
			teleported = true
			return wrap(
				encode(tMsg, int(msg[0]), int(msg[1])),
				encode(rMsg, int(msg[1]), int(msg[0])),
			)
		}
		return wrap()
	}
	fmt.Printf("%s(%3d)|%s\n", origin, len(decoded), cmd)
	fmt.Printf("%s>0x%x\n", origin, decoded)

	if cmd.name == "xemoji" {
		tMsg := teleportMsg(checkpoints[currentCheckpoint])
		rMsg := []byte{0x52, byte(currentCheckpoint)}
		fmt.Printf("%s(%3d)|%s\n", "MODSRV", len(tMsg), parseServerMsg(tMsg))
		fmt.Printf("%s(%3d)|%10s|%2d\n", "MODSRV", 2, "checkpoint", byte(currentCheckpoint))
		currentCheckpoint = (currentCheckpoint + 1) % len(checkpoints)
		return wrap(
			encode(tMsg, int(msg[0]), int(msg[1])),
			encode(rMsg, int(msg[1]), int(msg[0])),
		)
	}
	return wrap(encode(decoded, int(msg[0]), int(msg[1])))
}

func clientMsgInterceptor(origin string, msg []byte) [][]byte {
	decoded := decode(msg)
	cmd := parseClientMsg(decoded)
	// use emoji commands to teleport
	// we can intercept sendEmoji and send fake teleport messages
	// some of the update fields have to be emoji related
	if cmd.name == "xheartbeat" {
		// 1 byte command + 8 bytes secret + 4 bytes timestamp + 4 bytes first blank + 4 bytes first number = 21 bytes offset
		// height := binary.LittleEndian.Uint32(decoded[21:25])
		// binary.LittleEndian.PutUint32(decoded[21:25], height+50000)
		// TODO: make sure to include the secret in the client message - probably yes

		// this doesn't seem to do anything
		//decoded[len(decoded)-3] = 0x1
		//decoded[len(decoded)-4] = 0xc8
		// try to increase speed
		// maybe speed needs to be correlated with the position update?

		// locked emoji 0x38 => 0xe
		return wrap(msg)
	}
	cmd = parseClientMsg(decoded)

	fmt.Printf("%s(%3d)|%s\n", origin, len(decoded), cmd)
	fmt.Printf("%s>0x%x\n", origin, decoded)

	return wrap(encode(decoded, int(msg[0]), int(msg[1])))
}
