package handleQuotes

import (
	"github.com/bwmarrin/discordgo"
	"strings"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
    "fmt"
    "regexp"
    "math/rand"
    "time"
    "strconv"
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
    quoteId string
    tblName string
}

var sess  *session.Session
var svc *dynamodb.DynamoDB

func ParseQuote(message *discordgo.Message,guild_id string)string{
    sess = session.Must(session.NewSession())
	svc = dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-2"))

    
    content:= strings.TrimSpace(message.Content)
    content= strings.TrimPrefix(content,"!aq")
    content= strings.TrimSpace(content)
    short_form:=""
    said_by:=""
    quote:=""
    addedBy:=""
    
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
    
   //the query id will be the same as the message ID   
    quoteInput := QuoteStruct{
        quoteId: message.ID,
        said_by: said_by,
        addedByUsername: addedBy,
        quote: quote,
        short_form: short_form,
        tblName: guild_id+"QUOTES",
        Out: "ADDED QUOTE:\nid: "+message.ID +" said by: " + said_by + "\nadded by: " + addedBy+ "\nquote: "+ quote+"\nshort quote: "+ short_form,
    }     
    
	return    addQuote(&quoteInput)
}

func addQuote(qs *QuoteStruct) string{
    fmt.Println("top of add Quote")
    quoteItemIn := &dynamodb.PutItemInput{ 
    ReturnConsumedCapacity: aws.String("TOTAL"),
    TableName:              aws.String(qs.tblName),
        
        Item: map[string]*dynamodb.AttributeValue{
            "Quote_ID": {
                S: aws.String(qs.quoteId),
            },
            "Said_By": {
                S: aws.String(qs.said_by),
            },
            "Quote": {
                S: aws.String(qs.quote),
            },
             "addedByUsername": {
                S: aws.String(qs.addedByUsername),
            },
             "Short_Quote": {
                S: aws.String(qs.short_form),
            },
             "Output": {
                S: aws.String(qs.Out),
            },
            
        },
    }
    
    result, err := svc.PutItem(quoteItemIn)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case dynamodb.ErrCodeConditionalCheckFailedException:
                qs.Out = (dynamodb.ErrCodeConditionalCheckFailedException + aerr.Error())
            case dynamodb.ErrCodeProvisionedThroughputExceededException:
                qs.Out = (dynamodb.ErrCodeProvisionedThroughputExceededException + aerr.Error())
            case dynamodb.ErrCodeResourceNotFoundException:
                qs.Out = (dynamodb.ErrCodeResourceNotFoundException + aerr.Error())
            case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
                qs.Out = (dynamodb.ErrCodeItemCollectionSizeLimitExceededException + aerr.Error())
            case dynamodb.ErrCodeInternalServerError:
                qs.Out = (dynamodb.ErrCodeInternalServerError + aerr.Error())
            default:
                qs.Out = (aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            qs.Out = (err.Error())
        }
    }
    fmt.Println(result)
    return qs.Out
}

func GetQuote(message string,guild_ID string) string{
    //either an id or a name that will be quried 
    //expected message !quote 3193328238923
    //expected message !quote dough
    //expected message !quote !short_form
    fmt.Println("top of get quote")
    content:=strings.TrimPrefix(message,"!quote")
    content=strings.TrimSpace(content)
    //anything longer than 1 word will be denied
    matched, err := regexp.MatchString("[a-zA-Z]+[[:space:]]+[a-zA-Z]+", content)
    if (matched || content=="") {
        return "Please enter a person's name, a quote ID or a quote Name"
    }
    sess = session.Must(session.NewSession())
	svc = dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-2"))

    queryInput:=&dynamodb.QueryInput{}   
    fmt.Println(content)

    _, err = strconv.ParseFloat(content, 64)
    if(err == nil) {
        fmt.Println("quote id sent")
        //query quote
        queryInput = &dynamodb.QueryInput{
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":v1": {
                S: aws.String(message),
            },
        },
        KeyConditionExpression: aws.String("Quote_ID = :v1"),
        ProjectionExpression:   aws.String("Quote"),
        TableName:              aws.String(guild_ID+"QUOTES"),
        }               
    } else {
        //check for persons or short form
        //person: dough
        //short_form !qte
        fmt.Println("quote is not a number")
        if  strings.HasPrefix(content, "!") {   
             //query quick quote
            queryInput = &dynamodb.QueryInput{
            ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
                ":v1": {
                    S: aws.String(content),
                },
            },
            IndexName: aws.String("Short_Quote_IDX"),
            KeyConditionExpression: aws.String("Short_Quote = :v1"),
            ProjectionExpression:   aws.String("Quote"),
            TableName:              aws.String(guild_ID+"QUOTES"),
            }
            //this is a short form that has to be queried
            content=strings.TrimPrefix(content,"!")
            content= strings.TrimSpace(content)
            
        }else {
            //query person
            queryInput = &dynamodb.QueryInput{
            ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
                ":v1": {
                    S: aws.String(content),
                },
            },
            IndexName: aws.String("Said_By_IDX"),
            KeyConditionExpression: aws.String("Said_By = :v1"),
            ProjectionExpression:   aws.String("Quote"),
            TableName:              aws.String(guild_ID+"QUOTES"),

            }
        }           
    }
    
    result, err := svc.Query(queryInput)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case dynamodb.ErrCodeProvisionedThroughputExceededException:
                 fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
            case dynamodb.ErrCodeResourceNotFoundException:
                fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
            case dynamodb.ErrCodeInternalServerError:
                fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }   
    }
   
    rand.Seed(time.Now().UTC().UnixNano())
    var i int
    var output string
    if(0==len(result.Items)){
         i=0
         output="no quote found"
    }else{
         i=rand.Intn(len(result.Items))//return quote in between range [0-n)
         attributeValue:=result.Items[i]
         output=*attributeValue["Quote"].S
    }
    fmt.Println(output)
 

return output
}







