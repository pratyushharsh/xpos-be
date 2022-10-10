import { CreateAuthChallengeTriggerEvent } from 'aws-lambda';
import AWS from 'aws-sdk';

const cognito = new AWS.CognitoIdentityServiceProvider({
    apiVersion: "2016-04-18",
});


exports.handler = async (event: CreateAuthChallengeTriggerEvent) => {
    
    console.log("CUSTOM_CHALLENGE_LAMBDA", JSON.stringify(event));

    // Create OTP If it is not present in the session.
    // Filter the session to find the OTP_CHALLANGE
    const otpChallange = event.request.session.filter((val) => val.challengeName === 'CUSTOM_CHALLENGE' && val.challengeMetadata?.match('OTP_CODE'));
    if (otpChallange.length == 0) {
        event.response = {
            privateChallengeParameters: {
                'challenge_step': 'VERIFY_OTP',
                'secret_code': '666666'
            },
            challengeMetadata: 'OTP_CODE-666666',
            publicChallengeParameters: {
                'challenge_step': 'VERIFY_OTP',
                'number_of_retries': `${3 - otpChallange.length}`
            }
        }
        return event;
    }

    // If It is retry then fetch the existing match parameter.
    const otpVerified = event.request.session.find((val) => val.challengeName === 'CUSTOM_CHALLENGE' && val.challengeMetadata?.match('OTP_CODE') && val.challengeResult)

    // If OTP is verified then sent the list of loggedin device
    if (!otpVerified) {
        // Fetch the existing parameter
        let exisitingChallange = otpChallange[0]
        event.response = {
            privateChallengeParameters: {
                'challenge_step': 'VERIFY_OTP',
                'secret_code': '666666'
            },
            challengeMetadata: exisitingChallange.challengeMetadata!,
            publicChallengeParameters: {
                'challenge_step': 'VERIFY_OTP',
                'number_of_retries': `${3 - otpChallange.length}`
            }
        }
        return event;
    } else {

        let existingDevice: any[] | undefined = []
        try {
            let devices = await cognito.adminListDevices({
                Username: event.userName,
                UserPoolId: event.userPoolId
            }).promise();
            console.log(JSON.stringify(devices));
            existingDevice = devices.Devices?.map((d) => {
                return {
                    'device_key': d.DeviceKey,
                    'device_name': d.DeviceAttributes?.find((e) => e.Name == 'device_name')?.Value
                }
            }) || []
        } catch(e) {
            console.error(e)
        }
        

        event.response = {
            privateChallengeParameters: {
                'challenge_step': 'VERIFY_DEVICE',
                'device_list': JSON.stringify(existingDevice)
            },
            challengeMetadata: 'LIMIT_DEVICE',
            publicChallengeParameters: {
                'challenge_step': 'VERIFY_DEVICE',
                'no_of_existing_device': `${existingDevice.length}`,
                'device_list': JSON.stringify(existingDevice)
            }
        }
        return event;
    }
}