# expense-bot

this repo is deployed on google cloud run for me to log my expense

note 

sheetrange means tab in spreadsheet

## Date Format

this project use MM/DD/YYYY formate


## helpful gcloud command

if it asked for permission
 gcloud auth login

if cannot pull image from google cloud
 docker login -u oauth2accesstoken -p "$(gcloud auth print-access-token)" https://gcr.io
