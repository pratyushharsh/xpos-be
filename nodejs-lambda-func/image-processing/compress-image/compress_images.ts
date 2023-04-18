import { S3Event } from "aws-lambda";

const TARGET_BUCKET = process.env.TARGET_BUCKET
const QUALITY = process.env.OUT_IMAGE_QUALITY || 80

const s3Service = require('./service/s3Service').s3Service
const sharpService = require('./service/sharpService').sharpService


// Process Each Record

async function processRecord(bucket: string, key: string) {
    const imageData = await s3Service.getObject(bucket, key)
    // Compress the image using sharp
    const convertedImageData = await sharpService.convertImage(imageData, parseInt(`${QUALITY}`))
    // Save the image to the target bucket
    await s3Service.saveImage(TARGET_BUCKET, key, convertedImageData)

    // Delete Object After Processing
}


exports.handler = async (event: S3Event) => {
    // For each record
    let PQ = []
    for (const record of event.Records) {
        // Get the bucket name
        const bucket = record.s3.bucket.name;
        // Get the object key
        const key = decodeURIComponent(record.s3.object.key.replace(/\+/g, " "));
        console.log(`Compressing Bucket: ${bucket} Key: ${key}`);
        PQ.push(processRecord(bucket, key))
    }
    await Promise.all(PQ)
}