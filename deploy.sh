GOOGLE_PROJECT_ID='tixologi-partner'
CLOUD_RUN_SERVICE='partner-prod'

gcloud builds submit --tag gcr.io/$GOOGLE_PROJECT_ID/$CLOUD_RUN_SERVICE \
  --project=$GOOGLE_PROJECT_ID

gcloud run deploy $CLOUD_RUN_SERVICE \
  --image gcr.io/$GOOGLE_PROJECT_ID/$CLOUD_RUN_SERVICE \
  --add-cloudsql-instances $INSTANCE_CONNECTION_NAME \
  --platform managed \
  --region us-west2 \
  --allow-unauthenticated \
  --project=$GOOGLE_PROJECT_ID