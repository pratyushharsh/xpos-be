import { APIGatewayProxyEventV2 } from 'aws-lambda';
import AWS from 'aws-sdk';

const s3 = new AWS.S3()
// Create DynamoDB document client
let docClient = new AWS.DynamoDB.DocumentClient({apiVersion: '2012-08-10'});

const URL_EXPIRATION_SECONDS = process.env.URL_EXPIRATION_TIME || 300
const IMAGE_IMAGE_BUCKET = process.env.IMAGE_IMAGE_BUCKET


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

    const { businessId } = event.pathParameters!
    const Key = `raw/${businessId}/LOGO.jpg`


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

    // Get signed URL fro zS3
    const s3Params = {
        Bucket: IMAGE_IMAGE_BUCKET,
        Key,
        Expires: URL_EXPIRATION_SECONDS,
        ContentType: 'image/jpeg',
    }
    try {
        const uploadURL = await s3.getSignedUrlPromise('putObject', s3Params)

        return {
            statusCode: 200,
            body: JSON.stringify({
                uploadURL: uploadURL,
                filename: Key
            })
        }
    } catch (e) {
        console.error(e);
        return {
            statusCode: 500,
            body: JSON.stringify({
                "message": "Unable to generate Url"
            })
        }
    }
}