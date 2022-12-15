package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var email_data string = `Message-ID: <29790972.1075855665306.JavaMail.evans@thyme>
Date: Wed, 13 Dec 2000 18:41:00 -0800 (PST)
From: 1.11913372.-2@multexinvestornetwork.com
To: pallen@enron.com
Subject: December 14, 2000 - Bear Stearns' predictions for telecom in Latin
 America
Mime-Version: 1.0
Content-Type: text/plain; charset=us-ascii
Content-Transfer-Encoding: 7bit
X-From: Multex Investor <1.11913372.-2@multexinvestornetwork.com>
X-To: <pallen@enron.com>
X-cc: 
X-bcc: 
X-Folder: \Phillip_Allen_Dec2000\Notes Folders\All documents
X-Origin: Allen-P
X-FileName: pallen.nsf

Some content Juan Manuel Ticona Pacheco 
es contratado por Truora, por el buen desempoño mostrado.`

func TestParser_v2(t *testing.T) {
	data := DataRegexParcer(email_data)
	t.Run("return the expect quantity of rows.", func(t *testing.T) {
		assert.Equal(t, 16, len(data))
	})

	t.Run("DataRegexParcer", func(t *testing.T) {
		expect_data := `Some content Juan Manuel Ticona Pacheco 
es contratado por Truora, por el buen desempoño mostrado.`
		assert.Equal(t, expect_data, data["content"])
	})

	t.Run("DataRegexParcer empty value", func(t *testing.T) {
		assert.Equal(t, "", data["x-cc"])
	})

	t.Run("DataRegexParcer but content no one have next-line", func(t *testing.T) {
		assert.NotContains(t, data["subject"], "\n")
		assert.Contains(t, data["content"], "\n")
	})
}
