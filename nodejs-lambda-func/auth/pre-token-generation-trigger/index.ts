import { PreTokenGenerationTriggerEvent } from 'aws-lambda';

import AWS from 'aws-sdk';

const cognito = new AWS.CognitoIdentityServiceProvider({
    apiVersion: "2016-04-18",
});

exports.handler = async (event: PreTokenGenerationTriggerEvent) => {
    console.log(JSON.stringify(typeof event));
    console.log(JSON.stringify(event));

    event.response = {
        claimsOverrideDetails: {
            claimsToAddOrOverride: {
                "deviceKey": "customClaimValue"
            }
        }
    }

    return event
}