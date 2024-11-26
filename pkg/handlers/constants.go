package handlers

import "gopkg.in/telebot.v4"

var WaitingForUserMessage = map[int64]bool{}
var WaitingForMessage = map[int64]bool{}
var AwaitngForward = map[int64]bool{}
var OriginalUserId int64
var ForwardedMsg *telebot.Message
var ChatID int64
var WaitingForOrganizationInfoMsg = map[int64]bool{}

var (
	StartMsg              = "Hi. I am a bot that can help to find out details about HS code and give out certain information about it such as:\n–whether this code is categorized on the sanctions list\n–category danger class."
	StartMsgWithOrgMsg    = "Hi. I am a bot that can help to find out details about HS code and give out certain information about it such as:\n–whether this code is categorized on the sanctions list\n–category danger class.\n\nNote: Please include the name of the organization in the following message, it is not necessary, if you don't want to just send a blank message."
	OrgMsg                = "Thank you!"
	BaseMsg               = "Use '/hs' command for check HS code"
	WaitingHsCodeMsg      = "Type HS code"
	CannotForwardedMsg    = "Sorry, I can't sent message to support chat"
	CompletlyForwardedMsg = "You message forwarded to support channel. Thank you!"
	HelpCommandMsg        = "The following message will be sent to the support chat"
)
