package Tables

import (
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws"
)



var QuoteTableIn = &dynamodb.CreateTableInput{
    AttributeDefinitions: []*dynamodb.AttributeDefinition{
      {
        AttributeName: aws.String("Quote_ID"),
        AttributeType: aws.String("S"),
      },
      {
      AttributeName: aws.String("Said_By"),
      AttributeType: aws.String("S"),
      },
      
     {
        AttributeName: aws.String("Short_Quote"),
        AttributeType: aws.String("S"),
      },    
    },
  
    TableName: aws.String("temp"),

      //everyQuote should have a unique ID
    KeySchema: []*dynamodb.KeySchemaElement{
      {
        AttributeName: aws.String("Quote_ID"),
        KeyType: aws.String("HASH"),
      },
    },
    
     ProvisionedThroughput: &dynamodb.ProvisionedThroughput {
      ReadCapacityUnits: aws.Int64(1),
      WriteCapacityUnits: aws.Int64(1),
    },

    
    //Can Query by said by and added by
    GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
        {
              IndexName: aws.String("Said_By_IDX"),

              KeySchema: []*dynamodb.KeySchemaElement{
                {
                  AttributeName: aws.String("Said_By"),
                  KeyType: aws.String("HASH"),
                },
                {
                  AttributeName: aws.String("Quote_ID"),
                  KeyType: aws.String("RANGE"),
                },
                  
              },

              ProvisionedThroughput:&dynamodb.ProvisionedThroughput {
                ReadCapacityUnits: aws.Int64(1),
                WriteCapacityUnits: aws.Int64(1),
              },

                Projection:  &dynamodb.Projection{
                ProjectionType : aws.String("ALL"),
                },
          },
          {
              IndexName: aws.String("Short_Quote_IDX"),

              KeySchema: []*dynamodb.KeySchemaElement{
                {
                  AttributeName: aws.String("Short_Quote"),
                  KeyType: aws.String("HASH"),
                },
              },

              ProvisionedThroughput:&dynamodb.ProvisionedThroughput {
                ReadCapacityUnits: aws.Int64(1),
                WriteCapacityUnits: aws.Int64(1),
              },

                Projection:  &dynamodb.Projection{
                ProjectionType : aws.String("ALL"),
                },
          },
    },
}