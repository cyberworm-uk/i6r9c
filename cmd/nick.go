package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
)

// randomName returns to a pointer to a string containing a pseudo-random name.
func randomName() *string {
	rand.Seed(time.Now().Unix())
	nick := strings.Replace(namesgenerator.GetRandomName(0), "_", "", 1)
	return &nick
}
