package core

const (
	// Telnet commands. rfc854, rfc1116, rfc1123
	_     = 0
	EOF   = 0xec // end of file
	SUSP  = 0xed // suspend process
	ABORT = 0xee // abort process
	EOR   = 0xef // end of record (transparent mode, used for prompt marking)
	SE    = 0xf0 // end sub negotiation
	NOP   = 0xf1 // nop (used for keep alive messages	)
	DM    = 0xf2 // data mark--for connect. cleaning
	BREAK = 0xf3 // break
	IP    = 0xf4 // interrupt process (permanently)
	AO    = 0xf5 // abort output (but let program finish)
	AYT   = 0xf6 // are you there
	EC    = 0xf7 // erase the current character
	EL    = 0xf7 // erase the current line
	GA    = 0xf8 // you may reverse the line (used for prompt marking)
	SB    = 0xfa // interpret as subnegotiation
	WILL  = 0xfb // I will use option
	WONT  = 0xfc // I won"t use option
	DO    = 0xfd // please, you use option
	DONT  = 0xfe // you are not to use option
	IAC   = 0xff // interpret as command

	// Telnet Options. rfc855
	_                           = 0
	OptEcho                     = 0x01
	OptSuppressGoAhead          = 0x03
	OptTimingMark               = 0x06 // timing mark. rfc860
	OptTerminalType             = 0x18 // terminal type. rfc930, rfc1091
	OptEndOfRecord              = 0x19 // end of record. rfc885
	OptNegotiateAboutWindowSize = 0x1f // negotiate about window size. rfc1073
	OptLineMode                 = 0x22 // linemode. rfc1184
	OptEnviron                  = 0x24 // environment option. rfc1408
	OptNewEnviron               = 0x27 // new environment option. rfc1572
	OptCharset                  = 0x2a // character set. rfc2066

	// Options Terminal Type
	_                = 0
	TerminalTypeIs   = 0x00
	TerminalTypeSend = 0x01
)
