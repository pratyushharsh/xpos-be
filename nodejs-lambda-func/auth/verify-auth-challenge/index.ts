import { VerifyAuthChallengeResponseTriggerEvent } from 'aws-lambda';

import AWS from 'aws-sdk';

const cognito = new AWS.CognitoIdentityServiceProvider({
    apiVersion: "2016-04-18",
});

exports.handler = async (event: VerifyAuthChallengeResponseTriggerEvent) => {
    console.log(JSON.stringify(event));

    event.response.answerCorrect = false;
    // Find the type of verification step
    let challangeStep: string| undefined = event.request.privateChallengeParameters['challenge_step']

    if (challangeStep && challangeStep === 'VERIFY_OTP') {
        const expectedAnswer = event.request.privateChallengeParameters.secret_code;
        if (event.request.challengeAnswer && event.request.challengeAnswer === expectedAnswer) {
            event.response.answerCorrect = true;
        } else {
            event.response.answerCorrect = false;
        }
        return event;
    } else if (challangeStep && challangeStep === 'VERIFY_DEVICE') {
        // Unregister device
        let req = event.request.challengeAnswer && event.request.challengeAnswer;
        let devices = req.split(";")
        let PQ: any[] = []
        try {
            devices.forEach(key => {
                PQ.push(cognito.adminForgetDevice({
                    UserPoolId: event.userPoolId,
                    Username: event.userName,
                    DeviceKey: key
                }).promise())
            });
            await Promise.all(PQ);
        } catch (e) {
            console.error(e);
        }
        event.response.answerCorrect = true;
        return event;
    }
    
    return event;
}