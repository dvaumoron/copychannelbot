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
	"path"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	basePath string
	data     map[string]any
}

func ReadConfig() Config {
	confPath := "bot.yaml"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}

	log.Println("Load configuration from", confPath)
	confBody, err := os.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprint("Unable to read configuration :", err))
	}

	confData := map[string]any{}
	if err = yaml.Unmarshal(confBody, confData); err != nil {
		panic(fmt.Sprint("Unable to parse configuration :", err))
	}
	return Config{basePath: path.Dir(confPath), data: confData}
}

func (c Config) Get(valueConfName string) string {
	value, _ := c.data[valueConfName].(string)
	return value
}

func (c Config) GetInt(valueConfName string) int64 {
	value, _ := c.data[valueConfName].(int64)
	return value
}

func (c Config) Require(valueConfName string) string {
	value := c.Get(valueConfName)
	if value == "" {
		panic("Configuration value is missing : " + valueConfName)
	}
	return value
}

func (c Config) RequireInt(valueConfName string) int64 {
	value, ok := c.data[valueConfName].(int64)
	if !ok {
		panic("Configuration value is missing or incorrect : " + valueConfName)
	}
	return value
}
