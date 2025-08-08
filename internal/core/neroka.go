package core

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Request struct {
	GuildId   string `json:"guild_id"`
	ChannelId string `json:"channel_id"`
	Content   string `json:"content"`

	//Message         discordgo.Message             `json:"message"`
	//Guild           discordgo.Guild               `json:"guild"`
	//Attachments     []discordgo.MessageAttachment `json:"attachments"`
	//UserName        string                        `json:"user_name"`
	//IsDirectMessage bool                          `json:"is_direct_message"`
	//UserId          int64                         `json:"user_id"`
}

type Conversations map[string]bool
type RecentParticipants map[string][]string

type Neroka struct {
	cancelChan  chan struct{}
	requestChan chan Request

	wg          sync.WaitGroup
	workerCount int

	session *discordgo.Session

	conversations      Conversations
	recentParticipants RecentParticipants
}

func NewNeroka(requestSize int, workerCount int) (*Neroka, error) {
	return &Neroka{
		cancelChan:  make(chan struct{}),
		requestChan: make(chan Request, requestSize),
		wg:          sync.WaitGroup{},
		workerCount: workerCount,
		//
		conversations: Conversations{},
	}, nil
}

func (n *Neroka) Start() {
	for i := 0; i < n.workerCount; i++ {
		n.wg.Add(1)
		go n.worker()
	}
}

func (n *Neroka) Stop() {
	close(n.cancelChan)
	n.wg.Wait()
}

func (n *Neroka) worker() {
	defer n.wg.Done()

	for {
		select {
		case <-n.cancelChan:
			return

		case request := <-n.requestChan:
			n.processSingleRequest(request)
		}
	}
}

func (n *Neroka) AddRequest(request Request) bool {
	select {
	case n.requestChan <- request:
		return true
	default:
		return false
	}
}

func (q *Neroka) processSingleRequest(request Request) {
	log.Printf("%+v\n", request)
}
