package event

import "context"

type Conn interface {
	SendMessage(content string, sendReturn ...bool)
	SendMultipleMessage(content []string, sendReturn ...bool)

	Screen() string

	PrepareWait(name string)
	WaitUpdate(ctx context.Context, name string)
}

type GenericEvent struct {
	conn Conn
	wg   string
}

func (g *GenericEvent) Conn() Conn {
	return g.conn
}
func (g *GenericEvent) SetConn(conn Conn) {
	g.conn = conn
}

func (g *GenericEvent) Wait(ctx context.Context) {
	g.conn.WaitUpdate(ctx, g.wg)
}

func (g *GenericEvent) PrepareWait() {
	g.conn.PrepareWait(g.wg)
}
