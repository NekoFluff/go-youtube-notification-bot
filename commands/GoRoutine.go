package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func consume(name string, x chan int, done chan bool) {
	for {
		value, more := <-x
		if more {
			fmt.Printf("%s: %v\n", name, value)
		} else {
			fmt.Printf("%s: Done\n", name)
			done <- true
			return
		}
		time.Sleep(time.Duration(rand.Float64() / time.Nanosecond.Seconds() * 10))
	}
}

var GoRoutine = &DiscordCommand{Command: "!goroutine", Execute: func(s *discordgo.Session, m *discordgo.MessageCreate) {
	x := make(chan int)
	done := make(chan bool)
	start := time.Now()

	go consume("Consumer1", x, done)
	go consume("Consumer2", x, done)
	go consume("Consumer3", x, done)

	count := 24
	for i := 0; i < count; i++ {
		x <- i + 1
	}
	close(x)
	<-done
	<-done
	<-done

	duration := time.Since(start)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Start: %v | Duration: %v", start, duration))
}}
