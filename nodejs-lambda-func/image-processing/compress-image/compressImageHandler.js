const AWS = require("aws-sdk");

let docClient = new AWS.DynamoDB.DocumentClient()

const s3Service = require('./service/s3Service').s3Service
const sharpService = require('./service/sharpService').sharpService

const INPUT_BUCKET = process.env.INPUT_BUCKET
const OUTPUT_BUCKET = process.env.OUTPUT_BUCKET
const TABLE_NAME = "XPOS_DEV"

const imageSizes = {
    logo_small: [100, 100],
    logo_medium: [300, 300],
    logo_large: [600, 600],
    banner: [1400, 350],
    content: [274, 377]
}

async function convertImage(imageData, quality) {
    const convertedImageData = await sharpService.convertImage(imageData, quality)

    // Build and return the image key.
    const image_key = `${key.split('.').slice(0, -1).join('.')}_${size[0]}x${size[1]}.jpg`
    await s3Service.saveImage(OUTPUT_BUCKET, `${image_key}`, convertedImageData)

    // Insert Into Dynamodb
    return image_key
}

async function insertLogoIntoDynamoDb(images, store_id) {
    console.log(images)

    try {
        const config = {
            TableName: TABLE_NAME,
            Key: {
                PK: `STORE#${store_id}`,
                SK: `STORE#${store_id}`
            },
            ConditionExpression: "attribute_exists(PK)",
            UpdateExpression: "SET #logo = :logo",
            ExpressionAttributeNames: {
                "#logo": "logo"
            },
            ExpressionAttributeValues: {
                ":logo": {
                    small: `https://${OUTPUT_BUCKET}.s3.ap-south-1.amazonaws.com/${await images[0]}`
                },
            },
            ReturnValues:"UPDATED_NEW"
        };
        console.log(JSON.stringify(config))
        let resp = await docClient.update(config).promise()
        console.log(resp)
    } catch (e) {
        console.log(e)
        console.error(`Unable to update in dynamodb`)
    }

}

async function processRecord(record) {
    try {
        const splitPath = record.s3.object.key.split('/')
        const imageData = await s3Service.getObject(INPUT_BUCKET, record.s3.object.key)

        // Process Record According To The Type of data
        const QUEUE = []
        QUEUE.push(convertImage(imageData, "out/" + record.s3.object.key, imageSizes.logo_medium))
        await Promise.all(QUEUE);
        // insert into dynamodb
        await insertLogoIntoDynamoDb(QUEUE, splitPath[1])


        // else if (splitPath[1] === 'MENU') {
        //     // build key for image-processing ${vendorid}/MENU/${category}/${itemid}.jpg
        //     const item_id = splitPath[3].replace(".jpg", '')
        //     const image_key = `${splitPath[0]}/MENU/${splitPath[3]}`
        //     const MENU_QUEUE = []
        //     MENU_QUEUE.push(convertImage(imageData, image_key, imageSizes.logo_small))
        //     MENU_QUEUE.push(convertImage(imageData, image_key, imageSizes.logo_medium))
        //     MENU_QUEUE.push(convertImage(imageData, image_key, imageSizes.logo_large))
        //     await Promise.all(MENU_QUEUE);
        //     // insert into dynamodb
        //     await insertMenuItemDynamoDb(MENU_QUEUE, splitPath[0], splitPath[2], item_id)
        // }

    } catch (error) {
        console.error(`Error thrown while processing the ${record.s3.object.key} image. ${error.message}`, error)
        return { statusCode: 500, body: JSON.stringify({ message: error.message }) }
    }
}

exports.main = async (event) => {
    console.log(JSON.stringify(event))
    if (!event || !event.Records ||
        !event.Records[0] || !event.Records[0].s3 ||
        !event.Records[0].s3.object || !event.Records[0].s3.object.key) {
        const errorMessage = 'Please provide a valid S3 trigger event'
        console.error(errorMessage)
        return { statusCode: 500, body: JSON.stringify({ message: errorMessage }) }
    }

    const PROMISE_QUEUE = []

    for (let i = 0; i < event.Records.length; i++) {
        const record = event.Records[i];
        PROMISE_QUEUE.push(processRecord(record))
    }
    await Promise.all(PROMISE_QUEUE);
}