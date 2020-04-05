#!/bin/bash

## TODO before launching the import
# 1. insure the DB is empty
# 2. double check json info are OK
# 3. comment out the crdedential checks on the post => this does not work yet with the script

function migrateDate() {

    fieldName=$1
    legacyObj=$2

    unixTime=$(echo $legacyObj | jq --arg k $fieldName '.[$k] | tonumber')
    formattedDate=$(date -d @${unixTime} +'%Y-%m-%dT%H:%M:%S.%N+02:00')
    newObj=$(echo $legacyObj | jq --arg k $fieldName --arg v $formattedDate '. | .[$k] = $v')
    echo $newObj
}

function migrateDateAndKey() {

    fieldName=$1
    nfName=$2
    legacyObj=$3

    unixTime=$(echo $legacyObj | jq --arg k $fieldName '.[$k] | tonumber')
    formattedDate=$(date -d @${unixTime} +'%Y-%m-%dT%H:%M:%S.%N+02:00')
    newObj=$(echo $legacyObj | jq --arg k $nfName --arg v $formattedDate '. | .[$k] = $v')
    echo $newObj
}


echo "... Starting legacy data import"

fname=./upala-data-reformatted.json

slugs=$(jq '.posts[] | .slug' $fname)

index=0
for currSlug in $slugs; do
    echo "Importing $currSlug"

    # if [ $currSlug == "\"pensee\"" ]; then
    currPost=$(jq --arg ii $index '.posts[$ii | tonumber]' $fname)
    currPost=$(migrateDateAndKey "publicationDate" "publishedAt" "$currPost")
    currPost=$(migrateDateAndKey "createdOn" "createdAt" "$currPost")
    currPost=$(migrateDateAndKey "updatedOn" "updatedAt" "$currPost")

    # echo $currPost

    tokenStr="eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo5MDkyZGY4NC1jODY3LTRkYTAtODhmYi00MzhkYjQ0YTljMTQiLCJ0eXAiOiJKV1QifQ.eyJhdF9oYXNoIjoiOGpNV0FzRjRVWExNZ3lQY294Z19uZyIsImF1ZCI6WyJ2aXRybngiXSwiYXV0aF90aW1lIjoxNTc1OTI0MTM1LCJleHAiOjE1NzU5Mjc3MzgsImlhdCI6MTU3NTkyNDEzOCwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL29pZGMvIiwianRpIjoiMzQ0MTlkYjQtZGUxMi00NTNmLTllZmItNWNjZTY5NjMzNWViIiwibmFtZSI6ImFkbWluIiwibm9uY2UiOiIiLCJyYXQiOjE1NzU5MjQxMjksInNpZCI6ImNkOTY3MmY0LWY0Y2ItNDRjNS05NmZiLWRlNDRkYzg1NDljMiIsInN1YiI6IjE2YWJjYjU2LWRmNTYtNGZiMy1hODM5LWU0NWU2ZjI0MGIxNyJ9.vN4CsWX98bNrl_RDoiHHL_aHyn7dO9Kvl_3UIRjrB70FRjPytLYfTy-jmZ3P6gRBfuELfwPK3KYUqKTTesciAeudDHvCr7zIWQxBfihSXQtURyT-8vc0NJiZM4Bjg6uQ_jfQII7jm_34HOU41UfvCQ4ppnCIXUonYmRPn67cq054P8hhElYAX3Le1bWkzUlmJeHTJMsF3W2EGINZng5DnOVQo897RZiVY4lF_NJtmbz0nDc39K5KBVI_P_Y8W30HUBtn7vawwG2lWriFJpIbRegcuA_We4IVaNjVoWmDmhKyBxoPnaMCAXB7lu7PbcieRdZqM51XzmfBw6w1i3pn5lA-OpTTCEahXgiFhLSCUBZw9YFO9ZBpEhhkg5r8Y4cx_bp5VKZxlJ2vbqZI1yW5LQQBCI85_ZiedLV1t58EHXjMN-1iRkH_EI1lyxDl8FyfVADl9pVi9k2Ed4ty6mu2URQRKeRPEtr-e0tJ-hlypaq3Ge48imHg_S4iVzlbhFeGJgYjybjxQkuPPZ4avU-QnKpea5cAvMLi4iz_cNk7Iw7Yt2y0kvPaZPhmjzLKvR51C46ZGzm9Asv5iyzS2Ml1RrTfjOujDL7t-a8l1nFqikDOAYUqzSnLgMBHcZa3tQFvzKHcOnVY9uocqlx8J7ZvV7hbnr1o67fiyxDeM3rDkG0"
    curl -i -X POST -H "Content-Type:application/json" -H "Authorization:$tokenStr" http://localhost:3000/api/posts -d "$currPost"
    echo ""
    # fi
    ((index++))
    # exit early for testing purposes
    # exit 0
done

## Sample curls.

# To retrieve legacy data from a running instance
# curl -i -X GET -H "Content-Type:application/json" http://localhost:8888/api/posts
#
