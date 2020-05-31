# Checkpoints

 { 203.880,   0.046, 193.930 },
 { 180.660,   0.046, 179.010 },
 { 173.030,   0.046, 208.420 },
 { 187.960,   0.046, 233.300 },
 { 165.230,   0.046, 232.210 },
 { 150.660,   0.046, 186.220 },
 { 180.000,   0.046, 162.200 },
 { 165.200,   0.046, 118.500 },
 { 121.300,   0.046,  96.600 },
 { 119.998,   0.046, 126.010 },
 { 112.000,   0.046, 194.000 },
 {  75.600,   0.046, 208.700 },
 {  60.480,   0.046, 208.960 },

# Available controls

DLL injection to execute code on a thread in the client
Pointermap to get a stable reference to the ServerManager instance
Proxy to intercept and rewrite client and server messages

The server keeps track of the client position but does not return it to the client
The clients send position updates and has internal teleport controls that the server accepts
The server sends a packet every time the client gets a checkpoint

Start the race at the final checkpoint? No teleporting, just go to final checkpoint, send the server packet, then run to the finish.

Why does DLL injection break the proxy? Change sockets when injecting DLL? 
=> Do injection before starting the proxy; change DLL code to get server manager on demand
What happens to the server comms when we use the teleport in the client?
=> teleport a few times and how the client updates change
How do we get a speed hack?
=> try to correlate position updates with velocity
=> load a separate client and look at the character while manipulating the y position in the proxy packets - float!

Race manager always reports 'too slow'
Teleporting onto checkpoints and disabling server teleport doesn't win the race
CheatEngine has speed hacking - it will probably screw up the heartbeats?
Turn online play off and do the race with the speed hack?

Can we manipulate packets to slow down client time?
Does the server always return the client heartbeat it received, along with its heartbeat?
Does the diff between the heartbeats need to match?
Can we record a race and then replay it faster?
Can we generate race data and then pick our speed?