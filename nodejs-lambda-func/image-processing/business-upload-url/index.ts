import {APIGatewayProxyEventV2} from 'aws-lambda';
import AWS from 'aws-sdk';

const s3 = new AWS.S3()
// Create DynamoDB document client
let docClient = new AWS.DynamoDB.DocumentClient({apiVersion: '2012-08-10'});

const URL_EXPIRATION_SECONDS = process.env.URL_EXPIRATION_TIME || 300
const IMAGE_IMAGE_BUCKET = process.env.IMAGE_IMAGE_BUCKET
const BUCKET_PREFIX = process.env.BUCKET_PREFIX || 'raw'


async function checkIfBusinessExist(businessId: string) {

    let inp = {
        TableName: "XPOS_DEV",
        Key: {
            PK: `STORE#${businessId}`,
            SK: `STORE#${businessId}`
        }
    }
    try {
        let res = await docClient.get(inp).promise();
        if (res.Item) {
            return true;
        }
    } catch (e) {
        console.log(e);
        return false;
    }
    return false;
}

exports.handler = async (event: APIGatewayProxyEventV2) => {

    console.log(event);
    const {businessId} = event.pathParameters!

    // Check if the business exist
    let exist = await checkIfBusinessExist(businessId!);
    if (!exist) {
        return {
            statusCode: 404,
            body: JSON.stringify({
                "message": "Business Does not exist."
            })
        }
    }

    if (!event.queryStringParameters || !event.queryStringParameters.type) {
        return {
            statusCode: 400,
            body: JSON.stringify({
                "message": "query param type is required."
            })
        }
    }

    // Use Query parameter to get the corresponding signed url.
    const {type, fileName} = event.queryStringParameters!

    let s3Params = {
        Bucket: IMAGE_IMAGE_BUCKET,
        Key: '',
        Expires: URL_EXPIRATION_SECONDS,
        ContentType: 'image/jpeg',
    }

    switch (type) {
        case 'logo':
            s3Params.Key = `${BUCKET_PREFIX}/${businessId}/logo/logo.jpg`
            break;
        case 'bulkProductImage':
            s3Params.ContentType = 'multipart/x-zip'
            s3Params.Key = `${BUCKET_PREFIX}/${businessId}/bulkProductImage/${fileName || `${Date.now()}.zip`}`
            break;
        default:
            return {
                statusCode: 400,
                body: JSON.stringify({
                    "message": "Invalid type"
                })
            }
    }

    const uploadURL = await s3.getSignedUrlPromise('putObject', s3Params)
    return {
        statusCode: 200,
        body: JSON.stringify({
            uploadURL: uploadURL,
            filename: s3Params.Key
        })
    }
}