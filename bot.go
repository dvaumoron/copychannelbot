/*
 *
 * Copyright 2025 copychannelbot authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	discordgo "github.com/bwmarrin/discordgo"
)

func main() {
	config := ReadConfig()

	guildId := config.Require("GUILD_ID")
	srcChannelName := config.Require("SRC_CHANNEL")
	destChannelName := config.Require("DEST_CHANNEL")

	session, err := discordgo.New("Bot " + config.Require("BOT_TOKEN"))
	if err != nil {
		panic(fmt.Sprint("Cannot create the bot :", err))
	}
	session.Identify.Intents |= discordgo.IntentMessageContent

	err = session.Open()
	if err != nil {
		panic(fmt.Sprint("Cannot open the session :", err))
	}
	defer session.Close()

	guildChannels, err := session.GuildChannels(guildId)
	if err != nil {
		panic(fmt.Sprint("Cannot retrieve the guild channels :", err))
	}

	srcChannelId := ""
	destChannelId := ""
	for _, channel := range guildChannels {
		switch channel.Name {
		case srcChannelName:
			srcChannelId = channel.ID
		case destChannelName:
			destChannelId = channel.ID
		}
	}

	switch {
	case srcChannelId == "":
		panic("Cannot retrieve the guild channel for source : " + srcChannelName)
	case destChannelId == "":
		panic("Cannot retrieve the guild channel for destination : " + destChannelName)
	}
	// for GC cleaning
	guildChannels = nil
	srcChannelName = ""
	destChannelName = ""

	session.AddHandler(func(s *discordgo.Session, u *discordgo.MessageCreate) {
		if u.ChannelID == srcChannelId {
			session.ChannelMessageSend(destChannelId, u.Content)
		}
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Started successfully")
	fmt.Println("Press Ctrl+C to exit")
	<-stop
}
