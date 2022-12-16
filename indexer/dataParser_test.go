package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var email_data string = `Message-ID: <4350849.1075857628871.JavaMail.evans@thyme>
Date: Wed, 2 May 2001 07:35:00 -0700 (PDT)
From: john.arnold@enron.com
To: jeffrey.shankman@enron.com
Subject: 
Mime-Version: 1.0
Content-Type: text/plain; charset=us-ascii
Content-Transfer-Encoding: 7bit
X-From: John Arnold
X-To: Jeffrey A Shankman
X-cc: 
X-bcc: 
X-Folder: \John_Arnold_Jun2001\Notes Folders\Discussion threads
X-Origin: Arnold-J
X-FileName: Jarnold.nsf

Jeff:
To explain the P&L of -349,000  : 
We executed the trade when you gave the order (the delta anyway), first thing 
in the morning.  The market rallied 8 cents from the morning, with the back 
rallying about 2.5 cents.  On 904 PV contracts, curve shift was -226,000.  
The balance, $123,000,  is almost exactly $.01 bid/mid, which I think is 
pretty fair considering the tenor of the deal and that it included price and 
vol.  Cal 3 straddles, for instance, are $1.39 / $1.45.

Looking out for you bubbeh:
John`

func TestParser_v2(t *testing.T) {
	data, _ := DataRegexParser(email_data)

	t.Run("return the expect quantity of rows.", func(t *testing.T) {
		assert.Equal(t, 16, len(data))
	})

	t.Run("DataRegexParcer", func(t *testing.T) {
		expect_data := `Jeff:
To explain the P&L of -349,000  : 
We executed the trade when you gave the order (the delta anyway), first thing 
in the morning.  The market rallied 8 cents from the morning, with the back 
rallying about 2.5 cents.  On 904 PV contracts, curve shift was -226,000.  
The balance, $123,000,  is almost exactly $.01 bid/mid, which I think is 
pretty fair considering the tenor of the deal and that it included price and 
vol.  Cal 3 straddles, for instance, are $1.39 / $1.45.

Looking out for you bubbeh:
John`

		assert.Equal(t, expect_data, data["content"])
	})

	//================================================================
	t.Run("read the quantity by limit parameter", func(t *testing.T) {
		limit := 18000
		base_path := "/home/juancho/Programming/GoProjects/enron_mail_20110402/maildir/"
		readed_files := ListFiles(base_path, limit)

		assert.Equal(t, limit, len(readed_files))
	})

	t.Run("all have 16 fields", func(t *testing.T) {
		limit := 6000
		base_path := "/home/juancho/Programming/GoProjects/enron_mail_20110402/maildir/"
		readed_files := ListFiles(base_path, limit)

		for _, file := range readed_files {
			str, _ := os.ReadFile(file)
			lst, err := DataRegexParser(string(str))

			if err != nil {
				log.Fatal(err, file)
			}
			var vals = []int{15, 16}
			assert.Contains(t, vals, len(lst))
		}
	})

	t.Run("DataRegexParcer empty value", func(t *testing.T) {
		assert.Equal(t, "", data["x-cc"])
	})

	t.Run("DataRegexParcer but content no one have next-line", func(t *testing.T) {
		assert.NotContains(t, data["subject"], "\n")
		assert.Contains(t, data["content"], "\n")
	})
}
