import { DefineAuthChallengeTriggerEvent } from 'aws-lambda';

exports.handler = async (event: DefineAuthChallengeTriggerEvent) => {

    console.log(JSON.stringify(event));

    if (event.request.userNotFound) {
        event.response.failAuthentication = true;
        event.response.issueTokens = false;
        return event
    }

    const limitDevice = event.request.session.filter((val) => val.challengeName === 'CUSTOM_CHALLENGE' && val.challengeMetadata?.match('LIMIT_DEVICE'));
    if (limitDevice && limitDevice.length > 0) {
        event.response = {
            challengeName: 'CUSTOM_CHALLENGE',
            failAuthentication: false,
            issueTokens: limitDevice[0].challengeResult
        }
        return event
    }

    // Filter the session to find the OTP_CHALLANGE
    const otpChallange = event.request.session.filter((val) => val.challengeName === 'CUSTOM_CHALLENGE' && val.challengeMetadata?.match('OTP_CODE'));
    // If OTP flow is not triggred.
    if (otpChallange.length === 0) {
        event.response = {
            challengeName: 'CUSTOM_CHALLENGE',
            failAuthentication: false,
            issueTokens: false
        }
        return event
    }

    const otpVerified = event.request.session.find((val) => val.challengeName === 'CUSTOM_CHALLENGE' && val.challengeMetadata?.match('OTP_CODE') && val.challengeResult)

    // OTP Already Verified Proceed with Next Step.
    if (otpVerified) {
        event.response = {
            challengeName: 'CUSTOM_CHALLENGE',
            failAuthentication: false,
            issueTokens: false
        }
        return event
    }

    // OTP is not verified 3 times it will stop the request here
    if (otpChallange.length >= 3) {
        event.response = {
            challengeName: 'CUSTOM_CHALLENGE',
            failAuthentication: true,
            issueTokens: false
        }
        return event
    } else {
        event.response = {
            challengeName: 'CUSTOM_CHALLENGE',
            failAuthentication: false,
            issueTokens: false
        }
        return event
    }



}