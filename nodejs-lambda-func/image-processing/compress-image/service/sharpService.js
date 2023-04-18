const sharp = require('sharp')

class SharpService {

    constructor(sharp) {
        this.sharp = sharp
    }

    async convertImage(imageData, quality = 50) {
        return await this.sharp(imageData)
            .png({lossless: true, quality: quality})
            .toBuffer()
    }

    async isValidSize(imageData, imagesSize) {
        const metaData = await this.sharp(imageData).metadata()
        return metaData.width >= imagesSize[0] && metaData.height >= imagesSize[1]
    }
}
module.exports = {
    sharpService: new SharpService(sharp),
    SharpService: SharpService
}