package handleQuotes

import (
	"github.com/bwmarrin/discordgo"
	"strings"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws"
)

/*
this function reads the quote message and parses into seperate vairables
*/

//quote example: !aq !dough: this is a great quote
//quote example: !aq ..short_form !dough  this is an even better quote because it has a short form
//quote example: !aq !dough ..short_form this is an even better quote because it has a short form

//quote example: !aq this is a great quote by someone unknown
//quote example: !aq ..dhort_form: this is a great quote by an unknown with a short form

type  QuoteStruct struct{
    said_by string
    addedByUsername string
    quote string
    short_form string
    Out string  
    quoteId int64
}

func ParseQuote(message *discordgo.Message,guild_id string)QuoteStruct{
    content:= strings.TrimSpace(message.Content)
    content= strings.TrimPrefix(content,"!aq")
    content= strings.TrimSpace(content)
    short_form:=""
    said_by:=""
    quote:=""
    addedBy:=""
    //dateAdded:=
    
    //check for ! and : to determine short form or said by   
    if  strings.HasPrefix(content, "!") {
        said_by=  content[1:strings.Index(content," ")]
        content=strings.TrimPrefix(content,"!"+said_by)
        content= strings.TrimSpace(content)
     
        if strings.HasPrefix(content,"..") {
            short_form=  content[2:strings.Index(content," ")]
            content=strings.TrimPrefix(content,".."+short_form)
            content= strings.TrimSpace(content)
            }       
    }else if strings.HasPrefix(content,"..") {
        short_form=  content[2:strings.Index(content," ")]
        content=strings.TrimPrefix(content,".."+short_form)
        content= strings.TrimSpace(content)
       
        if strings.HasPrefix(content, "!") {
            said_by=  content[1:strings.Index(content," ")]
            content=strings.TrimPrefix(content,"!"+said_by)
            content= strings.TrimSpace(content)
            }
    }    
    quote=content
    addedBy=message.Author.Username
    if (said_by==""){
        said_by="NA"
    }
    if (addedBy==""){
        addedBy="NA"
    }
    if(quote==""){
        quote="NA"
    }
    if(short_form==""){
        short_form="NA"
    }
    
    
	output:= ("said by: " + said_by + "\nadded by: " +
              addedBy+ "\nquote: "+ quote+"\nshort quote: "+ short_form )
    
    
    quoteInput := QuoteStruct{} 
    quoteInput.said_by= said_by-
    quoteInput.addedByUsername= addedBy
    quoteInput.quote =quote
    quoteInput.short_form = short_form
    quoteInput.Out = output
    
   
   // addQuote(&quoteInput,guild_id+"QUOTES")
	return  quoteInput
}

func addQuote(qs *quoteStruct,quoteTableName string){
    //Query from table in order to get the max id then add one to it
    
    
    
    
    
    
    quoteItemIn := &dynamodb.PutItemInput{ 
    ReturnConsumedCapacity: aws.String("TOTAL"),
    TableName:              aws.String(quoteTableName),
    }
    quoteItemIn.Item = make(map[string]*dynamodb.AttributeValue)
    quoteItemIn.Item ["Quote_ID"]=aws.Int64()//find number
    quoteItemIn.Item ["dateAdded"]=aws.String(qs.said_by)
    quoteItemIn.Item ["Said_By"]=aws.String(qs.said_by)
    quoteItemIn.Item ["Quote"]=aws.String(qs.quote)
    quoteItemIn.Item ["addedByUsername"]=aws.String(qs.addedByUsername)
    quoteItemIn.Item ["Short_Quote"]=aws.String(qs.short_form)
    quoteItemIn.Item ["Output"]=aws.String(qs.Out)


    
}

