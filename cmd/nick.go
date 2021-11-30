package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
)

func randomName() *string {
	rand.Seed(time.Now().Unix())
	nick := strings.Replace(namesgenerator.GetRandomName(0), "_", "", 1)
	return &nick
}
