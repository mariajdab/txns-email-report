# Txns-email-report

## Starting 🚀

The csv file to process needs to be on the root repository, in this case the file is ```txns.csv```

git clone https://github.com/mariajdab/txns-email-report.git

### Open a new terminal window 
Copy and run the command ```docker-compose up --build```

### Open a new terminal window 

Copy and run the command ```docker run --rm --net=txns-email-report_default -v "${PWD}:/input" txns-reporter ./input/txns.csv```

#### What happened!? 🚀

The csv file is validated, the first line should be the expected: ```Id,Date,Transaction```, also each field or row is parsed in order to be validated. 
In the other hand, we are creating the report or transactions summary. When we have the report done then we use a fake email sender (mailtrap) in order mock the 
acccion of sending the report via email. At the end we insert the data in to a database, in this case PostgreSQL.

## Note
If we check ```.env```, you'll see the mailtrap configuration information needed, that match with the mailtrap account that I created, certainly that information will change in a few days.

