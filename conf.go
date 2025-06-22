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

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	botToken       string `yaml:"BOT_TOKEN"`
	guildId        string `yaml:"GUILD_ID"`
	srcChannelName string `yaml:"SRC_CHANNEL"`
	refreshRate    int    `yaml:"REFRESH_RATE"`
	port           int    `yaml:"PORT"`
	tmplPath       string `yaml:"TMPL_PATH"`
	cutUntil       string `yaml:"CUT_UNTIL"`
}

func ReadConfig() *Config {
	confPath := "bot.yaml"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}

	log.Println("Load configuration from", confPath)
	confBody, err := os.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprint("Unable to read configuration :", err))
	}

	var conf Config
	if err = yaml.Unmarshal(confBody, &conf); err != nil {
		panic(fmt.Sprint("Unable to parse configuration :", err))
	}

	conf.validate()

	return &conf
}

func (c *Config) validate() {
	switch {
	case c.botToken == "":
		panic("Configuration value is missing : BOT_TOKEN")
	case c.guildId == "":
		panic("Configuration value is missing : GUILD_ID")
	case c.srcChannelName == "":
		panic("Configuration value is missing : SRC_CHANNEL")
	case c.refreshRate <= 0:
		panic("Configuration value is missing or incorrect : REFRESH_RATE")
	case c.port <= 0 || c.port > 65535:
		panic("Configuration value is missing or incorrect : PORT")
	case c.tmplPath == "":
		panic("Configuration value is missing : TMPL_PATH")
	}
}
