package bot

import (
	"../config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
  "os"
	"os/signal"
	"syscall"
	"../handleQuotes"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"../Tables"
)

var BotID string
var goBot *discordgo.Session
var killBot chan bool
var err error
var sess  *session.Session
var svc *dynamodb.DynamoDB

//initalize the bot check if the guild/server is stored in aws
func Init() {
	goBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}


	sess = session.Must(session.NewSession())
	// Create a DynamoDB client with additional configuration
	svc = dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-2"))

 	killBot := make(chan bool)
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func(){
		fmt.Println(killBot)
		<-sc
		killBot <- true
			go func(){
				<-sc
				os.Exit(1)
			}()
}()

	go main()
	<- killBot
	fmt.Println("at the bottom of the the start")
}

func main(){

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		if m.Author.ID == BotID {
			return
		}
        
		//create new aws tables
		if m.Content == config.BotPrefix+"init" {
				currentChannel,err := goBot.Channel(m.ChannelID)

				if err != nil{
					fmt.Println(err)
					_, _ = s.ChannelMessageSend(m.ChannelID, "error initalizing")
				}else{                    
					message:= CreateTables(currentChannel.GuildID)
					_, _ = s.ChannelMessageSend(m.ChannelID, message)
                    fmt.Println(message)
				}
		}
        if m.Content == config.BotPrefix+"kill" {
			_, _ = s.ChannelMessageSend(m.ChannelID, "killed")
			kill()
		}
        
		if  strings.HasPrefix(m.Content, config.BotPrefix+"aq") {
            //parse quote
            currentChannel,err := goBot.Channel(m.ChannelID)
				if err != nil{
					fmt.Println(err)
					_, _ = s.ChannelMessageSend(m.ChannelID, "error adding quote")
				}
            
    
            output:= handleQuotes.ParseQuote(m.Message,currentChannel.GuildID)
            fmt.Println(output)   
			_, _ = s.ChannelMessageSend(m.ChannelID, currentChannel.GuildID)
            _, _ = s.ChannelMessageSend(m.ChannelID, "quote added")

		}

	}
}

//create all tables that the bot will use
//if a table already exists then do not create a new one
func CreateTables(guild_id string) string{
fmt.Println("create table")

	quoteTableName:=guild_id+"QUOTES"
	describeTableInput := &dynamodb.DescribeTableInput{TableName: aws.String(quoteTableName)}
    fmt.Println("Past describeTableInput")

    
  _,err :=svc.DescribeTable(describeTableInput)
        fmt.Println("tried to describe table")


	//determine
	if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
                //this is the one that we want to come true to create new tables
            case dynamodb.ErrCodeResourceNotFoundException:
            fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
            fmt.Println("Resource not found.  This is what we want to happen")

                 return createTable(quoteTableName)
            case dynamodb.ErrCodeInternalServerError:
                fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
                return aerr.Error()
                //	return "Server Issues, try again latter"
            default:
                fmt.Println(aerr.Error())
                return aerr.Error()

        }
    }
	}


	return "table(s) already exist"
}

func createTable(tblName string)string {
            fmt.Println("try to create " + tblName)

quoteTableIn := Tables.QuoteTableIn
quoteTableIn.TableName = aws.String(tblName)

_, err := svc.CreateTable(quoteTableIn)

output:= "tables initalized"

if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
			output=aerr.Error()
    } else {
      output= err.Error()
    }
}
return output
}


func kill(){
	os.Exit(1)
}



