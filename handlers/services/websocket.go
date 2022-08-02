package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Shuixingchen/year-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type BlockHandler struct {
	Client  *Client
	BlockCh chan int
	WsConns map[*websocket.Conn]bool
	Mux     sync.Mutex
	CancleF context.CancelFunc
}

func NewBlockHandler() *BlockHandler {
	b := &BlockHandler{
		BlockCh: make(chan int),
		WsConns: make(map[*websocket.Conn]bool),
	}
	b.Client = NewClient(utils.Config.Nodes[80001])
	ctx, cancel := context.WithCancel(context.Background())
	b.CancleF = cancel
	go b.getLateBlockLoop(ctx)
	go b.broadcaster(ctx)
	return b
}

func (h *BlockHandler) LatestBlock(c *gin.Context) {
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": err})
		return
	}
	h.Mux.Lock()
	h.WsConns[wsconn] = true
	fmt.Printf("add wsconn %p", wsconn)
	h.Mux.Unlock()
	for {
		if _, _, err = wsconn.ReadMessage(); err != nil { //一直读消息，没有消息就阻塞,返回err表示断开连接
			h.Mux.Lock()
			delete(h.WsConns, wsconn)
			fmt.Printf("close wsconn %p", wsconn)
			h.Mux.Unlock()
			return
		}
	}
}

func (h *BlockHandler) getLateBlockLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("getLateBlockLoop stop")
			return
		default:
			number, err := h.Client.GetLateBlock()
			if err != nil {
				continue
			}
			h.BlockCh <- int(number)
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *BlockHandler) broadcaster(ctx context.Context) {
	for {
		select {
		case number := <-h.BlockCh:
			for c, _ := range h.WsConns {
				c.WriteMessage(1, []byte(strconv.Itoa(number)))
			}
		case <-ctx.Done():
			fmt.Println("broadcaster stop")
			return
		}
	}
}
